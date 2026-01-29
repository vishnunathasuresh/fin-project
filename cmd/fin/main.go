package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/vishnunath-suresh/fin-project/internal/ast"
	"github.com/vishnunath-suresh/fin-project/internal/generator"
	"github.com/vishnunath-suresh/fin-project/internal/lexer"
	"github.com/vishnunath-suresh/fin-project/internal/parser"
	"github.com/vishnunath-suresh/fin-project/internal/sema"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "build":
		buildCmd(os.Args[2:])
	case "check":
		checkCmd(os.Args[2:])
	case "ast":
		astCmd(os.Args[2:])
	case "version":
		fmt.Println(version)
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  fin build <file.fin> [-o output.bat]\n")
	fmt.Fprintf(os.Stderr, "  fin check <file.fin>\n")
	fmt.Fprintf(os.Stderr, "  fin ast <file.fin>\n")
	fmt.Fprintf(os.Stderr, "  fin version\n")
}

func buildCmd(args []string) {
	flags := flag.NewFlagSet("build", flag.ExitOnError)
	var outPath string
	flags.StringVar(&outPath, "o", "", "output batch file")
	if err := flags.Parse(args); err != nil {
		os.Exit(1)
	}
	if flags.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "build requires exactly one input file")
		os.Exit(1)
	}
	inPath := flags.Arg(0)

	prog, err := loadAndAnalyze(inPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	out, err := generate(prog)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if outPath == "" {
		base := filepath.Base(inPath)
		outPath = base[:len(base)-len(filepath.Ext(base))] + ".bat"
	}
	if err := ioutil.WriteFile(outPath, []byte(out), 0644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

func checkCmd(args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "check requires exactly one input file")
		os.Exit(1)
	}
	prog, err := loadAndAnalyze(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// If generate detects unsupported nodes, surface it as an error even in check.
	if _, err := generate(prog); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

func astCmd(args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "ast requires exactly one input file")
		os.Exit(1)
	}
	prog, err := loadAndAnalyze(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(ast.Format(prog))
	os.Exit(0)
}

func loadAndAnalyze(path string) (*ast.Program, error) {
	src, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	l := lexer.New(string(src))
	toks := parser.CollectTokens(l)
	p := parser.New(toks)
	prog := p.ParseProgram()
	if perrs := p.Errors(); len(perrs) > 0 {
		return nil, errors.Join(perrs...)
	}

	a := sema.New()
	if err := a.Analyze(prog); err != nil {
		return nil, err
	}

	return prog, nil
}

func generate(prog *ast.Program) (string, error) {
	g := generator.NewBatchGenerator()
	return g.Generate(prog)
}
