package set

import (
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/tars47/go-expense-tracker/db"
	"github.com/tars47/go-expense-tracker/printf"
	"github.com/tars47/go-expense-tracker/util"
)

func init() {

	d := time.Now().Format("02/01/2006")

	BudgetCmd.Flags().StringP(
		"month",
		"m",
		d[3:],
		"specify the month to set budget for",
	)
}

var BudgetCmd = &cobra.Command{
	Use:   "budget [-m MM/YYYY or MM]",
	Short: "Sets the budget for a given month(MM) or month/year(MM/YYYY) defaults to current month",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		var b db.Budget
		var err error

		b.Budget, err = strconv.Atoi(args[0])
		if err != nil || b.Budget < 1 {
			printf.Red("Error: budget value should be a number")
			return
		}

		m, _ := cmd.Flags().GetString("month")

		date, err := util.ValidateMonthFlag(m)
		if err != nil {
			return
		}

		b.Month, b.Year = int(date.Month()), date.Year()

		e, err := db.OpenDB()
		if err != nil {
			printf.Red("Error: There was an error accessing the database, %v", err)
			return
		}
		defer e.Close()

		err = e.UpsertBudget(b)
		if err != nil {
			printf.Red("Error: There was an error when adding budget to the database, %v", err)
			return
		}
		printf.Green("Budget for %v updated successfully", m)
	},
}
