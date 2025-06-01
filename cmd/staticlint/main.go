package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"

	"github.com/aube/url-shortener/internal/staticlint/slerrors"
)

// Config is the name of the configuration file that the tool expects to find
// in the same directory as the executable.
const Config = `config.json`

// ConfigData describes the structure of the configuration file.
// It contains settings for which staticcheck analyzers to enable.
type ConfigData struct {
	Staticcheck []string // List of staticcheck analyzer names or "*" for all
}

// main is the entry point for the static analysis tool.
// It performs the following steps:
// 1. Locates and reads the configuration file
// 2. Sets up the default analyzers (error checks, printf, shadow, etc.)
// 3. Adds staticcheck analyzers based on configuration
// 4. Runs all selected analyzers using multichecker
func main() {
	appfile, err := os.Executable()
	if err != nil {
		panic(err)
	}
	data, err := os.ReadFile(filepath.Join(filepath.Dir(appfile), Config))
	if err != nil {
		panic(err)
	}
	var cfg ConfigData
	if err = json.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}
	mychecks := []*analysis.Analyzer{
		ErrCheckAnalyzer,
		ErrRunErrOSExit,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
	}
	checks := make(map[string]bool)
	for _, v := range cfg.Staticcheck {
		checks[v] = true
	}
	// добавляем анализаторы из staticcheck, которые указаны в файле конфигурации
	for _, v := range staticcheck.Analyzers {
		if checks["*"] || checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	multichecker.Main(
		mychecks...,
	)
}

// ErrCheckAnalyzer is a custom analyzer that checks for unchecked errors in function calls.
// It reports cases where errors returned from functions are not properly handled.
var ErrCheckAnalyzer = &analysis.Analyzer{
	Name: "errcheck",
	Doc:  "check for unchecked errors",
	Run:  slerrors.RunErrUnchecked,
}

// ErrRunErrOSExit is a custom analyzer that checks for usage of os.Exit in functions.
// It reports cases where os.Exit is called, which is generally discouraged
// in functions other than main().
var ErrRunErrOSExit = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "check for os.Exit",
	Run:  slerrors.RunErrOSExit,
}
