package main

import (
	"fmt"
	"strconv"
)

// TokenKind 词法单元种类
type TokenKind int

const (
	TokEOF TokenKind = iota
	TokDef
	TokExtern
	TokIf
	TokThen
	TokElse
	TokFor
	TokIn
	TokUnary
	TokBinary
	TokVar
	TokIdent
	TokNumber
	TokOp
	TokLParen
	TokRParen
	TokComma
	TokComment
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
