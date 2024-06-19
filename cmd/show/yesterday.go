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

	YesterdayCmd.Flags().StringP(
		"group",
		"g",
		"",
		"specify a group condition can be any one of [date, category]",
	)
}

var YesterdayCmd = &cobra.Command{
	Use:   "yesterday [--group GROUP]",
	Short: "View your yesterday's tracked expenses",
	Long: `
View your yesterday's tracked expenses!  Use expense show today with optional filters:

--group GROUP to categorize expenses by a group, can be any of date or category.
	
Examples:

expense show yesterday
expense show yesterday -g date
expense show yesterday -g category
expense show yesterday --group date
expense show yesterday --group category
	`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		g, _ := cmd.Flags().GetString("group")

		if len(g) > 0 && slices.Index([]string{"date", "category"}, g) == -1 {
			printf.Red("Error: GROUP value should be one of date or category")
			return
		}

		y := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

		e, err := db.OpenDB()
		if err != nil {
			printf.Red("Error: There was an error accessing the database, %v", err)
			return
		}
		defer e.Close()

		var es []db.Expense

		es, err = e.ListRange(y, y, g)
		if err != nil {
			printf.Red("Error: There was an error querying expense from the database, %v", err)
			return
		}

		s := util.PrintTable(g, es)

		printf.BlueS("Spent\t\t", util.FmtNum(s))
	},
}
