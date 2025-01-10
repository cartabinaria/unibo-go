package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/cartabinaria/unibo-go/timetable"
)

var cmdTimetable = &cobra.Command{
	Use:     "timetable courseType courseId year [curriculum]",
	Short:   "fetches the timetable of a degree course",
	Example: "unibo timetable laurea informatica 1",
	Aliases: []string{"t", "tt"},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return fmt.Errorf("requires at least 3 arguments")
		}
		return nil
	},
	Run: runTimetable,
}

func init() {
	rootCmd.AddCommand(cmdTimetable)
}

func runTimetable(cmd *cobra.Command, args []string) {

	courseType := args[0]
	courseId := args[1]
	yearStr := args[2]

	year64, err := strconv.ParseInt(yearStr, 10, 32)
	if err != nil {
		Errorln("year must be a number")
		cmd.Help()
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
		Errorf("error fetching timetable: %v\n", err)
		return
	}

	if len(tt) == 0 {
		fmt.Println(yellowFmt("No lessons found\n"))
		return
	}

	for _, e := range tt {
		fmt.Printf("- %s -> %s: %-50s %-30s (%s)\n",
			greenFmt(e.Start.Format("15:04")), redFmt(e.End.Format("15:04")),
			e.Title, yellowFmt(e.Teacher), grayFmt(e.CodModulo))
	}
}
