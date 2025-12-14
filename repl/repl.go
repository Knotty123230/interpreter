package repl

import (
	"bufio"
	"fmt"
	"interpreter/evaluator"
	"interpreter/lexer"
	"interpreter/object"
	"interpreter/parser"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	env := object.NewEnvironment()
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprintf(out, PROMPT)
		scanner.Scan()
		line := scanner.Text()

		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors) > 0 {
			printParserErrors(out, p.Errors)
			continue
		}
		obj := evaluator.Eval(program, env)
		if obj != nil {
			io.WriteString(out, obj.Inspect())
			io.WriteString(out, "\n")
		}
	}

}

func printParserErrors(out io.Writer, s []string) {
	for _, msg := range s {
		io.WriteString(out, "\t"+msg+"\n")
		// fmt.Fprintf(out, "\t"+msg+"\n")
	}
}
