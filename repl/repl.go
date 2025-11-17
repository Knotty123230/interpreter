package repl

import (
	"bufio"
	"fmt"
	"interpreter/lexer"
	"interpreter/token"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	fmt.Fprintf(out, PROMPT)
	scanner := bufio.NewScanner(in)
	scan := scanner.Scan()
	if !scan {
		return
	}

	text := scanner.Text()
	l := lexer.New(text)

	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		fmt.Fprintf(out, "%+v\n", tok)
	}

}
