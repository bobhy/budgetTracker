package config

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// helper to manage environment isolation
func withIsolatedEnv(t *testing.T, fn func(tempDir string)) {
	t.Helper()

	// 1. Reset Global Config
	Current = nil

	// 2. Mock Flags
	oldCommandLine := flag.CommandLine
	// Reset the FlagSet completely
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// 3. Mock Args/Env
	oldArgs := os.Args
	oldEnv := os.Environ()
	os.Clearenv() // Start clean

	// 4. Create Temp Dir
	tempDir, err := os.MkdirTemp("", "config_test_isolated")
	assert.NoError(t, err)

	// Set essential path vars
	os.Setenv("HOME", tempDir)
	os.Setenv("XDG_CONFIG_HOME", filepath.Join(tempDir, "config"))
	os.Setenv("XDG_DATA_HOME", filepath.Join(tempDir, "data"))

	defer func() {
		// Restore
		flag.CommandLine = oldCommandLine
		os.Args = oldArgs
		os.Clearenv()
		for _, e := range oldEnv {
			pair := splitEnv(e)
			os.Setenv(pair[0], pair[1])
		}
		os.RemoveAll(tempDir)
	}()

	fn(tempDir)
}

func splitEnv(s string) []string {
	for i := 0; i < len(s); i++ {
		if s[i] == '=' {
			return []string{s[:i], s[i+1:]}
		}
	}
	return []string{s, ""}
}

func TestInit_TableDriven(t *testing.T) {
	tests := []struct {
		name              string
		args              []string
		env               map[string]string
		configFileContent string
		configFilePath    string // optional, relative to tempDir

		expectedDb     string
		expectedImport string
		expectedArgs   []string
	}{
		{
			name:           "Defaults",
			args:           []string{"cmd"},
			expectedDb:     "/data/budgetTracker/budget.db", // "data" is from XDG_DATA_HOME mock
			expectedImport: "/data/budgetTracker/importHistory",
			expectedArgs:   []string{},
		},
		{
			name: "Environment Override",
			args: []string{"cmd"},
			env: map[string]string{
				"budgetTracker_database":     "/env/db.db",
				"budgetTracker_importFolder": "/env/import",
			},
			expectedDb:     "/env/db.db",
			expectedImport: "/env/import",
			expectedArgs:   []string{},
		},
		{
			name:           "Flag Override",
			args:           []string{"cmd", "-database", "/flag/db.db", "-importFolder", "/flag/import"},
			expectedDb:     "/flag/db.db",
			expectedImport: "/flag/import",
			expectedArgs:   []string{},
		},
		{
			name: "Precedence Flag > Env",
			args: []string{"cmd", "-database", "/flag/db.db"},
			env: map[string]string{
				"budgetTracker_database": "/env/db.db",
			},
			expectedDb:     "/flag/db.db", // Flag wins
			expectedImport: "/data/budgetTracker/importHistory",
			expectedArgs:   []string{},
		},
		{
			name: "Precedence Env > Config File",
			args: []string{"cmd"},
			env: map[string]string{
				"budgetTracker_database": "/env/db.db",
			},
			configFileContent: `database = "/file/db.db"`,
			configFilePath:    "config/budgetTracker/config.toml", // Default config path
			expectedDb:        "/env/db.db",                       // Env wins
			expectedImport:    "/data/budgetTracker/importHistory",
			expectedArgs:      []string{},
		},
		{
			name: "Config File Loading (Default Path)",
			args: []string{"cmd"},
			configFileContent: `
database = "/file/db.db"
importFolder = "/file/import"
`,
			configFilePath: "config/budgetTracker/config.toml",
			expectedDb:     "/file/db.db",
			expectedImport: "/file/import",
			expectedArgs:   []string{},
		},
		{
			name: "Config File via Flag",
			args: []string{"cmd", "-config", "custom_config.toml"},
			configFileContent: `
database = "/custom/db.db"
`,
			configFilePath: "custom_config.toml", // relative to tempDir root
			expectedDb:     "/custom/db.db",
			expectedImport: "/data/budgetTracker/importHistory",
			expectedArgs:   []string{},
		},
		{
			name:         "Non-Flag Arguments",
			args:         []string{"cmd", "-database", "/db", "arg1", "arg2"},
			expectedDb:   "/db",
			expectedArgs: []string{"arg1", "arg2"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			withIsolatedEnv(t, func(tempDir string) {
				// Setup Env
				for k, v := range tc.env {
					os.Setenv(k, v)
				}

				// Setup Config File
				if tc.configFileContent != "" {
					fullPath := filepath.Join(tempDir, tc.configFilePath)
					err := os.MkdirAll(filepath.Dir(fullPath), 0755)
					assert.NoError(t, err)
					err = os.WriteFile(fullPath, []byte(tc.configFileContent), 0644)
					assert.NoError(t, err)

					// If using -config flag with relative path, we might need to adjust or ensure Cwd matches?
					// But our test code does not chdir. It uses absolute paths if absolute.
					// The "Flag" test case passes "custom_config.toml"
					// config.go logic doesn't join with Cwd, it assumes os.Stat works.
					// So if we pass a filename, it looks in Cwd.
					// We are not changing Cwd in withIsolatedEnv.
					// So for "Config File via Flag" to work with just "custom_config.toml", we should probably write it to actual Cwd?
					// OR, easier: update the args in the test case to be absolute path to temp file.
					// Let's dynamically update args if we detect we need the temp path.
				}

				// Fixup args if they reference our config file and we want it to point to tempDir location
				argsToUse := make([]string, len(tc.args))
				copy(argsToUse, tc.args)
				for i, arg := range argsToUse {
					if arg == "custom_config.toml" {
						argsToUse[i] = filepath.Join(tempDir, tc.configFilePath)
					}
				}
				os.Args = argsToUse

				// Run Init
				remainingArgs := Init()

				// Assert
				// Handle dynamic defaults in expectations if needed (e.g. they contain "/data/budgetTracker")
				// We used "/data" as the mocked XDG_DATA_HOME.
				// However, the test expectations above hardcoded "/data/..." which matches our mock data home.
				// So we should be fine.

				// But wait, the mock uses a random temp dir for "data".
				// "filepath.Join(tempDir, "data")" will result in something like "/tmp/config_test_isolated123/data".
				// My test table expectations utilize literal "/data/...".
				// I need to replace placeholders or resolve expectations dynamically.

				// Strategy: Resolve expectations relative to tempDir.
				mockDataHome := filepath.Join(tempDir, "data")

				resolve := func(path string) string {
					if path == "" {
						return ""
					}
					if len(path) > 5 && path[:6] == "/data/" {
						return filepath.Join(mockDataHome, path[6:])
					}
					return path
				}

				assert.Equal(t, resolve(tc.expectedDb), Current.DatabasePath, "DatabasePath mismatch")

				// Only check ImportPath if we have an expectation for it (or default checking)
				// My table has expectedImport defaults for most cases.
				expectedImportResolved := resolve(tc.expectedImport)
				if expectedImportResolved == "" {
					// Fallback for cases I didn't specify in table, though above I did specify defaults.
					// Let's assume table is complete or empty string means "don't care" (but here we care about defaults).
					// Actually, let's just use the resolve logic.
					expectedImportResolved = filepath.Join(mockDataHome, "budgetTracker", "importHistory")
				}

				if tc.expectedImport != "" {
					assert.Equal(t, resolve(tc.expectedImport), Current.ImportPath, "ImportPath mismatch")
				} else {
					// If not specified in test case, assume we expect default (or just don't check?)
					// Better to check default.
					assert.Equal(t, filepath.Join(mockDataHome, "budgetTracker", "importHistory"), Current.ImportPath, "ImportPath mismatch (default)")
				}

				assert.Equal(t, tc.expectedArgs, remainingArgs, "Args mismatch")
			})
		})
	}
}

func TestToJSON(t *testing.T) {
	config := &Config{
		DatabasePath: "/db",
		ImportPath:   "/import",
	}

	jsonStr, err := config.ToJSON()
	assert.NoError(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &parsed)
	assert.NoError(t, err)

	assert.Equal(t, "/db", parsed["database"])
	assert.Equal(t, "/import", parsed["importFolder"])
}
