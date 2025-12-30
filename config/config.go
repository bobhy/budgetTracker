package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var Current *Config

// Config holds the application configuration
type Config struct {
	DatabasePath string `toml:"database" json:"database"`
	ImportPath   string `toml:"importFolder" json:"importFolder"`
}

// ToJSON returns the configuration as a JSON string
func (c *Config) ToJSON() (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Init initializes the configuration system.
// It resolves settings from Defaults -> Config File -> Env Vars -> Flags.
// Returns the non-flag command line arguments.
func Init() []string {
	// 1. Initialize Defaults
	Current = &Config{}
	
	// Determine XDG paths
	configHome, _ := os.UserConfigDir() // e.g., ~/.config
	if configHome == "" {
		configHome = filepath.Join(os.Getenv("HOME"), ".config")
	}
	
	dataHome, _ := os.UserHomeDir() // Used for .local/share default
	// On Linux, data home is usually ~/.local/share, but go doesn't have a direct UserDataDir.
	// We'll stick to XDG spec default: $HOME/.local/share if XDG_DATA_HOME is not set.
	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome == "" {
		xdgDataHome = filepath.Join(dataHome, ".local", "share")
	}

	defaultConfigPath := filepath.Join(configHome, "budgetTracker", "config.toml")
	defaultDatabasePath := filepath.Join(xdgDataHome, "budgetTracker", "budget.db")
	defaultImportPath := filepath.Join(xdgDataHome, "budgetTracker", "importHistory")

	Current.DatabasePath = defaultDatabasePath
	Current.ImportPath = defaultImportPath

	// 2. Define Flags (but define them locally variables to hold flag values)
	// We don't want to overwrite defaults with empty strings if flags aren't set.
	var flagConfig, flagDatabase, flagImport string

	flag.StringVar(&flagConfig, "config", "", "Path to configuration file")
	flag.StringVar(&flagDatabase, "database", "", "Path to database file")
	flag.StringVar(&flagImport, "importFolder", "", "Path to import history folder")

	flag.Parse()

	// 3. Load Config File
	// Priority: Flag > Env > Default
	configPathToUse := defaultConfigPath
	
	envConfig := os.Getenv("budgetTracker_config")
	if envConfig != "" {
		configPathToUse = envConfig
	}
	if flagConfig != "" {
		configPathToUse = flagConfig
	}

	// Try to read config file
	if _, err := os.Stat(configPathToUse); err == nil {
		// Create a temp config to decode into
		var fileConfig Config
		if _, err := toml.DecodeFile(configPathToUse, &fileConfig); err == nil {
			// Merge: only update if value is present in file (though TOML decoding handles this naturally for filled fields)
			// Since we want file to override defaults:
			if fileConfig.DatabasePath != "" {
				Current.DatabasePath = fileConfig.DatabasePath
			}
			if fileConfig.ImportPath != "" {
				Current.ImportPath = fileConfig.ImportPath
			}
		} else {
			fmt.Printf("Warning: Failed to parse config file at %s: %v\n", configPathToUse, err)
		}
	}

	// 4. Apply Environment Variables (Override File/Defaults)
	if envDb := os.Getenv("budgetTracker_database"); envDb != "" {
		Current.DatabasePath = envDb
	}
	if envImport := os.Getenv("budgetTracker_importFolder"); envImport != "" {
		Current.ImportPath = envImport
	}

	// 5. Apply Flags (Override Env/File/Defaults)
	if flagDatabase != "" {
		Current.DatabasePath = flagDatabase
	}
	if flagImport != "" {
		Current.ImportPath = flagImport
	}

	// 6. Finalization: Ensure directories exist
	dbDir := filepath.Dir(Current.DatabasePath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		fmt.Printf("Warning: Failed to create database directory %s: %v\n", dbDir, err)
	}
	if err := os.MkdirAll(Current.ImportPath, 0755); err != nil {
		fmt.Printf("Warning: Failed to create import directory %s: %v\n", Current.ImportPath, err)
	}

	return flag.Args()
}
