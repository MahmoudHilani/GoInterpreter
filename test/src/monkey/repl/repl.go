package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"github.com/MahmoudHilani/GoInterpreter/test/src/monkey/evaluator"
	"github.com/MahmoudHilani/GoInterpreter/test/src/monkey/lexer"
	"github.com/MahmoudHilani/GoInterpreter/test/src/monkey/object"
	"github.com/MahmoudHilani/GoInterpreter/test/src/monkey/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}

	}
}

func StartAPI(in string, env *object.Environment) string {
		fmt.Print(PROMPT)
		l := lexer.New(in)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			return printParserErrorsAPI(p.Errors())
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			return evaluated.Inspect()
		}
	return "No result"
}

const MONKEY_FACE = `        __,__
	 .---. .-"   "-. .---.
	/ .. \/ .-. .-. \/ .. \
   |  |  ' | / Y \ | '  |  |
   |  \  \ \ 0 | 0 / /  /  |
    \'- ,\.-"""""""-./, -'/
     ' '-' /_ ^ ^ _\ '-' '
		  | \._ _./ |
		  \ \ '~' / /
		  '._'-=-'_.'
			'-----'
`
func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func printParserErrorsAPI(errors []string) string {


	out := strings.Join(errors, "\n")
	return out
}

