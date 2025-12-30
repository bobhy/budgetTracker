# Configuration Implementation Plan

## Overview
Implement a robust configuration system using the `config` package as specified in `.specs/configuration_plan.md`. This system will resolve settings from Defaults, Config File, Environment Variables, and Command Line Flags (in that order of priority).

## Steps

### 1. Dependencies
- Add the TOML parser dependency:
  `go get github.com/BurntSushi/toml`

### 2. Create `config` Package
Create directory `config/` and file `config/config.go`.

#### Struct Definition
```go
package config

var Current *Config

type Config struct {
    DatabasePath  string `toml:"database" json:"database"`
    ImportPath    string `toml:"importFolder" json:"importFolder"`
}
```

#### Init Function `func Init() []string`
1. **Initialize Defaults**:
    - Determine XDG Base Directories (`os.UserConfigDir`, `os.UserHomeDir` for `.local/share`).
    - Defaults:
        - Config: `$XDG_CONFIG_HOME/budgetTracker/config.toml`
        - Database: `$XDG_DATA_HOME/budgetTracker/budget.db`
        - Import: `$XDG_DATA_HOME/budgetTracker/importHistory`
2. **Parse Flags**:
    - Use `flag` package.
    - Define flags: `config`, `database`, `importFolder`.
    - Parse `flag.Parse()`.
3. **Load Config File**:
    - Determine Config Path (Priority: Flag > Env `BUDGETTRACKER_CONFIG` > Default).
    - If file exists, decode TOML into a temp struct and update `Current` (or merge).
4. **Apply Environment Variables**:
    - Check `BUDGETTRACKER_DATABASE`, `BUDGETTRACKER_IMPORTFOLDER`.
    - Override `Current` values if Env var is present.
5. **Apply Flags**:
    - If flags were set (non-empty), override `Current` values.
6. **Finalization**:
    - Ensure data directories exist (mkdir -p).
    - Return `flag.Args()` (non-flag arguments).

### 3. Update `main.go`
- Import `budgetTracker/config` (or appropriate module path).
- Call `config.Init()` at the start of `main()`.
- Update `models.NewService(...)` to use `config.Current.DatabasePath`.

### 4. Verification
- Verify `sampleConfig.toml` keys match struct tags (`database`, `importFolder`).
- Test with:
    - No args (Defaults).
    - `-database` flag.
    - `BUDGETTRACKER_DATABASE` env var.
    - `config.toml` file.

## Questions/Clarifications
- Spec says "exports ... methods to serialize Config as JSON". I will add a `ToJSON()` method or ensure the struct is json-serializable (add `json` tags).
- Project name is "budgetTracker" (from `wails.json`).
