package main

import (
	"fmt"
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
	defer f.Close()
	data := make([]byte, 1024)
	count, err := f.Read(data)
	if err != nil {
		panic(err)
	}

	s := scanner.NewScanner(string(data[:count]))
	fmt.Println(string(data[:count]))
	fmt.Println()
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
