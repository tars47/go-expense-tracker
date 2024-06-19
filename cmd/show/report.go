package show

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/tars47/go-expense-tracker/db"
	"github.com/tars47/go-expense-tracker/printf"
	"github.com/tars47/go-expense-tracker/util"
)

func init() {

	year, _, _ := time.Now().Date()

	ReportCmd.Flags().StringP(
		"year",
		"y",
		fmt.Sprint(year),
		"specify the year to generate report for",
	)
}

var ReportCmd = &cobra.Command{
	Use:   "report [-y YYYY]",
	Short: "Gets the yearly report for a given year, defaults to current year",
	Long: `
View your yearly expense report!  
Use expense show report with optional filters:
	
Examples:

expense show report
expense show report -y 2024
expense show report --year 2024
`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		y, _ := cmd.Flags().GetString("year")

		if yi, err := strconv.Atoi(y); err != nil || yi < 1 {
			printf.Red("Error: YEAR value should be of format YYYY")
			return
		}

		f, l := fmt.Sprintf("%v-01-01", y), fmt.Sprintf("%v-12-31", y)

		e, err := db.OpenDB()
		if err != nil {
			printf.Red("Error: There was an error accessing the database, %v", err)
			return
		}
		defer e.Close()

		rs, err := e.GetReport(f, l)
		if err != nil {
			printf.Red("Error: There was an error fetching reports from the database, %v", err)
			return
		}

		h := []string{"Month", "Year", "Spent", "Budget", "Saved", "Percentage"}
		var ro [][]string
		var tbu int
		var tsp int
		var tsa int

		for _, r := range rs {
			if r.Saved >= 0 {
				r.Percentage = fmt.Sprintf("saved %.1f%% of the budget", (float64(r.Saved)/float64(r.Budget))*100)
			} else {
				r.Percentage = fmt.Sprintf("exceeded the budget by %.1f%%", math.Abs((float64(r.Saved)/float64(r.Budget))*100))
			}
			ro = append(ro, []string{
				r.Month.String(),
				fmt.Sprint(r.Year),
				util.FmtNum(r.Spent),
				util.FmtNum(r.Budget),
				util.FmtNum(r.Saved),
				r.Percentage,
			})

			tbu += r.Budget
			tsp += r.Spent
			tsa += r.Saved
		}

		fmt.Println(util.Table(h, ro))

		printf.BlueS(" Spent\t\t", util.FmtNum(tsp))
		printf.BlueS(" Budget\t\t", util.FmtNum(tbu))

		if tbu-tsp > 0 {
			printf.BlueS(" Saved\t\t", util.FmtNum(tsa))
			printf.BlueS(" Summary\t", fmt.Sprintf("saved %.1f%% of the budget\n", (float64(tsa)/float64(tbu))*100))
		} else {
			printf.BlueS(" Saved\t\t", 0)
			printf.BlueS(" Summary\t", fmt.Sprintf("exceeded the budget by %.1f%%\n", math.Abs((float64(tsp)/float64(tbu))*100)))
		}
	},
}
