# WSJ TODO List

## High Priority - MVP Blockers

### Control Flow Statements (Required for mvp_001.wsj)
- [x] **If statements**: `if (condition) { ... } else { ... }` - COMPLETED
- [x] **For loops**: `for (init; condition; update) { ... }` - COMPLETED
- [x] **Block statements**: `{ statement1; statement2; }` - COMPLETED
- [x] **Continue and Break statements**: `continue;` and `break;` - COMPLETED
- [x] **Remove dead code**: Delete unused rules CallExpr, IndexExpr, MemberExpr (lines 199-228)
- [x] **Fix unreachable code**: Remove dead return statement on line 261 in null literal case  
- [x] **Add systematic whitespace handling** between tokens (top priority when implementing comments)

## Medium Priority

- [x] **Add boolean operators**: Implement `&&` (logical AND) and `||` (logical OR) operators for compound conditions - COMPLETED
- [ ] **Make operator precedence consistent**: Update Comparison rule to handle left-associative sequences like Term and Factor
- [x] **Standardize position tracking**: Use consistent approach for c.pos vs expression positions - COMPLETED
- [ ] **Improve parser error messages**: Replace cryptic pigeon-generated errors with user-friendly messages (implement after interpreter is complete)

## Low Priority

- [ ] **Template literals**: `` `Error: ${variable}` `` with expression interpolation (removed from mvp_001.wsj)
- [ ] **Improve string literal parsing** for proper escape sequence handling (will address with template literals)
- [ ] **Enhance reserved keyword handling**: Prevent keywords in broader contexts beyond Ident rule
- [x] **Remove or implement generalization comment** on line 192 for CallExpr callee

### MVP Script Issues to Fix
- [x] **Fix inconsistent index syntax in mvp_001.wsj**: Line 33 uses `[row, col]` but should use `[row][col]`

## Notes

- Helper functions are properly located in `parser/helpers.go` (no action needed)
- Number parsing is appropriate for Worldographer files (integers and floats only)
- String literal improvements will be addressed when implementing template literals
- Whitespace handling is top priority for future comment implementation
