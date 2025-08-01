# Grammar Improvement TODOs

## High Priority

- [ ] **Remove dead code**: Delete unused rules CallExpr, IndexExpr, MemberExpr (lines 199-228)
- [ ] **Fix unreachable code**: Remove dead return statement on line 261 in null literal case  
- [ ] **Add systematic whitespace handling** between tokens (top priority when implementing comments)

## Medium Priority

- [ ] **Make operator precedence consistent**: Update Comparison rule to handle left-associative sequences like Term and Factor
- [ ] **Standardize position tracking**: Use consistent approach for c.pos vs expression positions

## Low Priority

- [ ] **Improve string literal parsing** for proper escape sequence handling (will address with template literals)
- [ ] **Enhance reserved keyword handling**: Prevent keywords in broader contexts beyond Ident rule
- [ ] **Remove or implement generalization comment** on line 192 for CallExpr callee

## Notes

- Helper functions are properly located in `parser/helpers.go` (no action needed)
- Number parsing is appropriate for Worldographer files (integers and floats only)
- String literal improvements will be addressed when implementing template literals
- Whitespace handling is top priority for future comment implementation
