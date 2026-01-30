package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCLI_Build_Valid(t *testing.T) {
	tmp := t.TempDir()
	finPath := filepath.Join(tmp, "valid.fin")
	if err := os.WriteFile(finPath, []byte("set x 1\n"), 0644); err != nil {
		t.Fatalf("write fin: %v", err)
	}
	outPath := filepath.Join(tmp, "out.bat")
	cmd := exec.Command("go", "run", "./cmd/fin", "build", "-o", outPath, finPath)
	cmd.Dir = projectRoot(t)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build valid failed (code=%d): %v\noutput: %s", exitCode(err), err, output)
	}
	if _, err := os.Stat(outPath); err != nil {
		t.Fatalf("expected output file: %v", err)
	}
}

func TestCLI_Fmt(t *testing.T) {
	tmp := t.TempDir()
	finPath := filepath.Join(tmp, "fmt.fin")
	src := strings.Join([]string{
		"fn a",
		"    set x 1",
		"end",
		"fn b",
		"    for i in 1..3",
		"        echo $i",
		"    end",
		"end",
		"",
	}, "\n")
	if err := os.WriteFile(finPath, []byte(src), 0644); err != nil {
		t.Fatalf("write fin: %v", err)
	}
	cmd := exec.Command("go", "run", "./cmd/fin", "fmt", finPath)
	cmd.Dir = projectRoot(t)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected fmt to succeed (code=%d): %v\noutput: %s", exitCode(err), err, output)
	}
	expected := strings.Join([]string{
		"fn a",
		"    set x 1",
		"end",
		"",
		"fn b",
		"    for i in 1 .. 3",
		"        echo $i",
		"    end",
		"end",
	}, "\n")
	if string(output) != expected {
		t.Fatalf("unexpected fmt output:\nexpected:\n%s\ngot:\n%s", expected, string(output))
	}
}

func TestCLI_Build_Invalid(t *testing.T) {
	tmp := t.TempDir()
	finPath := filepath.Join(tmp, "invalid.fin")
	if err := os.WriteFile(finPath, []byte("return 1\n"), 0644); err != nil {
		t.Fatalf("write fin: %v", err)
	}
	cmd := exec.Command("go", "run", "./cmd/fin", "build", finPath)
	cmd.Dir = projectRoot(t)
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected build invalid to fail; output: %s", output)
	}
	if code := exitCode(err); code == 0 {
		t.Fatalf("expected non-zero exit code, got 0; output: %s", output)
	}
	if !strings.Contains(string(output), "return") {
		t.Fatalf("expected return error in output, got: %s", output)
	}
}

func TestCLI_Check_Invalid(t *testing.T) {
	tmp := t.TempDir()
	finPath := filepath.Join(tmp, "invalid.fin")
	if err := os.WriteFile(finPath, []byte("return 1\n"), 0644); err != nil {
		t.Fatalf("write fin: %v", err)
	}
	cmd := exec.Command("go", "run", "./cmd/fin", "check", finPath)
	cmd.Dir = projectRoot(t)
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected check invalid to fail; output: %s", output)
	}
	if code := exitCode(err); code == 0 {
		t.Fatalf("expected non-zero exit code, got 0; output: %s", output)
	}
	if !strings.Contains(string(output), "return") {
		t.Fatalf("expected return error in output, got: %s", output)
	}
}

func TestCLI_AST_Valid(t *testing.T) {
	tmp := t.TempDir()
	finPath := filepath.Join(tmp, "valid.fin")
	if err := os.WriteFile(finPath, []byte("set x 1\n"), 0644); err != nil {
		t.Fatalf("write fin: %v", err)
	}
	cmd := exec.Command("go", "run", "./cmd/fin", "ast", finPath)
	cmd.Dir = projectRoot(t)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected ast to succeed (code=%d): %v\noutput: %s", exitCode(err), err, output)
	}
	if !strings.Contains(string(output), "Program") {
		t.Fatalf("expected AST output, got: %s", output)
	}
}

func projectRoot(t *testing.T) string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("cannot determine caller")
	}
	root := filepath.Join(filepath.Dir(file), "..", "..")
	abspath, err := filepath.Abs(root)
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	return abspath
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	var ee *exec.ExitError
	if ok := errors.As(err, &ee); ok {
		return ee.ExitCode()
	}
	return -1
}
