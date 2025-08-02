# WSJ Interpreter Design Specification

## Overview

The WSJ Interpreter is a tree-walking interpreter designed to execute WSJ programs for Worldographer map automation. It prioritizes simplicity and direct execution over performance optimization.

## Core Architecture

### Interpreter Structure

```go
type Interpreter struct {
    globals map[string]interface{} // Global symbol table
    output  io.Writer              // Output destination (stdout, files, etc.)
}
```

**Design Decisions:**
- No built-in world map state - users load maps as variables
- Global symbol table supports multiple loaded maps simultaneously
- No function call stack (functions not implemented in MVP)
- Clean separation between interpreter state and user data

### Symbol Table

The `globals` map serves as the primary symbol table, storing:

```go
// Example after user executes:
// a := load("region1.wxx")
// b := load("region2.wxx") 
// playerName := "Aragorn"
globals = map[string]interface{}{
    "a":          *wxx.Map,    // Loaded Worldographer map
    "b":          *wxx.Map,    // Second map for comparison/merging
    "playerName": "Aragorn",   // String variable
    "maxLevel":   20,          // Numeric variable
}
```

### Variable Resolution

```go
func (interp *Interpreter) getVariable(name string) (interface{}, error) {
    if val, exists := interp.globals[name]; exists {
        return val, nil
    }
    return nil, fmt.Errorf("undefined variable: %s", name)
}

func (interp *Interpreter) setVariable(name string, value interface{}) {
    interp.globals[name] = value
}
```

## Execution Model

### Script Execution
- Create fresh interpreter instance per script
- Execute entire Program AST
- Interpreter is discarded after completion
- Isolated execution environment

```go
// Script mode
func runScriptFile(filename string, debug bool) error {
    // Parse script to AST
    interp := interpreter.New()  // Fresh interpreter
    return interp.Execute(program)
}
```

### REPL Execution
- Single persistent interpreter instance
- Each input line/block parsed as Program
- Interpreter state maintained between commands
- Variables persist across executions

```go
// REPL mode
func runREPL(debug bool) error {
    interp := interpreter.New()  // Persistent interpreter
    for {
        // Parse user input to Program
        interp.Execute(program)  // Reuse same interpreter
    }
}
```

## MVP Implementation

### Core Evaluator

```go
func (interp *Interpreter) Execute(program *ast.Program) error {
    for _, stmt := range program.Statements {
        if err := interp.executeStatement(stmt); err != nil {
            return err
        }
    }
    return nil
}

func (interp *Interpreter) executeStatement(stmt ast.Statement) error {
    switch s := stmt.(type) {
    case *ast.AssignmentStatement:
        return interp.executeAssignment(s)
    case *ast.ExpressionStatement:
        _, err := interp.evaluateExpression(s.Expression)
        return err
    case *ast.IfStatement:
        return interp.executeIf(s)
    // ... other statement types
    }
}
```

### Expression Evaluation

```go
func (interp *Interpreter) evaluateExpression(expr ast.Expression) (interface{}, error) {
    switch e := expr.(type) {
    case *ast.IntegerLiteral:
        return e.Value, nil
    case *ast.StringLiteral:
        return e.Value, nil
    case *ast.Identifier:
        return interp.getVariable(e.Value)
    case *ast.CallExpression:
        return interp.executeCall(e)
    // ... other expression types
    }
}
```

### Built-in Functions

Worldographer-specific functions available to WSJ programs:

| Function | Purpose | Example |
|----------|---------|---------|
| `load(filename)` | Load .wxx map file | `map := load("dungeon.wxx")` |
| `save(map, filename)` | Save map to file | `save(map, "updated.wxx")` |
| `print(...)` | Output to console | `print("Processing hex", x, y)` |
| `getHex(map, x, y)` | Get hex at coordinates | `hex := getHex(map, 5, 10)` |
| `setHex(map, x, y, terrain)` | Set hex terrain | `setHex(map, 5, 10, "forest")` |
| `distance(x1, y1, x2, y2)` | Hex distance | `d := distance(0, 0, 3, 4)` |

### Error Handling

- Runtime errors include line numbers from AST
- Type checking at execution time
- Graceful error recovery in REPL mode
- Fatal errors exit scripts

## Integration Points

### Parser Integration
```go
// Current interface works perfectly
result, err := parser.Parse("", []byte(input))
program := result.(*ast.Program)
interp.Execute(program)
```

### REPL Integration
```go
// Persistent interpreter instance
interp, _ := interpreter.New()
// Each command block
runProgram(interp, input, debug)
```

### Dependency Management
- Import `github.com/maloquacious/wxx` for map file I/O
- Import `github.com/maloquacious/hexg` for coordinate calculations
- Standard library for basic operations

## Future Extensions

### Function Support (Post-MVP)
```go
type Interpreter struct {
    globals map[string]interface{}
    stack   []map[string]interface{} // Call stack for function scoping
    functions map[string]*ast.FunctionLiteral
}
```

### Advanced Features
- User-defined functions with local scopes
- Module system for reusable WSJ libraries
- Debugging hooks and breakpoints
- Performance optimizations (bytecode compilation)

## Testing Strategy

### Unit Tests
- Individual expression evaluation
- Statement execution
- Built-in function behavior
- Error condition handling

### Integration Tests
- Complete program execution
- REPL session simulation
- File I/O operations with sample .wxx files
- Multi-map operations

### Test Data
- Sample .wxx files in `testdata/`
- Expected output for various WSJ programs
- Error case validation
