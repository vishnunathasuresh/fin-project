# Fin Language

Fin is a Fish-inspired DSL that compiles to Windows Batch. This page summarizes the language surface and canonical formatting.

## Files
- Extension: `.fin`
- UTF-8 text
- Statements separated by newlines
- Blocks terminated explicitly with `end`

## Comments
```
# line comment
```

## Variables
```
set name "Vishnu"
set count 10
```
- All runtime values are strings
- References use `$name`
- Shadowing allowed inside functions

## Literals
- String: `"hello"`
- Number: `123`
- Boolean: `true`, `false`
- List: `[1, 2, 3]`
- Map: `{name: "bob", age: 20}`

## Expressions
Operators (precedence high→low):
```
!                unary not
**               exponent (right-assoc)
* /              multiply/divide
+ -              add/sub
< <= > >=        comparisons
== !=            equality
&& ||            logical and/or
```
Other forms: indexing (`$nums[0]`), property (`$user.name`), grouping (`(expr)`), existence condition (`exists "file.txt"`).

## Statements
- `set name expr`
- `echo expr`
- `run "command"`
- `if cond ... [else ...] end`
- `for i in expr .. expr ... end`
- `while expr ... end`
- `fn name {params} ... end`
- `return [expr]` (inside functions only)
- `break` / `continue` (inside loops)
- Call: `ident {expr}`

Conditions: `if exists "file.txt"` or `if expr`.

## Functions
```
fn greet name
    echo "Hello $name"
end

greet "Bob"
```
- Positional params only
- Locals shadow outer vars
- Return values reserved for future versions (return jumps to function end)

## Formatting (canonical)
`fin fmt` produces deterministic output:
- 4-space indentation
- Space around `..` in ranges (`for i in 1 .. 3`)
- Blank line between top-level functions
- Stable spacing in maps/lists (`{k: v}`, `[a, b]`)
- Unary/binary rendered with spaces (`(a + b)`, `!flag`)

## String interpolation (generator)
Inside string literals, `$name` becomes `%name%`; `$$` becomes literal `$`.

For full architectural rules, see [rules.md](rules.md).
