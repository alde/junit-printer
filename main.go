package main

import (
	"fmt"
	"os"

	gojunit "github.com/joshdk/go-junit"
	"github.com/logrusorgru/aurora"
	. "github.com/logrusorgru/aurora"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: junit-printer <junit file>")
		os.Exit(1)
	}
	file := os.Args[1]
	junitSuites, err := gojunit.IngestFile(file)
	if err != nil {
		logrus.WithError(err).Fatal("unable to parse junit file")
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Suite Name", "Status", "Passed", "Failed", "Skipped", "Error"})

	for _, s := range junitSuites {
		table.Append([]string{
			s.Name, status(s).String(), colorNumber(s.Totals.Passed, "passed"), colorNumber(s.Totals.Skipped, "skipped"), colorNumber(s.Totals.Failed, "failed"), colorNumber(s.Totals.Error, "error"),
		})
	}
	table.Render()
}

func status(s gojunit.Suite) aurora.Value {
	if s.Totals.Error == 0 && s.Totals.Failed == 0 {
		return Bold(Green("passed"))
	}
	if s.Totals.Error > 0 {
		return Bold(Red("error"))
	}
	if s.Totals.Failed > 0 {
		return Bold(Red("failed"))
	}
	return White("unknown")
}

func colorNumber(n int, numtype string) string {
	if n == 0 {
		return White(fmt.Sprintf("%d", n)).String()
	}
	if numtype == "skipped" {
		return Yellow(fmt.Sprintf("%d", n)).String()
	}
	if numtype == "failed" || numtype == "error" {
		return Red(fmt.Sprintf("%d", n)).String()
	}
	if numtype == "passed" {
		return Green(fmt.Sprintf("%d", n)).String()
	}
	return White(fmt.Sprintf("%d", n)).String()
}
