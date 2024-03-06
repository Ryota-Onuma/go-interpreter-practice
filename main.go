package main

import (
	"fmt"
	"lox-by-go/parser"
	"lox-by-go/scanner"
	"os"
)

func main() {
	fmt.Println()
	runFile()
	fmt.Println()
}

func runCmd() {
	fmt.Println()
	s := scanner.NewScanner("var a = 1;")
	s.ScanTokens()
	if len(s.GetErrors()) > 0 {
		for _, err := range s.GetErrors() {
			fmt.Println(err)
		}
	}

	for _, token := range s.Tokens() {
		fmt.Println(token)
	}
	fmt.Println()
}

func runFile() {
	fmt.Println()
	f, err := os.Open("test.onu")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	data := make([]byte, 1024)
	count, err := f.Read(data)
	if err != nil {
		panic(err)
	}

	s := scanner.NewScanner(string(data[:count]))
	parser := parser.NewParser(s)
	program, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(program)

	fmt.Println()
}
