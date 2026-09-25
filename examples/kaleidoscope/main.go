package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	var (
		eval            = flag.String("e", "", "求值单个表达式后退出")
		displayLexer    = flag.Bool("dl", false, "显示词法分析结果")
		displayParser   = flag.Bool("dp", false, "显示语法分析结果")
		displayCompiler = flag.Bool("dc", false, "显示生成的 LLVM IR")
	)
	flag.Parse()

	session, err := NewSession(Options{
		DisplayLexer:    *displayLexer,
		DisplayParser:   *displayParser,
		DisplayCompiler: *displayCompiler,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	defer session.Close()

	if *eval != "" {
		runLine(os.Stdout, session, *eval)
		return
	}

	fmt.Println("Kaleidoscope REPL（输入 exit 或 quit 退出）")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("?> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "exit") || strings.HasPrefix(trimmed, "quit") {
			break
		}
		runLine(os.Stdout, session, line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
	}
}

// runLine 求值一行输入并打印结果；错误按词法/语法与编译两段归类（对齐 inkwell 输出）
func runLine(w io.Writer, s *Session, input string) {
	val, isExpr, err := s.Eval(input)
	if err != nil {
		var lexErr *LexError
		var parseErr *ParseError
		if errors.As(err, &lexErr) || errors.As(err, &parseErr) {
			fmt.Fprintln(w, "!> Error parsing expression:", err)
		} else {
			fmt.Fprintln(w, "!> Error compiling function:", err)
		}
		return
	}
	if isExpr {
		// 对齐 inkwell/教程的数值显示（不用 %v 的科学计数法）
		fmt.Fprintf(w, "=> %s\n", strconv.FormatFloat(val, 'f', -1, 64))
	}
}
