package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"

	"github.com/fsrv-xyz/nyx/internal/check"
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
		var resultError = "-"
		var resultHelp = "-"

		switch result.State {
		case check.StateOK:
			resultColor = color.FgGreen
		case check.StateWarning:
			resultColor = color.FgYellow
		default:
			resultColor = color.FgRed
		}

		if result.Error != nil {
			resultError = result.Error.Error()
			resultHelp = result.Help
		}

		if resultHelp == "" {
			resultHelp = "-"
		}
		table.Append([]string{result.Name, color.New(resultColor, color.Bold).Sprint(result.State), resultError, resultHelp})
	}
	table.Render()
}
