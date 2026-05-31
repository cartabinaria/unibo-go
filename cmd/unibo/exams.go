// SPDX-FileCopyrightText: 2024 - 2025 Eyad Issa <eyadlorenzo@gmail.com>
//
// SPDX-License-Identifier: MIT

package main

import (
	"encoding/csv"
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/karlseguin/ccache/v3"
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
	examsCmd.Flags().StringVarP(&outputFmt, "format", "f", "human", "output format (human, csv, json)")
}

var contacts = ccache.New(ccache.Configure[[]rubrica.Contact]().MaxSize(1000))

func runExams(cmd *cobra.Command, args []string) {
	if outputFmt != "human" && outputFmt != "csv" && outputFmt != "json" {
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

	type entry struct {
		Subject      string `json:"subject"`
		Teacher      string `json:"teacher"`
		TeacherEmail string `json:"teacher_email,omitempty"`
		Date         string `json:"date"`
		Location     string `json:"location"`
		Type         string `json:"type"`
	}

	var entries []entry
	if outputFmt == "csv" || outputFmt == "json" {
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
			var contact []rubrica.Contact
			cacheKey := firstName + " " + lastName
			if item := contacts.Get(cacheKey); item != nil {
				contact = item.Value()
			} else {
				contact, err = rubrica.Search(firstName, lastName)
				if err != nil {
					Errorln(err)
					return
				}
				contacts.Set(cacheKey, contact, time.Hour)
			}

			if len(contact) > 1 {
				cmd.PrintErrln("multiple contacts found for teacher", exam.Teacher)
				// do not set teacherEmail
			} else if len(contact) == 0 {
				cmd.PrintErrln("no contact found for teacher", exam.Teacher)
				// do not set teacherEmail
			} else if len(contact) == 1 {
				teacherEmail = contact[0].Email
			}

			entries = append(entries, entry{
				Subject:      exam.SubjectName,
				Teacher:      exam.Teacher,
				TeacherEmail: teacherEmail,
				Date:         exam.Date.Format(time.DateTime),
				Location:     exam.Location,
				Type:         exam.Type,
			})
		}
	}

	switch outputFmt {
	case "human":
		printToConsole(cmd, subjects, greenFmt, grayFmt, yellowFmt, e)
	case "csv":
		// print CSV header
		writer := csv.NewWriter(cmd.OutOrStdout())
		_ = writer.Write([]string{"Subject", "Teacher", "TeacherEmail", "Date", "Location", "Type"})
		for _, e := range entries {
			_ = writer.Write([]string{e.Subject, e.Teacher, e.TeacherEmail, e.Date, e.Location, e.Type})
		}
		writer.Flush()
	case "json":
		encoder := json.NewEncoder(cmd.OutOrStdout())
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(entries); err != nil {
			cmd.PrintErrln(err)
		}
	}
}

func printToConsole(
	cmd *cobra.Command,
	subjects map[string][]exams.Exam,
	greenFmt func(a ...any) string,
	grayFmt func(a ...any) string,
	yellowFmt func(a ...any) string,
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
