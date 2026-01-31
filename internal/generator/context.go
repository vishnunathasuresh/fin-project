package generator

import "strings"

// Context holds generator state: output buffer, indentation, and label counter.
// It is scoped to a generator instance to avoid globals and ensure deterministic output.
type Context struct {
	labelCounter int
	indent       int
	out          *strings.Builder
	loopStack    []loopLabels
	returnStack  []returnTarget
}

// NewContext constructs an empty generator context.
func NewContext() *Context {
	return &Context{out: &strings.Builder{}}
}

// pushIndent increases indentation level for subsequent lines.
func (c *Context) pushIndent() { c.indent++ }

// popIndent decreases indentation level (no-op at zero).
func (c *Context) popIndent() {
	if c.indent > 0 {
		c.indent--
	}
}

// emitLine writes a line with current indentation and a trailing newline.
func (c *Context) emitLine(s string) {
	for i := 0; i < c.indent; i++ {
		c.out.WriteString("    ")
	}
	c.out.WriteString(s)
	c.out.WriteString("\n")
}

// emitRawLine writes a line with no indentation (useful for labels).
func (c *Context) emitRawLine(s string) {
	c.out.WriteString(s)
	c.out.WriteString("\n")
}

// NextLabel returns a new deterministic label id.
func (c *Context) NextLabel() int {
	c.labelCounter++
	return c.labelCounter
}

// String returns the current output buffer.
func (c *Context) String() string { return c.out.String() }

// loopLabels represents break/continue targets for the current loop.
type loopLabels struct {
	breakLabel    string
	continueLabel string
}

func (c *Context) pushLoop(breakLabel, continueLabel string) {
	c.loopStack = append(c.loopStack, loopLabels{breakLabel: breakLabel, continueLabel: continueLabel})
}

func (c *Context) popLoop() {
	if len(c.loopStack) == 0 {
		return
	}
	c.loopStack = c.loopStack[:len(c.loopStack)-1]
}

func (c *Context) currentLoop() (loopLabels, bool) {
	if len(c.loopStack) == 0 {
		return loopLabels{}, false
	}
	return c.loopStack[len(c.loopStack)-1], true
}

type returnTarget struct {
	label   string
	tempVar string
	outVar  string
}

func (c *Context) pushReturn(label, tempVar, outVar string) {
	c.returnStack = append(c.returnStack, returnTarget{label: label, tempVar: tempVar, outVar: outVar})
}

func (c *Context) popReturn() {
	if len(c.returnStack) == 0 {
		return
	}
	c.returnStack = c.returnStack[:len(c.returnStack)-1]
}

func (c *Context) currentReturn() (returnTarget, bool) {
	if len(c.returnStack) == 0 {
		return returnTarget{}, false
	}
	return c.returnStack[len(c.returnStack)-1], true
}
