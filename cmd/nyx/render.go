package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"golang.fsrv.services/nyx/internal/check"
)

func renderOutput(checks []check.GenericCheck) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "State", "Error", "Help"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(true)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(false)

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)

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
			resultHelp = "-"
		default:
			resultError = result.Error.Error()
			resultHelp = result.Help
		}
		table.Append([]string{result.Name, color.New(resultColor, color.Bold).Sprint(result.State), resultError, resultHelp})
	}
	table.Render()
}
