package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"golang.fsrv.services/nyx/internal/check"
	"golang.fsrv.services/nyx/internal/util"
	"golang.fsrv.services/version"
)

type application struct {
	configFilePath string
}

// instantiate application state variable
var instance application

func init() {
	var printVersion bool

	// parse input parameters
	flag.StringVar(&instance.configFilePath, "config.file", "./nyx.json", "path to config file")
	flag.BoolVar(&printVersion, "version", false, "run as server and export metrics via http")
	flag.Parse()

	// version handling
	if printVersion {
		fmt.Println(version.Print("nyx"))
		os.Exit(0)
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
	if errors.Is(err, jsonDecodeError) {
		fmt.Printf("fail to decode config file %+q\n%+q\n", instance.configFilePath, jsonDecodeError)
		os.Exit(1)
	}
	if jsonDecodeError != nil {
		panic(jsonDecodeError)
	}

	m := runChecks(config)
	p := util.SortByCheckName(m)
	renderOutput(p)
}

func runChecks(config configuration) []check.GenericCheck {
	var overview []check.GenericCheck
	blaC := make(chan check.GenericCheck)

	go func() {
		for checkIndex := range config.Checks {
			go runCheck(config.Checks[checkIndex], blaC)
		}
	}()

	for result := range blaC {
		overview = append(overview, result)
		if len(overview) == len(config.Checks) {
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
	checkInstance.StartTiming()
	checkInstance.Run()
	checkInstance.FinishTiming()

	output <- checkInstance.Export()
}
