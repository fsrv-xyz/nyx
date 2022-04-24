package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "os"

    "golang.fsrv.services/nyx/internal/check"
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
    if err != nil {
        panic(err)
    }

    config := configuration{}
    jsonDecodeError := json.NewDecoder(file).Decode(&config)
    if jsonDecodeError != nil {
        panic(jsonDecodeError)
    }

    var overview []check.GenericCheck
    for _, che := range config.Checks {
        overview = append(overview, runCheck(che))
    }
    renderOutput(overview)
}

func runCheck(config checkConfiguration) check.GenericCheck {
    checkFactory, ok := check.RegistryInstance.Checks[config.Check]
    checkInstance := checkFactory()

    if !ok {
        return check.GenericCheck{Error: fmt.Errorf("%+q not implemented\n", config.Check)}
    }
    for key, value := range config.Parameter {
        checkInstance.SetParameter(key, value)
    }
    checkInstance.SetName(config.Name)
    checkInstance.SetHelp(config.Help)
    checkInstance.Run()

    return checkInstance.Export()
}
