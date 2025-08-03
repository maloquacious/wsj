# 📜 WSJ Specification: Variable Declarations with `let`

WSJ supports two distinct forms of variable declarations using `let`:

---

## Form A: **Single Variable Declarations (Preferred for Multiple Assignments)**

### Syntax

```js
let x = expr1, y = expr2, z = expr3;
```

* Each variable is bound to its own **independent expression**.
* All declarations are evaluated **left to right**.

### Examples

```js
let a = 1, b = 2;
let message = "hello", count = getCount();
```

### Semantics

* Evaluate each RHS expression individually.
* Bind each LHS variable in the current scope.
* The declarations are equivalent to multiple `let` statements.

---

## Form B: **Multiple Variables Assigned from Multi-Return Function**

### Syntax

```js
let x, y = functionCall();
```

* Only valid when the RHS is an expression that returns **multiple values** (e.g., a built-in or user-defined function).
* RHS **must return exactly as many values as variables declared**.

### Examples

```js
let map, err = wxx.Load("map.wxx");
let x, y, z = getCoordinates();
```

### Semantics

* Evaluate the RHS once.
* If the number of returned values does **not match** the number of LHS variables:

  * **Runtime error** with positional context.
* Variables are declared in current scope.

---

## ❌ Disallowed Form: Mixed Multi-Variable with Multiple Expressions

```js
let x, y = 5, 4;  // ❌ Error
```

* This form is **not allowed** in WSJ.
* Use **Form A** instead:

  ```js
  let x = 5, y = 4; // ✅
  ```

---

## Error Message Example

If a user writes:

```js
let x, y = 5, 4;
```

WSJ should report:

```
Error: Multiple variables may only be declared with a single expression (such as a function call) on the right-hand side. Use `let x = 5, y = 4;` instead.
```

---

## 🛠 Grammar Specification (PEG)

```peg
LetStmt <- "let" _ letList:LetList _ ";" {
    return letList, nil
}

LetList <- MultiReturnLet / SingleDeclList

// Try multi-return first (contains comma before =)
MultiReturnLet <- first:Ident rest:(_ "," _ Ident)+ _ "=" _ expr:Expression {
    identifiers := []*ast.Ident{first.(*ast.Ident)}
    for _, item := range rest.([]interface{}) {
        pair := item.([]interface{})
        ident := pair[3].(*ast.Ident) // skip whitespace, comma, whitespace
        identifiers = append(identifiers, ident)
    }
    return &ast.MultiReturnLetStmt{
        Names: identifiers,
        Value: expr.(ast.Expr),
        Pos: ast.Pos{Line: c.pos.line, Column: c.pos.col},
    }, nil
}

SingleDeclList <- first:SingleDecl rest:(_ "," _ SingleDecl)* {
    decls := []*ast.LetStmt{first.(*ast.LetStmt)}
    for _, item := range rest.([]interface{}) {
        pair := item.([]interface{})
        decl := pair[3].(*ast.LetStmt) // skip whitespace, comma, whitespace
        decls = append(decls, decl)
    }
    return &ast.SingleDeclListStmt{
        Declarations: decls,
        Pos: ast.Pos{Line: c.pos.line, Column: c.pos.col},
    }, nil
}

SingleDecl <- name:Ident _ "=" _ value:Expression {
    return &ast.LetStmt{
        Name: name.(*ast.Ident),
        Value: value.(ast.Expr),
        Pos: ast.Pos{Line: c.pos.line, Column: c.pos.col},
    }, nil
}
```

---

## 🧩 Parser Notes

| Case                     | How to Parse                                              |
| ------------------------ | --------------------------------------------------------- |
| `let x = expr, y = expr` | Parse as `SingleDeclList`                                 |
| `let x, y = expr`        | Parse as `MultiReturnLet`                                 |
| Disambiguation           | If any `,` appears **before** `=`, it's `MultiReturnLet`. |

---

## 🧮 VM Execution Summary

| Form           | Evaluation                                          |
| -------------- | --------------------------------------------------- |
| SingleDeclList | Evaluate each expr, assign to respective variable.  |
| MultiReturnLet | Evaluate RHS, check arity == number of identifiers. |
|                | Assign result\[i] to identifier\[i].                |

---

## 🧭 Design Goals Satisfied

* ✅ **Clarity**: Only one valid way to declare multiple variables per style.
* ✅ **Simplicity**: No tuple unpacking, no destructuring, easy to learn.
* ✅ **Go-inspired**: Consistent with Go's multi-return pattern.
* ✅ **No ambiguity**: Easy to parse, execute, and explain.

---

# Notes: Multiple Variable Assignment with `let`

WSJ supports **multiple variable assignment** using the `let` keyword, similar to Go. This allows a function to return multiple values, such as `(result, error)`, and for these to be assigned explicitly.

## Syntax

```js
let var1, var2, ..., varN = expression;
```

* Declares `var1` through `varN` in the current scope.
* Evaluates `expression` once and expects it to return **exactly N values**.
* Assigns value `i` to `vari`.

---

## Examples

```js
let map, err = wxx.Load("maps/forloriath.wxx");
if (err != null) {
  print(`Error: ${err}`);
}

let x, y, z = getCoordinates();
```

---

## Semantics

1. **Evaluation Order**:

   * The RHS `expression` is evaluated first.
   * The result **must be a tuple/list of values**, exactly matching the number of identifiers.

2. **Binding**:

   * Variables are **declared** and **assigned** in the current scope.
   * Rebinding (assigning again with `let`) is **not allowed**; use `x = value` to reassign.

3. **Arity Mismatch**:

   * If `expression` returns **fewer or more** than N values:

     * The VM **throws an error** at runtime with positional context.
     * Example: `let x, y = foo();` but `foo()` returns 1 value → error.

4. **Scoping**:

   * Variables declared with `let` are scoped to the **current block**.
   * Shadowing is allowed (e.g., redeclaring `x` inside a function).

---

## Runtime Implementation Notes (VM)

* Evaluate RHS → result slice/tuple.
* Verify `len(result) == len(IdentifierList)` → else error.
* Assign `result[i]` → identifier\[i] in current scope.

---

## Design Rationale

* ✔ One clear way to assign multiple values.
* ✔ Simpler parser and runtime logic.
* ✔ Easier for error reporting (no pattern matching).
* ✔ Avoids confusion with JavaScript destructuring.
* ✔ Compatible with Go-like function signatures.

---

## Error Reporting

**Compile-Time (optional)**:

* Enforce unique identifiers in the same `let` declaration.

**Runtime**:

* If arity mismatch, return error:

  ```
  Line 12: function wxx.Load returned 1 value(s), but 2 variables expected
  ```

---

## Future-Proofing (Explicit Non-Goals)

* **No destructuring**:

  * `let [x, y] = foo();` → invalid syntax.
* **No tuple unpacking from array**:

  * `let x, y = [1,2];` → invalid unless function returns multiple values.

---

# 🧩 Spec Clarification: Function Return Arity

## Must WSJ Functions Always Return the Same Number of Values?

### ✅ **Yes** — WSJ functions (including built-ins) must have a **fixed return arity**.

---

## Why Require Fixed Arity?

| Benefit                        | Description                                                                   |
| ------------------------------ | ----------------------------------------------------------------------------- |
| ✅ **Simple parser/VM**         | No need to inspect values dynamically or adjust assignment behavior.          |
| ✅ **Clear error messages**     | Easy to detect and report mismatch between `let a, b = ...` and return arity. |
| ✅ **Predictable for users**    | No surprise: you always know what a function returns.                         |
| ✅ **Matches Go semantics**     | Reinforces Go-like model: function signature defines return arity.            |
| ✅ **Avoids runtime ambiguity** | Eliminates “how many values did this function return?” problem.               |

---

## Language Rule

> Every function in WSJ **must declare and always return** the **same number of values** each time it is called.

---

## Function Categories

| Type         | Must Return            | Example                     |
| ------------ | ---------------------- | --------------------------- |
| Built-in     | Fixed number of values | `wxx.Load → (Map, Error)`   |
| User-defined | Fixed number of values | `fn foo() { return 1, 2; }` |
| Void         | Zero values            | `fn bar() { ... }`          |

---

## Enforcement

| Phase             | Enforcement Strategy                                                |
| ----------------- | ------------------------------------------------------------------- |
| Parse Time        | Require all `return` statements in a function to match arity.       |
| Compile Time (VM) | Mark function's arity in function object (e.g., `ReturnCount = 2`). |
| Runtime           | VM enforces arity in `let x, y = f()` and errors if mismatch.       |

---

## Example Errors

```js
fn foo() { return 1; } // Valid if foo always returns 1 value

let x, y = foo(); // Error: foo returns 1 value, but 2 were expected
```

---

## Optional: `return null` to fill out values?

Disallowed. Let’s **require explicit return of all values**:

```js
fn foo() {
  return 1, null; // required if 2 values expected
}
```

---

## Future Possibilities (Optional, but Not Now)

* Tuples or variadic return: `return tuple;` (Post-MVP).
* Default fill: missing values filled with `null`. (Discouraged for simplicity).

---

## Final Recommendation

**Functions must always return the same number of values**. It is **simpler**, **easier to document**, and **avoids ambiguity**. It’s the **right design choice for WSJ**, especially if you want Go-like clarity and **excellent error messages**.

---

Correct — under the current design philosophy, **`let x;` should be invalid**.

Let’s make this explicit and reason it through carefully.

---

# 🚫 `let x;` — **Not Allowed in WSJ**

## Design Rule

> **All variable declarations must include an initializer.**
> Declaration without initialization (e.g., `let x;`) is **not valid** in WSJ.

---

## 📜 Rationale

| Reason                      | Explanation                                                                                                  |
| --------------------------- | ------------------------------------------------------------------------------------------------------------ |
| ✅ **Clarity**               | Every variable has a **defined value** at declaration. No uninitialized state.                               |
| ✅ **Simpler runtime model** | No need for `undefined`/`nil` as a default placeholder.                                                      |
| ✅ **Error prevention**      | Prevents bugs from accidental use of uninitialized variables.                                                |
| ✅ **Go-inspired design**    | Aligns with Go's philosophy: every variable has a value (though Go allows zero value init, this is tighter). |
| ✅ **Smaller parser/VM**     | No special case for uninitialized storage.                                                                   |
| ❌ **Not JS-like**           | In JS, `let x;` is valid — but we explicitly **do not want this valley effect**.                             |

---

## 🛠 Specification Update

### Valid:

```js
let x = 0;
let msg = "hello";
let a = 1, b = 2;
let map, err = wxx.Load("map.wxx");
```

### Invalid:

```js
let x;        // ❌ Error
let a, b;     // ❌ Error
let x = 1, y; // ❌ Error
```

---

## 📢 Error Message

When a user writes:

```js
let x;
```

Return:

```
Error: Variable declarations must include an initializer. Use `let x = value;`.
```

---

## 🧩 Grammar Reinforcement (PEG)

The PEG grammar shown above **does not allow** `Identifier` without `= Expression`. Both `SingleDecl` and `MultiReturnLet` require an assignment expression, enforcing the "all variables must be initialized" rule.

---

## 🏆 Final Recommendation

* ✅ Require all `let` declarations to include an initializer.
* ✅ Disallow `let x;` to keep WSJ **simple, explicit, and reliable**.

---

# ❌ Why WSJ Does Not Support `let x;` (Uninitialized Declarations)

## Language Design Principle: **Every Variable Must Be Initialized**

In WSJ, all variable declarations **must include an initializer**. The following is **not allowed**:

```js
let x;        // ❌ Invalid: no initializer
let a, b;     // ❌ Invalid: no initializer for either variable
let y = 5, z; // ❌ Invalid: missing initializer for z
```

---

## ✅ Correct Declarations

WSJ requires **explicit initialization** of each variable:

```js
let x = 0;                      // ✅ Valid
let name = "Forloriath";        // ✅ Valid
let row = 10, col = 20;         // ✅ Valid
let map, err = wxx.Load(path);  // ✅ Valid (from multi-return function)
```

---

## 🎯 Design Goals Behind This Rule

| Goal                           | Why We Enforce It                                                     |
| ------------------------------ | --------------------------------------------------------------------- |
| ✅ **Clarity**                  | All variables have a defined value from the moment they are declared. |
| ✅ **Error Prevention**         | Eliminates bugs from using variables before assignment.               |
| ✅ **Simplicity**               | No need for a special “uninitialized” or “undefined” value.           |
| ✅ **No Implicit Defaults**     | Avoids hidden behavior or assumptions about default values.           |
| ✅ **Better Error Messages**    | Helps catch mistakes early during parsing or execution.               |
| ✅ **Smaller Language Surface** | One way to declare variables — easy to learn and document.            |

---

## 🧩 Comparison with Other Languages

| Language       | Behavior of `let x;`                         |
| -------------- | -------------------------------------------- |
| **JavaScript** | Valid; x is `undefined`.                     |
| **Go**         | Invalid; must use `var x T` or `x := value`. |
| **WSJ**        | ❌ Invalid — must use `let x = value;`.       |

WSJ **does not adopt JavaScript’s `undefined` behavior**. Instead, it chooses **explicitness over flexibility**, leading to **fewer bugs and better code readability**.

---

## 🛠 Implementation Simplicity

* No need for **special runtime values** for “not yet initialized.”
* No **branching in execution** for “is this variable initialized?”
* No **nullable binding state** in the interpreter or VM.

---

## 📢 Developer Feedback Example

> “Why can’t I write `let x;`?”
>
> Because WSJ requires you to always provide a value. Use `let x = null;` if you want to declare a variable explicitly intended to hold no meaningful value yet.

---

## 💡 Tip for Users

If you truly need a variable to start without a real value, you can initialize it to `null`:

```js
let x = null;
```

This communicates your intent clearly and keeps the language semantics simple and robust.

---

# 🚧 Implementation TODOs (Not This Sprint)

## Status: **DEFERRED** 
This enhanced `let` statement specification will **not be implemented this sprint**. The following issues must be resolved before implementation can begin.

---

## 🚩 Blocking Issues

### 1. Missing AST Node Types
The PEG grammar references AST nodes that don't exist in the current codebase:
- [ ] Add `ast.MultiReturnLetStmt` struct
- [ ] Add `ast.SingleDeclListStmt` struct  
- [ ] Update AST interfaces to support new node types

### 2. Undefined Runtime Behavior
Critical VM implementation details are missing:
- [ ] Define how multi-return values are represented in VM
- [ ] Specify what constitutes a "tuple/list of values"
- [ ] Design function arity tracking mechanism
- [ ] Define built-in function metadata system

### 3. Incomplete Function System
Multi-return requires function definition support:
- [ ] Define function declaration syntax (`fn name() { ... }`)
- [ ] Implement function arity validation at parse time
- [ ] Create built-in function registration system
- [ ] Design function signature metadata storage

### 4. Error Message Standards
Inconsistent error reporting throughout spec:
- [ ] Standardize error message format (e.g., `"Error at line X: message"`)
- [ ] Define error reporting patterns for parser vs runtime
- [ ] Create comprehensive error code system

### 5. VM Architecture Gaps
Current VM may not support required features:
- [ ] Multi-value return mechanism
- [ ] Runtime arity checking
- [ ] Function metadata lookup
- [ ] Tuple/array value handling

---

## 📋 Alternative MVP Approach

**Recommendation**: Implement simple `let x = value` only for MVP:

```js
// MVP Implementation
let x = 5;
let name = "test";
let a = 1, b = 2;  // Multiple single declarations
```

**Deferred for v1.1+**:
```js
// Post-MVP Features
let map, err = wxx.Load("file.wxx");  // Multi-return
```

This allows shipping a working interpreter while deferring the complex multi-return system.

---

## 🎯 Sprint Decision

**Current Sprint**: Focus on basic `let x = value` syntax to unblock interpreter development.

**Future Sprint**: Implement multi-return after VM architecture supports required features.

---
