package cmd

import (
	"slices"
	"time"

	"github.com/spf13/cobra"
	"github.com/tars47/go-expense-tracker/cmd/show"
	"github.com/tars47/go-expense-tracker/db"
	"github.com/tars47/go-expense-tracker/printf"
	"github.com/tars47/go-expense-tracker/util"
)

func init() {

	d := time.Now().Format("02/01/2006")

	rootCmd.AddCommand(showCmd)

	showCmd.Flags().StringP(
		"month",
		"m",
		d[3:],
		"specify the month for list of expenses",
	)
	showCmd.Flags().StringP(
		"group",
		"g",
		"",
		"specify a group condition can be any one of [date, category]",
	)

	showCmd.AddCommand(show.WeekCmd)
	showCmd.AddCommand(show.TodayCmd)
	showCmd.AddCommand(show.YesterdayCmd)
	showCmd.AddCommand(show.ReportCmd)

}

var showCmd = &cobra.Command{
	Use:   "show [--month MM/YYYY] [--group GROUP]",
	Short: "View your tracked expenses",
	Long: `
View your tracked expenses!  Use expense show with optional filters:

--month MM or --month MM/YYYY for a specific month (e.g., -m 06/2024).
--group GROUP to categorize expenses by a group, can be any of date or category.
	
Examples:

expense show
expense show -m 06
expense show -m 06/2024
expense show -m 06/2024 -g date
expense show -m 06/2024 -g category
expense show --month 06
expense show --month 06/2024
expense show --month 06/2024 --group date
expense show --month 06/2024 --group category
	`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		m, _ := cmd.Flags().GetString("month")
		g, _ := cmd.Flags().GetString("group")

		if len(g) > 0 && slices.Index([]string{"date", "category"}, g) == -1 {
			printf.Red("Error: GROUP value should be one of date or category")
			return
		}

		d, err := util.ValidateMonthFlag(m)
		if err != nil {
			return
		}

		f, l := getFirstAndLast(int(d.Month()), d.Year())

		e, err := db.OpenDB()
		if err != nil {
			printf.Red("Error: There was an error accessing the database, %v", err)
			return
		}
		defer e.Close()

		var es []db.Expense

		es, err = e.ListRange(f, l, g)
		if err != nil {
			printf.Red("Error: There was an error querying expense from the database, %v", err)
			return
		}

		s := util.PrintTable(g, es)

		util.PrintSummary(s, d, e)
	},
}

func getFirstAndLast(m, y int) (string, string) {

	loc := time.Now().Location()

	f := time.Date(y, time.Month(m), 1, 0, 0, 0, 0, loc)
	l := f.AddDate(0, 1, -1)

	return f.Format("2006-01-02"), l.Format("2006-01-02")
}
