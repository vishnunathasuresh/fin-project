# Fin v2 — System Instructions (Final, Locked)

> **This document is the authoritative and final specification for Fin v2.**
>
> All human and AI agents must follow this document strictly.
> Any deviation, extension, or reinterpretation requires an explicit design change and version bump.

---

## 1. What Fin v2 Is

Fin v2 is a **pythonic, statically typed automation DSL** that **transpiles** to multiple shell backends:

* **bash**
* **fish**
* **Windows Batch (.bat)**
* **PowerShell (.ps1)**

Fin v2 is not an interpreter and does not execute scripts itself. It always produces shell code.

### Core Properties

* Python-like syntax (low friction)
* Static typing (compile-time safety)
* Explicit side effects via `run()`
* Multi-backend transpilation
* Deterministic output
* Visualization as a first-class feature

Fin v2 is a **direct and intentional evolution of Fin v1**.

---

## 2. What Fin v2 Is NOT

To prevent scope creep, Fin v2 is explicitly **not**:

* An interactive shell
* A shell replacement (bash, zsh, PowerShell)
* A runtime scripting language
* A general-purpose programming language
* A filesystem or process abstraction layer

Fin v2 orchestrates shells; it does not reimplement them.

---

## 3. Design Principles (Non-Negotiable)

All language and tooling decisions must follow these principles:

1. **Determinism over convenience**
2. **Explicitness over magic**
3. **Static validation over runtime failure**
4. **Transpilation over interpretation**
5. **Tooling and visualization friendliness**

If a feature violates any of these principles, it must not be added.

---

## 4. Variable Declarations and Assignment (Locked)

### Declaration + Inference (`:=`)

```python
a := 2
name := "fin"
cmd := <grep "abc" file.txt>
(out, err) := run(platform=bash, cmd=cmd)
```

Rules:

* `:=` always **declares** new variable(s)
* Types are inferred from the RHS
* Redeclaration in the same scope is an error

---

### Assignment (`=`)

```python
a = a + 1
```

Rules:

* Variable must already exist
* Type must match
* Assignment never declares

---

## 5. Type System

Fin v2 is **statically typed**.

### Built-in Types

* `int`
* `float`
* `bool`
* `str`
* `list[T]`
* `map[K, V]`
* `command`
* `error`

### Type Rules

* All expressions are type-checked
* No implicit type coercion
* Commands and errors are first-class types

---

## 6. Commands as First-Class Values

Shell commands are **typed values**, not strings.

```python
cmd := <grep "abc" file.txt>
```

Rules:

* `<...>` produces a value of type `command`
* Commands are immutable
* Commands are backend-agnostic until lowered
* Commands cannot be concatenated or executed implicitly

---

## 7. `run()` — Explicit Effect Boundary

`run()` is the **only mechanism** that executes shell commands.

### Canonical Usage

```python
(out, err) := run(platform=bash, cmd=cmd)
```

Rules:

* `platform` is explicit
* Return values are typed
* `err` is non-nil on non-zero exit
* `run()` never executes implicitly

Go-like error handling is the intended model.

---

## 8. Control Flow

```python
if x > 0:
    run <echo "positive">
else:
    run <echo "zero">

for i in range(0, 3):
    run <echo i>

while ready:
    run <sleep 1>

```

Control flow is structural and statically analyzable.

---

## 9. Go-Style Types and Methods (No Classes)

Fin v2 supports **Go-style struct-like types with methods**.

### Function Definition

```python
def add(a: int, b: int) -> int:
    return a + b
```

### Type Definition

```python
type File:
    path: str
    size: int
```

### Method Definition

```python
def (f: File) is_large() -> bool:
    return f.size > 1024
```

Rules:

* No classes
* No inheritance
* No dynamic dispatch
* Methods are static functions with receiver sugar

Method calls are resolved at compile time.

---

## 10. Built-in and Standard Library (Locked)

### Built-in

Only **one** builtin exists:

```python
len(x)
```

Supported for:

* `str`
* `list[T]`
* `map[K, V]`

---

### Standard Library (`std`)

```python
import std
```

#### Numeric Aggregates

```python
std.max(nums)
std.min(nums)
```

#### Functional Operations (Pure Only)

```python
std.map(list, fn(x: T) -> U: expr)
std.filter(list, fn(x: T) -> bool: expr)
std.reduce(list, init, fn(acc: U, x: T) -> U: expr)
```

Rules:

* Functions must be pure
* No `run()` inside lambdas
* No mutation inside lambdas
* Lowered to loops in IR

---

## 11. What Is Explicitly Not Allowed

* Universal methods (e.g. `x.len()`)
* Method chaining on collections
* Built-in grep/find/fs APIs
* Side effects inside functional ops
* Implicit shell execution

Shell utilities must be invoked via `run()`.

---

## 12. Compiler Pipeline

```
.fin source
   ↓
Lexer
   ↓
Parser (AST)
   ↓
Semantic Analyzer (types & scopes)
   ↓
IR (lowered, backend-neutral)
   ↓
Backend Generator
   ↓
bash | fish | bat | ps1
```

Each stage must:

* Be deterministic
* Never panic
* Emit structured diagnostics

---

## 13. Visualization (First-Class Feature)

```bash
fin visualize ./main.fin
```

The visualizer must show:

* AST structure
* Control-flow graph
* Variable lifetimes
* Command values and execution points
* Backend lowering steps

Visualization reads compiler IR and never executes code.

---

## 14. Repository Structure

```
fin/
 ├─ cmd/fin/              # CLI
 ├─ internal/
 │   ├─ token/
 │   ├─ lexer/
 │   ├─ parser/
 │   ├─ ast/
 │   ├─ sema/
 │   ├─ ir/
 │   ├─ generator/
 │   ├─ visualize/
 │   └─ diagnostics/
 └─ AGENTS.md
```

---

## 15. Testing Requirements

All changes must include tests.

Required coverage:

* Lexer correctness
* Parser precedence & recovery
* Type checking
* IR lowering
* Backend golden outputs
* Visualization snapshots

Golden tests define backend behavior.

---

## 16. Versioning & Compatibility

* Grammar or type system changes → major version
* Backend additions → minor version
* Bug fixes → patch version

Breaking changes must be explicit and documented.

---

## 17. Agent Responsibilities

All agents must:

* Follow this document strictly
* Avoid scope creep
* Prefer clarity over cleverness
* Keep output deterministic
* Add tests for all behavior
* Escalate if a change risks turning Fin into a general-purpose language

---

## 18. Final Positioning (Locked)

> **Fin v2 is a statically typed, Pythonic automation language that transpiles to multiple shell backends with explicit execution and built-in visualization.**

This sentence defines the project. Do not dilute it.
