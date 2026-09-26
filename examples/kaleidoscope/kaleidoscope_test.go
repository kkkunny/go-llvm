package main

import (
	"bytes"
	"errors"
	"io"
	"math"
	"reflect"
	"strings"
	"testing"
)

func TestLex(t *testing.T) {
	tokens, err := Lex("def foo(x) x + 1.5 # 注释\n")
	if err != nil {
		t.Fatal(err)
	}
	want := []Token{
		{Kind: TokDef},
		{Kind: TokIdent, Ident: "foo"},
		{Kind: TokLParen},
		{Kind: TokIdent, Ident: "x"},
		{Kind: TokRParen},
		{Kind: TokIdent, Ident: "x"},
		{Kind: TokOp, Op: '+'},
		{Kind: TokNumber, Num: 1.5},
		{Kind: TokComment},
		{Kind: TokEOF},
	}
	if !reflect.DeepEqual(tokens, want) {
		t.Fatalf("tokens = %v, want %v", tokens, want)
	}
}

func TestLexError(t *testing.T) {
	for _, src := range []string{"1.2.3", "."} {
		var lexErr *LexError
		if _, err := Lex(src); !errors.As(err, &lexErr) {
			t.Fatalf("Lex(%q) error = %v, want *LexError", src, err)
		}
	}
}

func parseString(t *testing.T, input string) string {
	t.Helper()
	fn, err := Parse(input, DefaultPrecedence())
	if err != nil {
		t.Fatalf("Parse(%q) error: %v", input, err)
	}
	if fn.Body == nil {
		return "<extern " + fn.Proto.Name + ">"
	}
	return fn.Body.String()
}

func TestParseExpr(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"1 + 2 * 3", "(1 + (2 * 3))"},
		{"1 * 2 + 3", "((1 * 2) + 3)"},
		{"1 + 2 + 3", "((1 + 2) + 3)"},
		{"a = b = 3", "((a = b) = 3)"},
		{"1 + -2", "(1 + unary-(2))"},
		{"foo(1, 2.5)", "foo(1, 2.5)"},
		{"foo()", "foo()"},
		{"(1 + 2) * 3", "((1 + 2) * 3)"},
		{"if a then b else c", "(if a then b else c)"},
		{"if a < b then a else b + 1", "(if (a < b) then a else (b + 1))"},
		{"var a = 5, b in a * b", "(var a = 5, b in (a * b))"},
		{"for i = 1, i < 5, 2 in i", "(for i = 1, (i < 5), 2 in i)"},
		{"for i = 1, i < 5 in i + 1", "(for i = 1, (i < 5) in (i + 1))"},
	}
	for _, tt := range tests {
		if got := parseString(t, tt.input); got != tt.want {
			t.Errorf("Parse(%q) = %s, want %s", tt.input, got, tt.want)
		}
	}
}

func TestParsePrototype(t *testing.T) {
	fn, err := Parse("def add(a, b) a + b", DefaultPrecedence())
	if err != nil {
		t.Fatal(err)
	}
	if fn.Proto.Name != "add" || !reflect.DeepEqual(fn.Proto.Args, []string{"a", "b"}) || fn.IsAnon {
		t.Fatalf("unexpected def: %+v", fn)
	}

	fn, err = Parse("extern cos(x)", DefaultPrecedence())
	if err != nil {
		t.Fatal(err)
	}
	if fn.Body != nil || fn.Proto.Name != "cos" || !reflect.DeepEqual(fn.Proto.Args, []string{"x"}) {
		t.Fatalf("unexpected extern: %+v", fn)
	}

	prec := DefaultPrecedence()
	fn, err = Parse("def binary^ 40 (lhs, rhs) lhs * rhs", prec)
	if err != nil {
		t.Fatal(err)
	}
	if !fn.Proto.IsOp || fn.Proto.Name != "binary^" || fn.Proto.Prec != 40 {
		t.Fatalf("unexpected binary op: %+v", fn.Proto)
	}
	if prec['^'] != 40 {
		t.Fatalf("precedence of '^' = %d, want 40", prec['^'])
	}

	fn, err = Parse("def unary!(v) 0 - v", prec)
	if err != nil {
		t.Fatal(err)
	}
	if !fn.Proto.IsOp || fn.Proto.Name != "unary!" {
		t.Fatalf("unexpected unary op: %+v", fn.Proto)
	}
}

func TestParseAnon(t *testing.T) {
	fn, err := Parse("1 + 2", DefaultPrecedence())
	if err != nil {
		t.Fatal(err)
	}
	if !fn.IsAnon || len(fn.Proto.Args) != 0 {
		t.Fatalf("top-level expression should be anonymous: %+v", fn)
	}
}

func TestParseErrors(t *testing.T) {
	tests := []string{
		"1 +",            // 表达式不完整
		"(1",             // 缺右括号
		"1 2",            // 表达式后有多余 token
		"def 1(x) x",     // 签名缺标识符
		"def f(a,) a",    // 参数列表格式错误
		"foo(1 2)",       // 调用参数缺逗号
		"if a then b",    // 缺 else
		"for i = 1, 2 i", // 缺 in
		"var a = 1 a",    // 缺 in
		"binary % 40",    // 缺少函数签名
	}
	for _, input := range tests {
		var parseErr *ParseError
		if _, err := Parse(input, DefaultPrecedence()); !errors.As(err, &parseErr) {
			t.Errorf("Parse(%q) error = %v, want *ParseError", input, err)
		}
	}
}

func TestParseUnregisteredOpStops(t *testing.T) {
	// 未通过 binary 声明注册的算符不能出现在二元位置
	if _, err := Parse("1 % 2", DefaultPrecedence()); err == nil {
		t.Fatal("unregistered operator should fail to parse")
	}
}

// ===== 端到端：经 Session.Eval 走完整编译 + JIT 链路 =====

func newTestSession(t *testing.T, opts Options) *Session {
	t.Helper()
	if opts.Output == nil {
		opts.Output = io.Discard
	}
	s, err := NewSession(opts)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func evalExpr(t *testing.T, s *Session, input string) float64 {
	t.Helper()
	val, isExpr, err := s.Eval(input)
	if err != nil {
		t.Fatalf("Eval(%q) error: %v", input, err)
	}
	if !isExpr {
		t.Fatalf("Eval(%q) should be an expression", input)
	}
	return val
}

func evalDef(t *testing.T, s *Session, input string) {
	t.Helper()
	if _, isExpr, err := s.Eval(input); err != nil {
		t.Fatalf("Eval(%q) error: %v", input, err)
	} else if isExpr {
		t.Fatalf("Eval(%q) should be a definition", input)
	}
}

func TestEvalExpr(t *testing.T) {
	s := newTestSession(t, Options{})
	tests := []struct {
		input string
		want  float64
	}{
		{"1 + 2 * 2", 5},
		{"(1 + 2) * 3", 9},
		{"10 / 4", 2.5},
		{"1 - 2 - 3", -4}, // 左结合
		{"1 < 2", 1},
		{"2 < 1", 0},
		{"2 > 1", 1},
		{"if 1 then 2 else 3", 2},
		{"if 0 then 1 else 2", 2},
		{"if 1 then if 0 then 2 else 3 else 4", 3}, // 嵌套 if 的 PHI 前驱块
		{"if 0 then 1 else if 1 then 2 else 3", 2},
		{"1 + 1 # 尾随注释", 2},
	}
	for _, tt := range tests {
		if got := evalExpr(t, s, tt.input); got != tt.want {
			t.Errorf("Eval(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestEvalVariables(t *testing.T) {
	s := newTestSession(t, Options{})
	tests := []struct {
		input string
		want  float64
	}{
		{"var a = 5, b = 10 in a * b", 50},
		{"var a in a", 0},                                 // 无初值默认 0
		{"var a = 3 in a = a * a", 9},                     // 赋值返回右值
		{"var a = 1 in var a = 2 in a", 2},                // 内层遮蔽
		{"var a = 1 in (var a = 2 in a) + a", 3},          // 退出作用域后恢复旧绑定
		{"var i = 99 in (for i = 1, i < 2 in i) + i", 99}, // 循环变量遮蔽恢复
	}
	for _, tt := range tests {
		if got := evalExpr(t, s, tt.input); got != tt.want {
			t.Errorf("Eval(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestEvalFunction(t *testing.T) {
	s := newTestSession(t, Options{})
	evalDef(t, s, "def fib(n) if n < 2 then n else fib(n - 1) + fib(n - 2)")
	if got := evalExpr(t, s, "fib(10)"); got != 55 {
		t.Fatalf("fib(10) = %v, want 55", got)
	}

	// 同名重定义替换历史定义
	evalDef(t, s, "def answer() 41")
	if got := evalExpr(t, s, "answer()"); got != 41 {
		t.Fatalf("answer() = %v, want 41", got)
	}
	evalDef(t, s, "def answer() 42")
	if got := evalExpr(t, s, "answer()"); got != 42 {
		t.Fatalf("redefined answer() = %v, want 42", got)
	}

	// 反复求值：每次 eval 都经历 addModule/removeModule
	for i := 0; i < 3; i++ {
		if got := evalExpr(t, s, "answer() + 0"); got != 42 {
			t.Fatalf("repeat %d: answer() + 0 = %v, want 42", i, got)
		}
	}
}

func TestEvalCustomOperator(t *testing.T) {
	s := newTestSession(t, Options{})
	evalDef(t, s, "def binary^ 40 (a, b) a * b")
	if got := evalExpr(t, s, "2 ^ 3 + 1"); got != 7 { // ^ 优先级高于 +
		t.Fatalf("2 ^ 3 + 1 = %v, want 7", got)
	}

	evalDef(t, s, "def binary~ 40 (a, b) a - b")
	if got := evalExpr(t, s, "10 ~ 3 ~ 2"); got != 5 { // 同优先级左结合
		t.Fatalf("10 ~ 3 ~ 2 = %v, want 5", got)
	}

	evalDef(t, s, "def unary-(v) 0 - v")
	if got := evalExpr(t, s, "-3"); got != -3 {
		t.Fatalf("-3 = %v, want -3", got)
	}
	if got := evalExpr(t, s, "1 + -2"); got != -1 {
		t.Fatalf("1 + -2 = %v, want -1", got)
	}
}

func TestEvalForLoop(t *testing.T) {
	var buf bytes.Buffer
	s := newTestSession(t, Options{Output: &buf})

	evalDef(t, s, "extern putchard(x)")
	evalDef(t, s, "def star(n) for i = 1, i < n in putchard(42)")

	// 教程语义：先执行 body，再判断结束条件；body 至少执行一次
	if got := evalExpr(t, s, "star(4)"); got != 0 {
		t.Fatalf("star(4) = %v, want 0", got)
	}
	if got := buf.String(); got != "****" {
		t.Fatalf("star(4) printed %q, want %q", got, "****")
	}
	buf.Reset()
	if got := evalExpr(t, s, "star(1)"); got != 0 {
		t.Fatalf("star(1) = %v, want 0", got)
	}
	if got := buf.String(); got != "*" {
		t.Fatalf("star(1) printed %q, want %q", got, "*")
	}
}

func TestEvalExtern(t *testing.T) {
	s := newTestSession(t, Options{})
	evalDef(t, s, "extern sin(x)")
	evalDef(t, s, "extern cos(x)")
	if got := evalExpr(t, s, "sin(0)"); math.Abs(got) > 1e-9 {
		t.Fatalf("sin(0) = %v, want 0", got)
	}
	if got := evalExpr(t, s, "cos(0)"); math.Abs(got-1) > 1e-9 {
		t.Fatalf("cos(0) = %v, want 1", got)
	}

	// 用户函数调用 extern（进程符号解析 + 跨函数调用）
	evalDef(t, s, "def trig_identity(x) sin(x) * sin(x) + cos(x) * cos(x)")
	if got := evalExpr(t, s, "trig_identity(0.7)"); math.Abs(got-1) > 1e-9 {
		t.Fatalf("trig_identity(0.7) = %v, want 1", got)
	}
}

func TestEvalPrintd(t *testing.T) {
	var buf bytes.Buffer
	s := newTestSession(t, Options{Output: &buf})
	evalDef(t, s, "extern printd(x)")
	if got := evalExpr(t, s, "printd(42)"); got != 42 {
		t.Fatalf("printd(42) = %v, want 42", got)
	}
	if got := buf.String(); got != "42\n" {
		t.Fatalf("printd printed %q, want %q", got, "42\n")
	}
}

func TestEvalErrors(t *testing.T) {
	s := newTestSession(t, Options{})
	evalDef(t, s, "def one(x) x")

	tests := []string{
		"1 +",          // 语法错：表达式不完整
		"1 2",          // 语法错：多余 token
		"a",            // 未定义变量
		"foo(1)",       // 未定义函数
		"1 $ 2",        // 未注册算符
		"f(1) = 2",     // 赋值左值非法
		"one(1, 2)",    // 实参数量不符
		"var a = 1 in", // 语法错：var body 缺失
	}
	for _, input := range tests {
		if _, _, err := s.Eval(input); err == nil {
			t.Errorf("Eval(%q) should fail", input)
		}
	}
}

func TestRunLine(t *testing.T) {
	var buf bytes.Buffer
	s := newTestSession(t, Options{})

	runLine(&buf, s, "1 + 2 * 2")
	if got := buf.String(); got != "=> 5\n" {
		t.Fatalf("runLine output = %q", got)
	}

	buf.Reset()
	evalDef(t, s, "def fib(n) if n < 2 then n else fib(n - 1) + fib(n - 2)")
	runLine(&buf, s, "fib(40)")
	if got := buf.String(); got != "=> 102334155\n" { // 不使用科学计数法
		t.Fatalf("runLine fib output = %q", got)
	}

	buf.Reset()
	runLine(&buf, s, "1 +")
	if got := buf.String(); !strings.HasPrefix(got, "!> Error parsing expression: ") {
		t.Fatalf("parse error output = %q", got)
	}

	buf.Reset()
	runLine(&buf, s, "a")
	if got := buf.String(); !strings.HasPrefix(got, "!> Error compiling function: ") {
		t.Fatalf("compile error output = %q", got)
	}
}

func TestEvalDisplays(t *testing.T) {
	var buf bytes.Buffer
	s := newTestSession(t, Options{
		DisplayLexer:    true,
		DisplayParser:   true,
		DisplayCompiler: true,
		Output:          &buf,
	})
	if got := evalExpr(t, s, "1 + 2 * 2"); got != 5 {
		t.Fatalf("Eval = %v, want 5", got)
	}
	out := buf.String()
	for _, want := range []string{
		"-> tokens:",
		"Number(1)",
		"-> expression: (1 + (2 * 2))",
		"-> IR:",
		"define double @__anon_expr.1()",
		"ret double 5.000000e+00", // 优化管线常量折叠
	} {
		if !strings.Contains(out, want) {
			t.Errorf("display output missing %q:\n%s", want, out)
		}
	}
}
