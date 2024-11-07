package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cartabinaria/unibo-go/timetable"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var cmdTimetable = &cobra.Command{
	Use:     "timetable",
	Short:   "fetches the timetable of a degree course",
	Args:    cobra.RangeArgs(3, 4),
	Aliases: []string{"t"},
	Run:     runTimetable,
}

func init() {
	rootCmd.AddCommand(cmdTimetable)
}

var (
	yellowFmt = color.New(color.FgYellow).SprintFunc()
	greenFmt  = color.New(color.FgGreen).SprintFunc()
	redFmt    = color.New(color.FgRed).SprintFunc()
	grayFmt   = color.New(color.Faint).SprintFunc()
)

func runTimetable(cmd *cobra.Command, args []string) {

	courseType := args[0]
	courseId := args[1]
	yearStr := args[2]

	year64, err := strconv.ParseInt(yearStr, 10, 32)
	if err != nil {
		cmd.PrintErrln(fmt.Errorf("year must be a number"))
		return
	}

	year := int(year64)

	var curriculum string
	if len(args) == 4 {
		curriculum = args[3]
	}

	today := time.Now().Truncate(24 * time.Hour)

	interval := &timetable.Interval{Start: today, End: today}
	tt, err := timetable.FetchTimetable(courseType, courseId, curriculum, year, interval)
	if err != nil {
		cmd.PrintErrln(fmt.Errorf("error fetching timetable: %v", err))
		return
	}

	if len(tt) == 0 {
		cmd.Println(yellowFmt("No lessons found"))
		return
	}

	for _, e := range tt {
		cmd.Printf("- %s -> %s: %-50s %-30s (%s)\n",
			greenFmt(e.Start.Format("15:04")), redFmt(e.End.Format("15:04")),
			e.Title, yellowFmt(e.Teacher), grayFmt(e.CodModulo))
	}
}
