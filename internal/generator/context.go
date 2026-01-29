package generator

import "strings"

// Context holds generator state: output buffer, indentation, and label counter.
// It is scoped to a generator instance to avoid globals and ensure deterministic output.
type Context struct {
	labelCounter int
	indent       int
	out          *strings.Builder
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

// NextLabel returns a new deterministic label id.
func (c *Context) NextLabel() int {
	c.labelCounter++
	return c.labelCounter
}

// String returns the current output buffer.
func (c *Context) String() string { return c.out.String() }
