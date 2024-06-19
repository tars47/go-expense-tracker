package show

import (
	"slices"
	"time"

	"github.com/spf13/cobra"
	"github.com/tars47/go-expense-tracker/db"
	"github.com/tars47/go-expense-tracker/printf"
	"github.com/tars47/go-expense-tracker/util"
)

func init() {

	TodayCmd.Flags().StringP(
		"group",
		"g",
		"",
		"specify a group condition can be any one of [date, category]",
	)
}

var TodayCmd = &cobra.Command{
	Use:   "today [--group GROUP]",
	Short: "View your todays tracked expenses",
	Long: `
View your todays tracked expenses!  Use expense show today with optional filters:

--group GROUP to categorize expenses by a group, can be any of date or category.
	
Examples:

expense show today
expense show today -g date
expense show today -g category
expense show today --group date
expense show today --group category
	`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		g, _ := cmd.Flags().GetString("group")

		if len(g) > 0 && slices.Index([]string{"date", "category"}, g) == -1 {
			printf.Red("Error: GROUP value should be one of date or category")
			return
		}

		d := time.Now()
		t := d.Format("2006-01-02")

		e, err := db.OpenDB()
		if err != nil {
			printf.Red("Error: There was an error accessing the database, %v", err)
			return
		}
		defer e.Close()

		var es []db.Expense

		es, err = e.ListRange(t, t, g)
		if err != nil {
			printf.Red("Error: There was an error querying expense from the database, %v", err)
			return
		}

		s := util.PrintTable(g, es)

		util.PrintSummary(s, d, e)
	},
}
