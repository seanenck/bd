// Package main is the primary executable
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/seanenck/bd/internal/core"
	"github.com/seanenck/bd/internal/fetch"
)

const (
	check         = "check"
	upgrade       = "upgrade"
	configFileEnv = "BD_CONFIG_FILE"
	appFlag       = "applications"
	disableFlag   = "disable"
)

func defaultConfig() string {
	home := os.Getenv("HOME")
	return filepath.Join(home, ".config", "bd", "config.yaml")
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func helpLine(prefix bool, flag, text string) {
	f := flag
	if prefix {
		f = fmt.Sprintf("-%s", f)
	}
	fmt.Printf("  %-15s %s\n", f, text)
}

func help(msg string) error {
	if msg != "" {
		fmt.Printf("%s\n\n", msg)
	}
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	exe = filepath.Base(exe)
	fmt.Printf("%s\n", exe)
	helpLine(false, check, "check for updates")
	helpLine(false, upgrade, "upgrade packages")
	helpLine(true, appFlag, "specify a subset of packages (comma delimiter)")
	helpLine(true, disableFlag, "disable applications (comma delimiter)")
	fmt.Println()
	fmt.Printf("configuration file: %s\n", defaultConfig())
	fmt.Printf("  (override using %s)\n", configFileEnv)
	fmt.Println()
	fmt.Printf("to handle github rate limiting, specify a token in configuration or via env\n")
	fmt.Printf("  %s (directly or as reference to a file)\n", strings.Join(fetch.TokenOptions, ","))
	return nil
}

func run() error {
	args := os.Args
	if len(args) < 2 {
		return help("invalid arguments, missing command")
	}
	input := os.Getenv(configFileEnv)
	if input == "" {
		input = defaultConfig()
	}
	dryRun := true
	cmd := args[1]
	switch cmd {
	case "help":
		return help("")
	case check:
	case upgrade:
		dryRun = false
	default:
		return help(fmt.Sprintf("unknown argument: %s", cmd))
	}
	var appSet []string
	var disableSet []string
	if len(args) > 2 {
		set := flag.NewFlagSet("app", flag.ExitOnError)
		apps := set.String(appFlag, "", "limit application checks")
		disable := set.String(disableFlag, "", "disable applications")
		if err := set.Parse(args[2:]); err != nil {
			return err
		}
		appSet = commaList(apps)
		disableSet = commaList(disable)
		if len(appSet) > 0 && len(disableSet) > 0 {
			return help("can not limit applications and disable at the same time")
		}
	}
	if !core.PathExists(input) {
		return fmt.Errorf("config file does not exist: %s", input)
	}
	cfg, err := core.LoadConfig(input, dryRun, appSet, disableSet)
	if err != nil {
		return err
	}
	return cfg.Process(fetch.ResourceFetcher{})
}

func commaList(in *string) []string {
	if in == nil {
		return nil
	}
	v := strings.TrimSpace(*in)
	if v == "" {
		return nil
	}
	return strings.Split(v, ",")
}
