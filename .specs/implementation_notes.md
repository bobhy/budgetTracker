# Implementation Notes

## Testing Requirements

All new code must be accompanied by unit tests.

### Test Structure
Tests should follow the **Table Driven Tests** pattern (idiomatic Go).
- Define a slice of struct test cases.
- Iterate over the slice and run each case using `t.Run(tc.name, ...)`.
- **Convention**: The first field in the test case struct should be `name` (string), describing the subcase. This name will be used in `t.Run`.

### Assertion Library
Use the `github.com/stretchr/testify/assert` package for assertions.
- Use `assert.Equal(t, expected, actual)` for equality checks.
- Use `assert.NoError(t, err)` for error checking.
- Do not use manual `if err != nil { t.Errorf(...) }` unless strictly necessary for custom logic.

### Example
```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"Empty input", "", "default"},
        {"Valid input", "foo", "bar"},
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            result := MyFunction(tc.input)
            assert.Equal(t, tc.expected, result)
        })
    }
}
```

## Data Safety (Agent Testing)

The agent must **NEVER** run the application against default data locations (which likely contain the user's production data).

**Rules:**
1. **Always isolate data**: When invoking the app (via `wails dev`, `go run`, or compiled binary), explicit paths for the database and configuration must be provided.
2. **Use Testing Locations**: Use a dedicated `agent_testing/` directory in the project root or a transient system temp directory.
    - Example: `agent_testing/config.toml`
    - Example: `agent_testing/budget.db`
3. **Mechanism**: Override defaults using Command Line Flags or Environment Variables.

**Examples:**
- **Command Line**:
    ```bash
    wails dev --appargs "-database agent_testing/budget.db -importFolder agent_testing/imports"
    ```
- **Environment**:
    ```bash
    budgetTracker_database="agent_testing/budget.db" budgetTracker_importFolder="agent_testing/imports" wails dev
    ```
