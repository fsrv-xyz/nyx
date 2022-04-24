package main

import (
    "github.com/fatih/color"
    "github.com/rodaine/table"
    "golang.fsrv.services/nyx/internal/check"
)

func renderOutput(checks []check.GenericCheck) {
    headerFmt := color.New(color.FgWhite, color.Italic, color.Underline).SprintfFunc()
    columnFmt := color.New(color.FgHiWhite).SprintfFunc()

    tbl := table.New("Name", "Result", "Error", "Help")
    tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt).WithPadding(5)

    for _, result := range checks {
        var resultColor color.Attribute
        var resultError string
        var resultHelp string

        switch result.State {
        case check.StateOK:
            resultColor = color.FgGreen
        case check.StateWarning:
            resultColor = color.FgYellow
        default:
            resultColor = color.FgRed
        }

        switch result.Error {
        case nil:
            resultError = "-"
        default:
            resultError = result.Error.Error()
            resultHelp = result.Help
        }
        tbl.AddRow(result.Name, color.New(resultColor, color.Bold).Sprint(result.State), resultError, resultHelp)
    }
    tbl.Print()
}
