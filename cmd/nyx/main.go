package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/fsrv-xyz/version"

	"github.com/fsrv-xyz/nyx/internal/check"
	"github.com/fsrv-xyz/nyx/internal/util"
)

type application struct {
	configFilePath  string
	checkIdentifier string
}

// instantiate application state variable
var instance application

func init() {
	var printVersion bool

	// parse input parameters
	flag.BoolVar(&printVersion, "version", false, "run as server and export metrics via http")
	flag.StringVar(&instance.checkIdentifier, "identifier", "all", "name of check to return")
	flag.StringVar(&instance.configFilePath, "config.file", "./nyx.json", "path to config file; alternatively, use environment variable NYX_CONFIG")
	flag.Parse()

	// version handling
	if printVersion {
		fmt.Println(version.Print("nyx"))
		os.Exit(0)
	}

	configFileEnv := os.Getenv("NYX_CONFIG")
	if configFileEnv != "" {
		instance.configFilePath = configFileEnv
	}
}

func main() {
	file, err := os.Open(instance.configFilePath)

	// fail if config file is not found
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("config file %+q not found\n", instance.configFilePath)
		os.Exit(1)
	}
	// fail on other errors
	if err != nil {
		panic(err)
	}

	// decode configuration file to struct; error and exit if decoding fails
	config := configuration{}
	jsonDecodeError := json.NewDecoder(file).Decode(&config)
	if jsonDecodeError != nil {
		fmt.Printf("fail to decode config file %+q\n%+q\n", instance.configFilePath, jsonDecodeError)
		os.Exit(1)
	}

	// filter checks by identifier
	checks, checkSelectionError := config.filterCheckByIdentifier(instance.checkIdentifier)
	if checkSelectionError != nil {
		fmt.Println(checkSelectionError)
		os.Exit(127)
	}

	results := runChecks(checks)
	renderOutput(util.SortByCheckName(results))
}

// filter checks by identifier; return all checks if identifier is "all" or "any"
func (config *configuration) filterCheckByIdentifier(identifier string) ([]checkConfiguration, error) {
	// return all checks if identifier is "all"
	if identifier == "all" || identifier == "any" {
		return config.Checks, nil
	}

	var checks []checkConfiguration
	for configCheckIndex := range config.Checks {
		if config.Checks[configCheckIndex].Identifier == identifier {
			checks = append(checks, config.Checks[configCheckIndex])
		}
	}
	if len(checks) == 0 {
		return nil, fmt.Errorf("no check found with identifier %+q", identifier)
	}
	return checks, nil
}

func runChecks(checks []checkConfiguration) []check.GenericCheck {
	var overview []check.GenericCheck
	blaC := make(chan check.GenericCheck)

	go func() {
		for checkIndex := range checks {
			go runCheck(checks[checkIndex], blaC)
		}
	}()

	for result := range blaC {
		overview = append(overview, result)
		if len(overview) == len(checks) {
			close(blaC)
		}
	}
	return overview
}

func runCheck(config checkConfiguration, output chan check.GenericCheck) {
	checkFactory, ok := check.RegistryInstance.Checks[config.Check]

	if !ok {
		output <- check.GenericCheck{Error: fmt.Errorf("%+q not implemented", config.Check), State: check.StateUnable}
		return
	}

	checkInstance := checkFactory()

	for key, value := range config.Parameter {
		checkInstance.SetParameter(key, value)
	}
	checkInstance.SetName(config.Name)
	checkInstance.SetHelp(config.Help)
	checkInstance.SetIdentifier(config.Identifier)
	checkInstance.StartTiming()
	checkInstance.Run()
	checkInstance.FinishTiming()

	output <- checkInstance.Export()
}
