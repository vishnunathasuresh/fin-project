package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/vishnunath-suresh/fin-project/internal/ast"
	"github.com/vishnunath-suresh/fin-project/internal/format"
	"github.com/vishnunath-suresh/fin-project/internal/generator"
	"github.com/vishnunath-suresh/fin-project/internal/lexer"
	"github.com/vishnunath-suresh/fin-project/internal/parser"
	"github.com/vishnunath-suresh/fin-project/internal/sema"
	"github.com/vishnunath-suresh/fin-project/internal/version"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	cmd := os.Args[1]
	switch cmd {
	case "build":
		buildCmd(os.Args[2:])
	case "check":
		checkCmd(os.Args[2:])
	case "ast":
		astCmd(os.Args[2:])
	case "fmt":
		fmtCmd(os.Args[2:])
	case "version":
		fmt.Println(version.Version)
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		usage()
		os.Exit(2)
	}
}

// printDiagnostics renders errors with file:line:col style and grouping.
func printDiagnostics(w io.Writer, file string, err error) {
	if err == nil {
		return
	}
	// Flatten joined errors.
	var errs []error
	if j, ok := err.(interface{ Unwrap() []error }); ok {
		errs = j.Unwrap()
	} else {
		errs = []error{err}
	}
	for _, e := range errs {
		msg := strings.TrimSpace(e.Error())
		prefix := colorize("error:", red)
		switch v := e.(type) {
		case interface{ Pos() ast.Pos }:
			pos := v.Pos()
			fmt.Fprintf(w, "%s %s:%d:%d %s\n", prefix, file, pos.Line, pos.Column, msg)
		default:
			fmt.Fprintf(w, "%s %s\n", prefix, msg)
		}
	}
}

var (
	red     = "\x1b[31m"
	reset   = "\x1b[0m"
	noColor = os.Getenv("NO_COLOR") != ""
)

func colorize(s, c string) string {
	if noColor {
		return s
	}
	return c + s + reset
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  fin build <file.fin> [-o output.bat]\n")
	fmt.Fprintf(os.Stderr, "  fin check <file.fin>\n")
	fmt.Fprintf(os.Stderr, "  fin ast <file.fin>\n")
	fmt.Fprintf(os.Stderr, "  fin fmt [-w] <file.fin>\n")
	fmt.Fprintf(os.Stderr, "  fin version\n")
}

func buildCmd(args []string) {
	flags := flag.NewFlagSet("build", flag.ExitOnError)
	flags.SetOutput(os.Stderr)
	var outPath string
	flags.StringVar(&outPath, "o", "", "output batch file")
	if err := flags.Parse(args); err != nil {
		os.Exit(2)
	}
	if flags.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "build requires exactly one input file")
		os.Exit(2)
	}
	inPath := flags.Arg(0)
	if err := validateFinPath(inPath); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	prog, err := loadAndAnalyze(inPath)
	if err != nil {
		printDiagnostics(os.Stderr, inPath, err)
		os.Exit(1)
	}

	out, err := generate(prog)
	if err != nil {
		printDiagnostics(os.Stderr, inPath, err)
		os.Exit(1)
	}

	if outPath == "" {
		base := filepath.Base(inPath)
		outPath = base[:len(base)-len(filepath.Ext(base))] + ".bat"
	}
	if err := atomicWriteFile(outPath, []byte(out), 0644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

func checkCmd(args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "check requires exactly one input file")
		os.Exit(2)
	}
	if err := validateFinPath(args[0]); err != nil {
		printDiagnostics(os.Stderr, args[0], err)
		os.Exit(1)
	}
	prog, err := loadAndAnalyze(args[0])
	if err != nil {
		printDiagnostics(os.Stderr, args[0], err)
		os.Exit(1)
	}

	// If generate detects unsupported nodes, surface it as an error even in check.
	if _, err := generate(prog); err != nil {
		printDiagnostics(os.Stderr, args[0], err)
		os.Exit(1)
	}
	os.Exit(0)
}

func astCmd(args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "ast requires exactly one input file")
		os.Exit(2)
	}
	if err := validateFinPath(args[0]); err != nil {
		printDiagnostics(os.Stderr, args[0], err)
		os.Exit(1)
	}
	prog, err := loadAndAnalyze(args[0])
	if err != nil {
		printDiagnostics(os.Stderr, args[0], err)
		os.Exit(1)
	}
	fmt.Print(ast.Format(prog))
	os.Exit(0)
}

func fmtCmd(args []string) {
	flags := flag.NewFlagSet("fmt", flag.ExitOnError)
	flags.SetOutput(os.Stderr)
	write := flags.Bool("w", false, "write result to (overwrite) file instead of stdout")
	if err := flags.Parse(args); err != nil {
		os.Exit(2)
	}
	if flags.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "fmt requires exactly one input file")
		os.Exit(2)
	}
	path := flags.Arg(0)
	if err := validateFinPath(path); err != nil {
		printDiagnostics(os.Stderr, path, err)
		os.Exit(1)
	}
	prog, err := loadAndAnalyze(path)
	if err != nil {
		printDiagnostics(os.Stderr, path, err)
		os.Exit(1)
	}
	formatted := format.Format(prog)
	if *write {
		info, err := os.Stat(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := atomicWriteFile(path, []byte(formatted), info.Mode()); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		fmt.Print(formatted)
	}
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
		return nil, multiError("parse errors", perrs)
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

// printError renders single or joined errors with simple formatting.
func printError(w io.Writer, err error) {
	if err == nil {
		return
	}
	type unwrapper interface{ Unwrap() []error }
	if u, ok := err.(unwrapper); ok {
		fmt.Fprintln(w, err)
		for _, e := range u.Unwrap() {
			fmt.Fprintf(w, " - %v\n", e)
		}
		return
	}
	fmt.Fprintln(w, err)
}

func multiError(prefix string, errs []error) error {
	if len(errs) == 1 {
		return errs[0]
	}
	return fmt.Errorf("%s: %w", prefix, errors.Join(errs...))
}

func validateFinPath(path string) error {
	if filepath.Ext(path) != ".fin" {
		return fmt.Errorf("input must have .fin extension: %s", path)
	}
	return nil
}

// atomicWriteFile writes data to a temp file in the target directory and renames it into place.
func atomicWriteFile(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	tmp, err := ioutil.TempFile(dir, "fin-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmpPath, perm); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}
