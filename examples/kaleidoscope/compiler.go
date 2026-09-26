package main

import (
	"fmt"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/ir"
)

// Compiler 把 AST 编译为 LLVM IR；每个函数用一个实例，变量作用域随函数结束
type Compiler struct {
	ctx     *llvm.Context
	module  *ir.Module
	builder *ir.Builder
	f64     llvm.FloatType
	fn      ir.Function
	vars    map[string]ir.Alloca // 变量名 → 栈槽
}

// Compile 在模块中编译一个函数定义；Body 为 nil 时只声明（extern）
func Compile(ctx *llvm.Context, module *ir.Module, builder *ir.Builder, def *FuncDef) (ir.Function, error) {
	f64 := ctx.Float(llvm.FloatDouble)
	params := make([]llvm.AnyType, len(def.Proto.Args))
	for i := range params {
		params[i] = f64
	}
	fn := module.NewFunction(def.Proto.Name, ctx.Fn(f64, params, false))
	for i, name := range def.Proto.Args {
		fn.ParamAs[llvm.FloatT](uint(i)).SetName(name)
	}
	if def.Body == nil {
		return fn, nil
	}

	c := &Compiler{ctx: ctx, module: module, builder: builder, f64: f64, fn: fn, vars: map[string]ir.Alloca{}}
	c.builder.MoveToEnd(fn.NewBlock("entry"))

	// 参数先落到入口块栈槽，函数体内一切变量读写都走 alloca
	for i, name := range def.Proto.Args {
		alloca, err := c.createEntryAlloca(name)
		if err != nil {
			return fn, err
		}
		c.builder.Store(fn.ParamAs[llvm.FloatT](uint(i)), alloca)
		c.vars[name] = alloca
	}

	body, err := c.compileExpr(def.Body)
	if err != nil {
		return fn, err
	}
	c.builder.Ret(body)
	return fn, nil
}

// createEntryAlloca 在入口块开头分配栈槽（对应 inkwell 的 create_entry_block_alloca）：
// 用临时 builder 把 alloca 插到首条指令之前，避免占用条件分支的执行路径
func (c *Compiler) createEntryAlloca(name string) (ir.Alloca, error) {
	entry, ok := c.fn.EntryBlock()
	if !ok {
		return ir.Alloca{}, fmt.Errorf("function %q has no entry block", c.fn.Name())
	}
	tmp := ir.NewBuilder(c.ctx)
	defer tmp.Close()
	if first, ok := entry.FirstInst(); ok {
		tmp.MoveBefore(first)
	} else {
		tmp.MoveToEnd(entry)
	}
	return tmp.Alloca(c.f64, name), nil
}

// compileExpr 把表达式编译为 f64 值；错误为语义错误（未定义变量/函数等）
func (c *Compiler) compileExpr(expr Expr) (llvm.Value[llvm.FloatT], error) {
	switch e := expr.(type) {
	case *NumberExpr:
		return c.f64.Const(e.Value).Value, nil

	case *VariableExpr:
		alloca, ok := c.vars[e.Name]
		if !ok {
			return llvm.Value[llvm.FloatT]{}, fmt.Errorf("undefined variable %q", e.Name)
		}
		return c.builder.Load(alloca, c.f64, e.Name).Value, nil

	case *VarInExpr:
		return c.compileVarIn(e)

	case *BinaryExpr:
		return c.compileBinary(e)

	case *CallExpr:
		fn, ok := c.module.GetFunction(e.Callee)
		if !ok {
			return llvm.Value[llvm.FloatT]{}, fmt.Errorf("unknown function %q", e.Callee)
		}
		if want := len(fn.Signature().Params()); len(e.Args) != want {
			return llvm.Value[llvm.FloatT]{}, fmt.Errorf("%q expects %d argument(s), got %d", e.Callee, want, len(e.Args))
		}
		args := make([]llvm.AnyValue, len(e.Args))
		for i, arg := range e.Args {
			v, err := c.compileExpr(arg)
			if err != nil {
				return llvm.Value[llvm.FloatT]{}, err
			}
			args[i] = v.Dyn()
		}
		return c.builder.Call[llvm.FloatT](fn, args, "tmp").Value, nil

	case *ConditionalExpr:
		return c.compileConditional(e)

	case *ForExpr:
		return c.compileFor(e)

	default:
		return llvm.Value[llvm.FloatT]{}, fmt.Errorf("unsupported expression %T", expr)
	}
}

// compileVarIn 编译 var..in：每个变量在入口块分配栈槽，body 结束时恢复被遮蔽的绑定
func (c *Compiler) compileVarIn(e *VarInExpr) (llvm.Value[llvm.FloatT], error) {
	type savedBinding struct {
		name    string
		alloca  ir.Alloca
		existed bool
	}
	saved := make([]savedBinding, 0, len(e.Vars))

	for _, def := range e.Vars {
		init := c.f64.Const(0).Value
		if def.Init != nil {
			var err error
			if init, err = c.compileExpr(def.Init); err != nil {
				return llvm.Value[llvm.FloatT]{}, err
			}
		}
		alloca, err := c.createEntryAlloca(def.Name)
		if err != nil {
			return llvm.Value[llvm.FloatT]{}, err
		}
		c.builder.Store(init, alloca)

		old, existed := c.vars[def.Name]
		saved = append(saved, savedBinding{name: def.Name, alloca: old, existed: existed})
		c.vars[def.Name] = alloca
	}

	body, err := c.compileExpr(e.Body)
	if err != nil {
		return llvm.Value[llvm.FloatT]{}, err
	}
	for _, b := range saved {
		if b.existed {
			c.vars[b.name] = b.alloca
		} else {
			delete(c.vars, b.name)
		}
	}
	return body, nil
}

// compileBinary 编译二元表达式；'=' 为赋值，其余查找 binary<op> 函数
func (c *Compiler) compileBinary(e *BinaryExpr) (llvm.Value[llvm.FloatT], error) {
	if e.Op == '=' {
		target, ok := e.Left.(*VariableExpr)
		if !ok {
			return llvm.Value[llvm.FloatT]{}, fmt.Errorf("expected variable on the left-hand side of '='")
		}
		alloca, ok := c.vars[target.Name]
		if !ok {
			return llvm.Value[llvm.FloatT]{}, fmt.Errorf("undefined variable %q", target.Name)
		}
		val, err := c.compileExpr(e.Right)
		if err != nil {
			return llvm.Value[llvm.FloatT]{}, err
		}
		c.builder.Store(val, alloca)
		return val, nil // 赋值的值即右值（对齐 inkwell/教程）
	}

	lhs, err := c.compileExpr(e.Left)
	if err != nil {
		return llvm.Value[llvm.FloatT]{}, err
	}
	rhs, err := c.compileExpr(e.Right)
	if err != nil {
		return llvm.Value[llvm.FloatT]{}, err
	}

	switch e.Op {
	case '+':
		return c.builder.FAdd(lhs, rhs, "tmpadd"), nil
	case '-':
		return c.builder.FSub(lhs, rhs, "tmpsub"), nil
	case '*':
		return c.builder.FMul(lhs, rhs, "tmpmul"), nil
	case '/':
		return c.builder.FDiv(lhs, rhs, "tmpdiv"), nil
	case '<':
		return c.builder.UIToFP(c.builder.FCmp(llvm.FloatULT, lhs, rhs, "tmpcmp"), c.f64, "tmpbool"), nil
	case '>':
		return c.builder.UIToFP(c.builder.FCmp(llvm.FloatULT, rhs, lhs, "tmpcmp"), c.f64, "tmpbool"), nil
	default:
		fn, ok := c.module.GetFunction("binary" + string(e.Op))
		if !ok {
			return llvm.Value[llvm.FloatT]{}, fmt.Errorf("undefined binary operator %q", string(e.Op))
		}
		return c.builder.Call[llvm.FloatT](fn, []llvm.AnyValue{lhs.Dyn(), rhs.Dyn()}, "tmpbin").Value, nil
	}
}

// compileConditional 编译 if/then/else：then/else 写各自基本块，合并块用 PHI 取值
func (c *Compiler) compileConditional(e *ConditionalExpr) (llvm.Value[llvm.FloatT], error) {
	cond, err := c.compileExpr(e.Cond)
	if err != nil {
		return llvm.Value[llvm.FloatT]{}, err
	}
	// 非 0 为真（对齐 inkwell：ONE != 0.0）
	condInt := c.builder.FCmp(llvm.FloatONE, cond, c.f64.Const(0), "ifcond")

	thenBlk := c.fn.NewBlock("then")
	elseBlk := c.fn.NewBlock("else")
	contBlk := c.fn.NewBlock("ifcont")
	c.builder.CondBr(condInt, thenBlk, elseBlk)

	c.builder.MoveToEnd(thenBlk)
	thenVal, err := c.compileExpr(e.Then)
	if err != nil {
		return llvm.Value[llvm.FloatT]{}, err
	}
	c.builder.Br(contBlk)
	// 嵌套 if 会让构建器停在内层合并块，这里取真实前驱块
	thenEnd, _ := c.builder.CurrentBlock()

	c.builder.MoveToEnd(elseBlk)
	elseVal, err := c.compileExpr(e.Else)
	if err != nil {
		return llvm.Value[llvm.FloatT]{}, err
	}
	c.builder.Br(contBlk)
	elseEnd, _ := c.builder.CurrentBlock()

	c.builder.MoveToEnd(contBlk)
	phi := c.builder.PHI(c.f64, "iftmp")
	phi.AddIncoming(
		ir.Incoming[llvm.FloatT]{Value: thenVal, Block: thenEnd},
		ir.Incoming[llvm.FloatT]{Value: elseVal, Block: elseEnd},
	)
	return phi.Value, nil
}

// compileFor 编译 for 循环（对齐教程语义：先执行 body，再算结束条件，最后步进；
// 结束条件非 0 继续循环）；循环表达式的值恒为 0
func (c *Compiler) compileFor(e *ForExpr) (llvm.Value[llvm.FloatT], error) {
	start, err := c.compileExpr(e.Start)
	if err != nil {
		return llvm.Value[llvm.FloatT]{}, err
	}
	alloca, err := c.createEntryAlloca(e.Var)
	if err != nil {
		return llvm.Value[llvm.FloatT]{}, err
	}
	c.builder.Store(start, alloca)

	loopBlk := c.fn.NewBlock("loop")
	afterBlk := c.fn.NewBlock("afterloop")
	c.builder.Br(loopBlk)
	c.builder.MoveToEnd(loopBlk)

	old, existed := c.vars[e.Var]
	c.vars[e.Var] = alloca
	if _, err := c.compileExpr(e.Body); err != nil { // 循环体的值被丢弃
		return llvm.Value[llvm.FloatT]{}, err
	}

	step := c.f64.Const(1).Value
	if e.Step != nil {
		if step, err = c.compileExpr(e.Step); err != nil {
			return llvm.Value[llvm.FloatT]{}, err
		}
	}

	// 结束条件先于步进求值（end 引用的循环变量是步进前的值）
	endCond, err := c.compileExpr(e.End)
	if err != nil {
		return llvm.Value[llvm.FloatT]{}, err
	}
	cur := c.builder.Load(alloca, c.f64, e.Var).Value
	next := c.builder.FAdd(cur, step, "nextvar")
	c.builder.Store(next, alloca)

	loopCond := c.builder.FCmp(llvm.FloatONE, endCond, c.f64.Const(0), "loopcond")
	c.builder.CondBr(loopCond, loopBlk, afterBlk)
	c.builder.MoveToEnd(afterBlk)

	if existed {
		c.vars[e.Var] = old
	} else {
		delete(c.vars, e.Var)
	}
	return c.f64.Const(0).Value, nil
}
