package main

import (
	"fmt"
	"strconv"
)

// TokenKind 词法单元种类
type TokenKind int

// TokenKind 的取值分属五类：关键字、字面量、算符、标点与注释。
const (
	TokEOF     TokenKind = iota // 输入结束（token 流末尾的哨兵）
	TokDef                      // 关键字 def
	TokExtern                   // 关键字 extern
	TokIf                       // 关键字 if
	TokThen                     // 关键字 then
	TokElse                     // 关键字 else
	TokFor                      // 关键字 for
	TokIn                       // 关键字 in
	TokUnary                    // 关键字 unary（自定义一元算符声明）
	TokBinary                   // 关键字 binary（自定义二元算符声明）
	TokVar                      // 关键字 var
	TokIdent                    // 标识符
	TokNumber                   // 数字字面量（float64）
	TokOp                       // 单字符算符
	TokLParen                   // 左括号 (
	TokRParen                   // 右括号 )
	TokComma                    // 逗号 ,
	TokComment                  // 注释（# 至行尾；解析前剔除）
)

// tokenKindNames TokenKind 的展示名（索引与常量一一对应）
var tokenKindNames = [...]string{
	TokEOF:     "EOF",
	TokDef:     "Def",
	TokExtern:  "Extern",
	TokIf:      "If",
	TokThen:    "Then",
	TokElse:    "Else",
	TokFor:     "For",
	TokIn:      "In",
	TokUnary:   "Unary",
	TokBinary:  "Binary",
	TokVar:     "Var",
	TokIdent:   "Ident",
	TokNumber:  "Number",
	TokOp:      "Op",
	TokLParen:  "LParen",
	TokRParen:  "RParen",
	TokComma:   "Comma",
	TokComment: "Comment",
}

// Token 词法单元；Num/Ident/Op 仅在对应种类下有值
type Token struct {
	Kind  TokenKind
	Num   float64 // TokNumber 的字面量
	Ident string  // TokIdent 的文本
	Op    rune    // TokOp 的算符字符
}

// String 返回 token 的文本表示：数字、标识符与算符分别渲染为 Number(1.5)、
// Ident("foo") 与 Op('+') 这类形式，其余种类返回其展示名。
func (t Token) String() string {
	switch t.Kind {
	case TokNumber:
		return fmt.Sprintf("Number(%v)", t.Num)
	case TokIdent:
		return fmt.Sprintf("Ident(%q)", t.Ident)
	case TokOp:
		return fmt.Sprintf("Op(%q)", t.Op)
	default:
		return tokenKindNames[t.Kind]
	}
}

// LexError 词法错误；Index 为出错处的字节偏移
type LexError struct {
	Msg   string
	Index int
}

// Error 返回词法错误消息，格式为 "<Msg> (at index <Index>)"，Index 为出错处的字节偏移。
func (e *LexError) Error() string {
	return fmt.Sprintf("%s (at index %d)", e.Msg, e.Index)
}

// keywords 关键字表
var keywords = map[string]TokenKind{
	"def":    TokDef,
	"extern": TokExtern,
	"if":     TokIf,
	"then":   TokThen,
	"else":   TokElse,
	"for":    TokFor,
	"in":     TokIn,
	"unary":  TokUnary,
	"binary": TokBinary,
	"var":    TokVar,
}

// Lex 把源码切成 token 流（末尾恒有 EOF；注释作为 Comment 保留给 -dl 展示）
func Lex(src string) ([]Token, error) {
	var tokens []Token
	for i, n := 0, len(src); i < n; {
		ch := src[i]
		switch {
		case ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n':
			// 空白
			i++

		case ch == '#':
			// 注释直到行尾
			for i < n && src[i] != '\n' {
				i++
			}
			tokens = append(tokens, Token{Kind: TokComment})

		case ch == '(':
			tokens = append(tokens, Token{Kind: TokLParen})
			i++
		case ch == ')':
			tokens = append(tokens, Token{Kind: TokRParen})
			i++
		case ch == ',':
			tokens = append(tokens, Token{Kind: TokComma})
			i++

		case isDigit(ch) || ch == '.':
			// 数字字面量：连续扫 '.' 与十六进制位（与 inkwell 一致），再交给 ParseFloat
			start := i
			for i < n && (src[i] == '.' || isHexDigit(src[i])) {
				i++
			}
			text := src[start:i]
			v, err := strconv.ParseFloat(text, 64)
			if err != nil {
				return nil, &LexError{Msg: "invalid number literal " + strconv.Quote(text), Index: start}
			}
			tokens = append(tokens, Token{Kind: TokNumber, Num: v})

		case isIdentStart(ch):
			// 标识符或关键字
			start := i
			for i < n && isIdentCont(src[i]) {
				i++
			}
			word := src[start:i]
			if kind, ok := keywords[word]; ok {
				tokens = append(tokens, Token{Kind: kind})
			} else {
				tokens = append(tokens, Token{Kind: TokIdent, Ident: word})
			}

		default:
			// 单字符算符
			tokens = append(tokens, Token{Kind: TokOp, Op: rune(ch)})
			i++
		}
	}
	return append(tokens, Token{Kind: TokEOF}), nil
}

func isDigit(ch byte) bool { return ch >= '0' && ch <= '9' }

func isHexDigit(ch byte) bool {
	return isDigit(ch) || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

func isIdentStart(ch byte) bool {
	return ch == '_' || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isIdentCont(ch byte) bool { return isIdentStart(ch) || isDigit(ch) }
