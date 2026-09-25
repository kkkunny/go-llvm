package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kkkunny/go-llvm"
	"github.com/kkkunny/go-llvm/ir"
	"github.com/kkkunny/go-llvm/jit"
	"github.com/kkkunny/go-llvm/pass"
	"github.com/kkkunny/go-llvm/target"
)

// optPipeline 对齐 inkwell 示例的优化管线（opt -passes 语法）
const optPipeline = "instcombine,reassociate,gvn,simplifycfg,mem2reg"

// Options 会话选项
type Options struct {
	DisplayLexer    bool      // 显示词法分析结果
	DisplayParser   bool      // 显示语法分析结果
	DisplayCompiler bool      // 显示生成的 LLVM IR
	Output          io.Writer // 输出目标；nil 表示 os.Stdout
}

// Session 一次 REPL 会话：常驻 JIT + 历史定义 + 算符优先级表
type Session struct {
	jit     *jit.LLJIT
	opts    Options
	prec    map[rune]int
	defs    []*FuncDef // 按定义顺序保存；同名重定义时替换
	anonSeq int
}

// NewSession 创建会话：初始化宿主目标、JIT、宿主进程符号与 putchard/printd
func NewSession(opts Options) (*Session, error) {
	// JIT 依赖已注册的宿主目标
	if err := target.InitNative(); err != nil {
		return nil, err
	}
	j, err := jit.NewLLJIT()
	if err != nil {
		return nil, err
	}
	s := &Session{jit: j, opts: opts, prec: DefaultPrecedence()}
	if err := j.AddProcessSymbols(); err != nil {
		j.Close()
		return nil, err
	}
	if err := j.MapFunc("putchard", s.putchard); err != nil {
		j.Close()
		return nil, err
	}
	if err := j.MapFunc("printd", s.printd); err != nil {
		j.Close()
		return nil, err
	}
	return s, nil
}

// Close 释放 JIT
func (s *Session) Close() error {
	return s.jit.Close()
}

func (s *Session) out() io.Writer {
	if s.opts.Output != nil {
		return s.opts.Output
	}
	return os.Stdout
}

// putchard 打印字符并原样返回（JIT 内 extern putchard 声明可解析到本函数）
func (s *Session) putchard(x float64) float64 {
	fmt.Fprintf(s.out(), "%c", byte(x))
	return x
}

// printd 打印数值并原样返回
func (s *Session) printd(x float64) float64 {
	fmt.Fprintf(s.out(), "%v\n", x)
	return x
}

// Eval 解析并执行一行输入；isExpr 为 true 时返回匿名表达式的求值结果
func (s *Session) Eval(input string) (val float64, isExpr bool, err error) {
	tokens, err := Lex(input)
	if err != nil {
		return 0, false, err
	}
	if s.opts.DisplayLexer {
		fmt.Fprintf(s.out(), "-> tokens: %v\n", tokens)
	}

	def, err := ParseTokens(tokens, s.prec)
	if err != nil {
		return 0, false, err
	}
	if s.opts.DisplayParser {
		s.display(def)
	}

	// 每个求值单元独立 Context/Module：JIT 路径会消费两者，非 JIT 路径由 Close 级联释放。
	// 历史定义整体重编译进当前模块（对齐 inkwell），因此重定义天然生效、无跨模块符号冲突。
	ctx := llvm.NewContext()
	m := ir.NewModule(ctx, "kaleidoscope")
	b := ir.NewBuilder(ctx)
	// 非 JIT 路径的显式释放（builder 先于 module），避免调试层资源审计告警
	discard := func() {
		_ = b.Close()
		_ = m.Close()
		_ = ctx.Close()
	}
	for _, prev := range s.defs {
		if _, err := Compile(ctx, m, b, prev); err != nil {
			discard()
			return 0, false, fmt.Errorf("recompiling %q: %w", prev.Proto.Name, err)
		}
	}

	cur := def
	if def.IsAnon {
		s.anonSeq++
		// '.' 不会出现在标识符里，避免与用户函数重名
		cur = &FuncDef{
			Proto:  Prototype{Name: fmt.Sprintf("__anon_expr.%d", s.anonSeq)},
			Body:   def.Body,
			IsAnon: true,
		}
	}
	fn, err := Compile(ctx, m, b, cur)
	if err != nil {
		discard()
		return 0, false, err
	}
	if err := m.Verify(); err != nil {
		discard()
		return 0, false, err
	}
	if err := pass.RunPasses(m, optPipeline); err != nil {
		discard()
		return 0, false, err
	}
	if s.opts.DisplayCompiler {
		fmt.Fprintf(s.out(), "-> IR:\n%s\n", fn.String())
	}

	if !def.IsAnon {
		s.remember(def)
		discard()
		return 0, false, nil
	}

	// 只把本次求值的模块加入 JIT，求值完立即卸载（对应教程的 addModule/removeModule）
	rt := s.jit.NewResourceTracker()
	if err := rt.AddIRModule(m); err != nil {
		_ = rt.Remove()
		return 0, false, err
	}
	compiled, err := s.jit.Func[func() float64](cur.Proto.Name)
	if err != nil {
		_ = rt.Remove()
		return 0, false, err
	}
	val = compiled()
	if err := rt.Remove(); err != nil {
		return val, true, err
	}
	return val, true, nil
}

// remember 记录定义；同名时替换（有意优于 inkwell：重定义不残留旧函数）
func (s *Session) remember(def *FuncDef) {
	for i, prev := range s.defs {
		if prev.Proto.Name == def.Proto.Name {
			s.defs[i] = def
			return
		}
	}
	s.defs = append(s.defs, def)
}

// display 输出 -dp 展示
func (s *Session) display(def *FuncDef) {
	if def.IsAnon {
		fmt.Fprintf(s.out(), "-> expression: %s\n", def.Body)
		return
	}
	sig := fmt.Sprintf("%s(%s)", def.Proto.Name, strings.Join(def.Proto.Args, ", "))
	if def.Body == nil {
		fmt.Fprintf(s.out(), "-> extern: %s\n", sig)
	} else {
		fmt.Fprintf(s.out(), "-> function: %s = %s\n", sig, def.Body)
	}
}
