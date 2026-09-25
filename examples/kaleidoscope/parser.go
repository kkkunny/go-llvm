package main

import (
	"fmt"
	"strings"
)

// Expr 表达式节点；String 产出 -dp 展示用的括号化形式
type Expr interface {
	fmt.Stringer
}

// NumberExpr 数字字面量
type NumberExpr struct{ Value float64 }

// String 返回数字字面量的文本形式（float64 的默认格式）。
func (e *NumberExpr) String() string { return fmt.Sprintf("%v", e.Value) }

// VariableExpr 变量引用
type VariableExpr struct{ Name string }

// String 返回被引用变量的名字。
func (e *VariableExpr) String() string { return e.Name }

// BinaryExpr 二元表达式（含赋值 = 与自定义算符）
type BinaryExpr struct {
	Op          rune
	Left, Right Expr
}

// String 返回括号化的中缀形式，如 (1 + (2 * 3))。
func (e *BinaryExpr) String() string { return fmt.Sprintf("(%s %c %s)", e.Left, e.Op, e.Right) }

// CallExpr 函数调用
type CallExpr struct {
	Callee string
	Args   []Expr
}

// String 返回调用形式，实参以 ", " 连接，如 foo(1, 2.5)；无参时形如 foo()。
func (e *CallExpr) String() string {
	args := make([]string, len(e.Args))
	for i, arg := range e.Args {
		args[i] = arg.String()
	}
	return fmt.Sprintf("%s(%s)", e.Callee, strings.Join(args, ", "))
}

// ConditionalExpr if/then/else 表达式
type ConditionalExpr struct {
	Cond, Then, Else Expr
}

// String 返回括号化的条件表达式，如 (if a then b else c)。
func (e *ConditionalExpr) String() string {
	return fmt.Sprintf("(if %s then %s else %s)", e.Cond, e.Then, e.Else)
}

// ForExpr for/in 循环；Step 为 nil 时默认 1
type ForExpr struct {
	Var        string
	Start, End Expr
	Step       Expr
	Body       Expr
}

// String 返回括号化的 for 表达式：Step 非 nil 时为
// (for i = 1, (i < 5), 2 in body)，否则省略步长部分。
func (e *ForExpr) String() string {
	if e.Step != nil {
		return fmt.Sprintf("(for %s = %s, %s, %s in %s)", e.Var, e.Start, e.End, e.Step, e.Body)
	}
	return fmt.Sprintf("(for %s = %s, %s in %s)", e.Var, e.Start, e.End, e.Body)
}

// VarDef var 声明中的一个变量；Init 为 nil 时默认 0
type VarDef struct {
	Name string
	Init Expr
}

// VarInExpr var..in 作用域
type VarInExpr struct {
	Vars []VarDef
	Body Expr
}

// String 返回括号化的 var..in 表达式；每个绑定有初值时为 "name = init"，
// 否则只有 "name"，如 (var a = 5, b in (a * b))。
func (e *VarInExpr) String() string {
	defs := make([]string, len(e.Vars))
	for i, def := range e.Vars {
		if def.Init != nil {
			defs[i] = def.Name + " = " + def.Init.String()
		} else {
			defs[i] = def.Name
		}
	}
	return fmt.Sprintf("(var %s in %s)", strings.Join(defs, ", "), e.Body)
}

// Prototype 函数签名
type Prototype struct {
	Name string
	Args []string
	IsOp bool // 自定义 binary/unary 算符
	Prec int  // 二元算符优先级（仅 binary）
}

// FuncDef 顶层解析结果；Body 为 nil 表示 extern 声明
type FuncDef struct {
	Proto  Prototype
	Body   Expr
	IsAnon bool
}

// DefaultPrecedence 内置算符优先级（对齐 LLVM 教程：= 最低，* / 最高）
func DefaultPrecedence() map[rune]int {
	return map[rune]int{'=': 2, '<': 10, '>': 10, '+': 20, '-': 20, '*': 40, '/': 40}
}

// ParseError 语法错误；Pos 为出错 token 在（去注释后的）token 流中的下标
type ParseError struct {
	Msg string
	Pos int
}

// Error 返回语法错误消息，格式为 "<Msg> (at token <Pos>)"，Pos 为出错 token 在 token 流中的下标。
func (e *ParseError) Error() string { return fmt.Sprintf("%s (at token %d)", e.Msg, e.Pos) }

// Parser 递归下降 + 优先级爬升解析器
type Parser struct {
	tokens []Token
	pos    int
	prec   map[rune]int
}

// Parse 解析一行输入；prec 为算符优先级表，binary 声明会原地更新它
func Parse(input string, prec map[rune]int) (*FuncDef, error) {
	tokens, err := Lex(input)
	if err != nil {
		return nil, err
	}
	return ParseTokens(tokens, prec)
}

// ParseTokens 解析 token 流（供已做词法展示的调用方复用）
func ParseTokens(tokens []Token, prec map[rune]int) (*FuncDef, error) {
	p := &Parser{tokens: skipComments(tokens), prec: prec}

	fn, err := p.parseFunc()
	if err != nil {
		return nil, err
	}
	if tok := p.curr(); tok.Kind != TokEOF {
		return nil, p.errf("unexpected token %s after expression", tok)
	}
	return fn, nil
}

// skipComments 注释是展示用 trivia，解析前剔除（原地过滤）
func skipComments(tokens []Token) []Token {
	out := tokens[:0]
	for _, tok := range tokens {
		if tok.Kind != TokComment {
			out = append(out, tok)
		}
	}
	return out
}

// curr 当前 token；越界返回 EOF
func (p *Parser) curr() Token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return Token{Kind: TokEOF}
}

func (p *Parser) advance() {
	if p.pos < len(p.tokens) {
		p.pos++
	}
}

func (p *Parser) errf(format string, args ...any) error {
	return &ParseError{Msg: fmt.Sprintf(format, args...), Pos: p.pos}
}

// expect 要求当前 token 为指定种类并前进
func (p *Parser) expect(kind TokenKind) error {
	if p.curr().Kind != kind {
		return p.errf("expected %s, got %s", tokenKindNames[kind], p.curr())
	}
	p.advance()
	return nil
}

// getTokPrec 当前二元算符的优先级；-1 表示未注册（无法作为二元算符）
func (p *Parser) getTokPrec() int {
	if tok := p.curr(); tok.Kind == TokOp {
		if prec, ok := p.prec[tok.Op]; ok {
			return prec
		}
	}
	return -1
}

// parseFunc 按首 token 分派：def / extern / 顶层表达式
func (p *Parser) parseFunc() (*FuncDef, error) {
	switch p.curr().Kind {
	case TokDef:
		return p.parseDef()
	case TokExtern:
		return p.parseExtern()
	default:
		return p.parseToplevel()
	}
}

// parsePrototype 解析函数签名；binary 声明会注册算符优先级
func (p *Parser) parsePrototype() (Prototype, error) {
	var proto Prototype
	switch tok := p.curr(); tok.Kind {
	case TokIdent:
		p.advance()
		proto.Name = tok.Ident

	case TokBinary:
		p.advance()
		op, err := p.expectOp("binary")
		if err != nil {
			return Prototype{}, err
		}
		proto.Name = "binary" + string(op)
		proto.IsOp = true
		if tok := p.curr(); tok.Kind == TokNumber {
			proto.Prec = int(tok.Num)
			p.advance()
		}
		p.prec[op] = proto.Prec

	case TokUnary:
		p.advance()
		op, err := p.expectOp("unary")
		if err != nil {
			return Prototype{}, err
		}
		proto.Name = "unary" + string(op)
		proto.IsOp = true

	default:
		return Prototype{}, p.errf("expected identifier in prototype declaration, got %s", tok)
	}

	if err := p.expect(TokLParen); err != nil {
		return Prototype{}, err
	}
	if p.curr().Kind == TokRParen {
		p.advance()
		return proto, nil
	}
	for {
		tok := p.curr()
		if tok.Kind != TokIdent {
			return Prototype{}, p.errf("expected identifier in parameter declaration, got %s", tok)
		}
		proto.Args = append(proto.Args, tok.Ident)
		p.advance()

		switch p.curr().Kind {
		case TokRParen:
			p.advance()
			return proto, nil
		case TokComma:
			p.advance()
		default:
			return Prototype{}, p.errf("expected ',' or ')' in prototype declaration, got %s", p.curr())
		}
	}
}

func (p *Parser) expectOp(what string) (rune, error) {
	tok := p.curr()
	if tok.Kind != TokOp {
		return 0, p.errf("expected operator in custom %s declaration, got %s", what, tok)
	}
	p.advance()
	return tok.Op, nil
}

// parseDef 解析 def 函数定义
func (p *Parser) parseDef() (*FuncDef, error) {
	p.advance() // def
	proto, err := p.parsePrototype()
	if err != nil {
		return nil, err
	}
	body, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return &FuncDef{Proto: proto, Body: body}, nil
}

// parseExtern 解析 extern 函数声明
func (p *Parser) parseExtern() (*FuncDef, error) {
	p.advance() // extern
	proto, err := p.parsePrototype()
	if err != nil {
		return nil, err
	}
	return &FuncDef{Proto: proto}, nil
}

// parseToplevel 把顶层表达式包成匿名函数（函数名由 Session 赋唯一后缀）
func (p *Parser) parseToplevel() (*FuncDef, error) {
	body, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return &FuncDef{Proto: Prototype{Name: "__anon_expr"}, Body: body, IsAnon: true}, nil
}

// parseExpr 解析完整表达式：一元前缀 + 二元算符
func (p *Parser) parseExpr() (Expr, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	return p.parseBinary(0, left)
}

// parseUnary 解析一元算符前缀（编译为 unary<op> 调用）或主表达式
func (p *Parser) parseUnary() (Expr, error) {
	tok := p.curr()
	if tok.Kind != TokOp {
		return p.parsePrimary()
	}
	p.advance()
	operand, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	return &CallExpr{Callee: "unary" + string(tok.Op), Args: []Expr{operand}}, nil
}

// parsePrimary 解析主表达式
func (p *Parser) parsePrimary() (Expr, error) {
	switch tok := p.curr(); tok.Kind {
	case TokIdent:
		return p.parseIdent()
	case TokNumber:
		p.advance()
		return &NumberExpr{Value: tok.Num}, nil
	case TokLParen:
		return p.parseParen()
	case TokIf:
		return p.parseConditional()
	case TokFor:
		return p.parseFor()
	case TokVar:
		return p.parseVarIn()
	default:
		return nil, p.errf("unknown expression at %s", tok)
	}
}

// parseIdent 解析标识符：后跟 ( 为函数调用，否则为变量引用
func (p *Parser) parseIdent() (Expr, error) {
	name := p.curr().Ident
	p.advance()
	if p.curr().Kind != TokLParen {
		return &VariableExpr{Name: name}, nil
	}
	p.advance() // (

	if p.curr().Kind == TokRParen {
		p.advance()
		return &CallExpr{Callee: name}, nil
	}
	var args []Expr
	for {
		arg, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)

		switch p.curr().Kind {
		case TokComma:
			p.advance()
		case TokRParen:
			p.advance()
			return &CallExpr{Callee: name, Args: args}, nil
		default:
			return nil, p.errf("expected ',' or ')' in function call, got %s", p.curr())
		}
	}
}

// parseParen 解析括号表达式
func (p *Parser) parseParen() (Expr, error) {
	p.advance() // (
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if err := p.expect(TokRParen); err != nil {
		return nil, err
	}
	return expr, nil
}

// parseConditional 解析 if/then/else
func (p *Parser) parseConditional() (Expr, error) {
	p.advance() // if
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if err := p.expect(TokThen); err != nil {
		return nil, err
	}
	then, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if err := p.expect(TokElse); err != nil {
		return nil, err
	}
	els, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return &ConditionalExpr{Cond: cond, Then: then, Else: els}, nil
}

// parseFor 解析 for i = start, end[, step] in body
func (p *Parser) parseFor() (Expr, error) {
	p.advance() // for

	nameTok := p.curr()
	if nameTok.Kind != TokIdent {
		return nil, p.errf("expected identifier in for loop, got %s", nameTok)
	}
	p.advance()

	if tok := p.curr(); tok.Kind != TokOp || tok.Op != '=' {
		return nil, p.errf("expected '=' in for loop, got %s", tok)
	}
	p.advance()

	start, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if err := p.expect(TokComma); err != nil {
		return nil, err
	}
	end, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	var step Expr
	if p.curr().Kind == TokComma {
		p.advance()
		if step, err = p.parseExpr(); err != nil {
			return nil, err
		}
	}

	if err := p.expect(TokIn); err != nil {
		return nil, err
	}
	body, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return &ForExpr{Var: nameTok.Ident, Start: start, End: end, Step: step, Body: body}, nil
}

// parseVarIn 解析 var a = ..., b, ... in body
func (p *Parser) parseVarIn() (Expr, error) {
	p.advance() // var

	var vars []VarDef
	for {
		nameTok := p.curr()
		if nameTok.Kind != TokIdent {
			return nil, p.errf("expected identifier in 'var..in' declaration, got %s", nameTok)
		}
		p.advance()

		var def VarDef
		def.Name = nameTok.Ident
		if tok := p.curr(); tok.Kind == TokOp && tok.Op == '=' {
			p.advance()
			init, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			def.Init = init
		}
		vars = append(vars, def)

		switch p.curr().Kind {
		case TokComma:
			p.advance()
		case TokIn:
			p.advance()
			body, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			return &VarInExpr{Vars: vars, Body: body}, nil
		default:
			return nil, p.errf("expected ',' or 'in' in variable declaration, got %s", p.curr())
		}
	}
}

// parseBinary 优先级爬升：解析以 left 为左操作数的二元表达式
func (p *Parser) parseBinary(prec int, left Expr) (Expr, error) {
	for {
		currPrec := p.getTokPrec()
		if currPrec < prec {
			return left, nil
		}
		op := p.curr().Op
		p.advance()

		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		if nextPrec := p.getTokPrec(); currPrec < nextPrec {
			if right, err = p.parseBinary(currPrec+1, right); err != nil {
				return nil, err
			}
		}
		left = &BinaryExpr{Op: op, Left: left, Right: right}
	}
}
