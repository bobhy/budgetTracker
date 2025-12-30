---
trigger: manual
---

# Configuration System

How to configure runtime behavior in this project.

The project supports the [XDG Base Directory specification](https://linuxvox.com/blog/linux-xdg/).  It stores:
- a configuration file in `$XDG_CONFIG_HOME/<project>/config.toml`
- a database and other history files in `$XDG_DATA_HOME/<project>/`

Where `<project>` is the placeholder for the project name (as defined in `wails.json`), currently "budgetTracker".

## Config sample
A sample config file, showing all defined configuration items and their defaults, is in file://./sampleConfig.toml.

## Specifying configuration items

Configuration items can be specified in the following places:

1. on command line arguments.  
The name of the flag is the same as the field name in the configuration file.
1. in environment variables  
The name of the environment variable is `<project>_<field>`, where `<field>` is the field name in the configuration file.
1. in a configuration file

The configuration file settings can be overridden by environment variables, and both can be overridden by a command line flag.

# Package "config"
This package contains the types that define the configuration for the rest of the application.  

It uses package "burntsushi/toml" for toml handling.  It uses https://pkg.go.dev/flag to parse command line flags.

It exports a global, static Go struct `config.Config` and methods to serialize Config as JSON for use by the front end (which might not be written in Go).

This package also contains an `config.Init' function which:
1. Starts with a Config struct initialized with defaults.  
2. Determines name of project (placeholder `<project>`) below.
3. Looks for a command line flag "config" or the corresponding environment variable `<project>_config`, and updates the path to the config file.
4. Reads in the config file from the path determined above.
5. Scans all environment variables starting with `<project>_` and updates Config field accordingly.
6. Scans the flags in command line, overriding corresponding Config field.
7. Returns the command line arguments that were not flags.

 The app invokes `config.Init` to get a finalized Config struct and filtered list of command line arguments.  Current project doesn't this list yet.
