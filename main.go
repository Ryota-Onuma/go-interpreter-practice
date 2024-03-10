package main

import (
	"bufio"
	"io"
	"os"
	"os/signal"
	"syscall"

	"go-interpreter-practice/evaluator"
	"go-interpreter-practice/object"
	"go-interpreter-practice/parser"
	"go-interpreter-practice/scanner"
)

func main() {
	execWithFile("test.onu")
	// args := os.Args
	// if len(args) == 2 {
	// 	filePath := args[1]
	// 	execWithFile(filePath)
	// } else {
	// 	callREPL()
	// }
}

func execWithFile(filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		io.WriteString(os.Stdout, err.Error())
	}
	defer f.Close()
	var data = make([]byte, 1024)
	n, err := f.Read(data)
	if err != nil {
		io.WriteString(os.Stdout, err.Error())
	}
	s := scanner.NewScanner(string(data[:n]))
	parser := parser.NewParser(s)
	program, err := parser.Parse()
	if err != nil {
		io.WriteString(os.Stdout, err.Error())
		io.WriteString(os.Stdout, "\n")
		return
	}

	if len(parser.GetErrors()) > 0 {
		for _, e := range parser.GetErrors() {
			io.WriteString(os.Stdout, e.Error())
			io.WriteString(os.Stdout, "\n")
		}
		return
	}
	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		io.WriteString(os.Stdout, evaluated.String())
		io.WriteString(os.Stdout, "\n")
	}
}

func callREPL() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// 対話型のRPELループを開始
	reader := bufio.NewReader(os.Stdin)
	done := make(chan struct{})
	go func() {
		for {
			io.WriteString(os.Stdout, ">> ")
			input, err := reader.ReadString('\n')
			if err != nil {
				done <- struct{}{}
				return
			}

			// 入力が "exit" だったらループを終了
			if input == "exit" {
				io.WriteString(os.Stdout, "RPELを終了します。\n")
				done <- struct{}{}
				return
			}

			s := scanner.NewScanner(input)
			parser := parser.NewParser(s)
			program, err := parser.Parse()
			if err != nil {
				io.WriteString(os.Stdout, err.Error())
			}
			env := object.NewEnvironment()
			evaluated := evaluator.Eval(program, env)
			if evaluated != nil {
				io.WriteString(os.Stdout, evaluated.String())
				io.WriteString(os.Stdout, "\n")
			}
		}
	}()

	select {
	case <-sigChan:
		os.Exit(1)
	case <-done:
	}
}
