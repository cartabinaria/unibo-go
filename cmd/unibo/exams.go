package main

import (
	"regexp"

	"github.com/cartabinaria/unibo-go/exams"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var examsCmd = &cobra.Command{
	Use:   "exams courseType courseId [subjectRegex]",
	Short: "Get the exams for a degree",
	Args:  cobra.RangeArgs(2, 3),
	Run:   runExams,
}

func init() {
	rootCmd.AddCommand(examsCmd)
}

func runExams(cmd *cobra.Command, args []string) {
	var err error

	var subjectRegex *regexp.Regexp
	if len(args) == 3 {
		subjectRegex, err = regexp.Compile("(?i)" + args[2])
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
	}

	e, err := exams.GetExams(args[0], args[1])
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	// filter exams by subject
	if subjectRegex != nil {
		var filtered []exams.Exam
		for _, exam := range e {
			if subjectRegex.MatchString(exam.SubjectName) {
				filtered = append(filtered, exam)
			}
		}
		e = filtered
	}

	// group exams by subject
	subjects := make(map[string][]exams.Exam)
	for _, exam := range e {
		subjects[exam.SubjectName] = append(subjects[exam.SubjectName], exam)
	}

	greenFmt := color.New(color.FgGreen).SprintFunc()
	grayFmt := color.New(color.FgHiBlack).SprintFunc()
	yellowFmt := color.New(color.FgYellow).SprintFunc()

	for subject, exams := range subjects {
		cmd.Printf("%s (%s)\n", greenFmt(subject), exams[0].Teacher)
		for _, exam := range exams {
			cmd.Printf("- %s: %-30s ", grayFmt("data"), exam.Date)
			cmd.Printf("%s: %-10s ", grayFmt("luogo"), exam.Location)
			cmd.Printf("%s: %-10s\n", grayFmt("tipo"), yellowFmt(exam.Type))
		}
		cmd.Println()
	}

	cmd.Printf("Total exams: %d\n", len(e))
	cmd.Printf("Total subjects: %d\n", len(subjects))
}
