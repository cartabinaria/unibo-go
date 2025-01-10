package main

import (
	"encoding/csv"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/cartabinaria/unibo-go/exams"
	"github.com/cartabinaria/unibo-go/rubrica"
)

var examsCmd = &cobra.Command{
	Use:     "exams courseType courseId [subjectRegex]",
	Short:   "Get the exams for a degree",
	Aliases: []string{"e"},
	Args:    cobra.RangeArgs(2, 3),

	Run: runExams,
}

var outputFmt string

func init() {
	rootCmd.AddCommand(examsCmd)
	examsCmd.Flags().StringVarP(&outputFmt, "format", "f", "human", "output format (human, csv)")
}

func runExams(cmd *cobra.Command, args []string) {
	if outputFmt != "human" && outputFmt != "csv" {
		Errorln("invalid output format:", outputFmt)
		return
	}

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

	if outputFmt == "human" {
		printToConsole(cmd, subjects, greenFmt, grayFmt, yellowFmt, e)
	} else if outputFmt == "csv" {
		// map to CSV entries
		type csvEntry struct {
			Subject      string
			Teacher      string
			TeacherEmail string
			Date         string
			Location     string
		}

		var csvEntries []csvEntry
		for _, exam := range e {
			names := strings.Split(exam.Teacher, " ")
			numNames := len(names)

			var lastName, firstName string

			if numNames == 2 {
				lastName = names[0]
				firstName = names[1]
			} else {
				// TODO: we should try to split the name in different ways an try the lookup
				// as some teachers have multiple first names and some have multiple last names
				// for now, we just use the last word as the last name and the rest as the first name
				lastName = names[numNames-1]
				firstName = strings.Join(names[:numNames-1], " ")
			}

			lastName = strings.ToLower(lastName)
			firstName = strings.ToLower(firstName)

			var teacherEmail string

			// lookup teacher email
			contact, err := rubrica.Search(firstName, lastName)
			if err != nil {
				Errorln(err)
				return
			} else if len(contact) > 1 {
				cmd.PrintErrln("multiple contacts found for teacher", exam.Teacher)
				// do not set teacherEmail
			} else if len(contact) == 0 {
				cmd.PrintErrln("no contact found for teacher", exam.Teacher)
				// do not set teacherEmail
			} else if len(contact) == 1 {
				teacherEmail = contact[0].Email
			}

			csvEntries = append(csvEntries, csvEntry{
				Subject:      exam.SubjectName,
				Teacher:      exam.Teacher,
				TeacherEmail: teacherEmail,
				Date:         exam.Date.Format(time.DateTime),
				Location:     exam.Location,
			})
		}

		// print CSV header
		writer := csv.NewWriter(cmd.OutOrStdout())
		_ = writer.Write([]string{"Subject", "Teacher", "TeacherEmail", "Date", "Location"})
		for _, entry := range csvEntries {
			_ = writer.Write([]string{entry.Subject, entry.Teacher, entry.TeacherEmail, entry.Date, entry.Location})
		}
	}
}

func printToConsole(
	cmd *cobra.Command,
	subjects map[string][]exams.Exam,
	greenFmt func(a ...interface{}) string,
	grayFmt func(a ...interface{}) string,
	yellowFmt func(a ...interface{}) string,
	e []exams.Exam,
) {

	formatDate := func(date time.Time) string {
		return date.Format(time.DateTime)
	}

	formatLocation := func(location string) string {
		return location
	}

	formatType := func(examType string) string {
		return examType
	}

	for subject, subjectExams := range subjects {
		teacher := subjectExams[0].Teacher

		longestDate := 0
		longestLocation := 0
		longestType := 0

		for _, exam := range subjectExams {
			if len(formatDate(exam.Date)) > longestDate {
				longestDate = len(formatDate(exam.Date))
			}
			if len(formatLocation(exam.Location)) > longestLocation {
				longestLocation = len(formatLocation(exam.Location))
			}
			if len(formatType(exam.Type)) > longestType {
				longestType = len(formatType(exam.Type))
			}
		}

		cmd.Printf("%s (%s)\n", greenFmt(subject), teacher)
		for _, exam := range subjectExams {

			date := formatDate(exam.Date)
			location := formatLocation(exam.Location)
			examType := formatType(exam.Type)

			cmd.Printf("- %s: %-*s %s: %-*s %s: %-*s\n",
				grayFmt("data"), longestDate, yellowFmt(date),
				grayFmt("luogo"), longestLocation, location,
				grayFmt("tipo"), longestType, yellowFmt(examType))
		}
		cmd.Println()
	}

	cmd.Printf("Total exams: %d\n", len(e))
	cmd.Printf("Total subjects: %d\n", len(subjects))
}
