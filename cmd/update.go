package cmd

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tars47/go-expense-tracker/db"
	"github.com/tars47/go-expense-tracker/printf"
	"github.com/tars47/go-expense-tracker/util"
)

func init() {

	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().IntP(
		"amount",
		"p",
		0,
		"specify the amount of the expense item",
	)
	updateCmd.Flags().StringP(
		"date",
		"d",
		"",
		"specify a date of the expense item",
	)
	updateCmd.Flags().StringP(
		"category",
		"c",
		"",
		"specify a category of the expense item",
	)
}

var updateCmd = &cobra.Command{
	Use:   "update ID [ITEM] [-amount AMOUNT] [-date DATE] [-category CATEGORY]",
	Short: "Update existing entries",
	Long: `
Update existing entries!  
Use expense update ID [ITEM] [-amount AMOUNT] [-date DATE] [-category CATEGORY].  
Replace ID with the entry's number, brackets indicate optional arguments.
	
Examples: assuming current date is 12/06/2024

expense update 123 -a 500 -c food
expense update 123 -a 500 -d 12 -c food
expense update 123 --amount 500 --category food
expense update 123 --amount 500 --date 12 -category food
expense update 123 --amount 500 --date 12/06/2024 -category food

all the above commands will produce the same results
	`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		var err error
		ex := db.Expense{}

		ex.Id, err = strconv.Atoi(args[0])
		if err != nil || ex.Id < 1 {
			printf.Red("Error: ID value is invalid, it should be a number")
			return
		}

		if len(args) > 1 {
			args = args[1:]
			ex.Item = strings.Join(args, " ")
		}

		ex.Amount, _ = cmd.Flags().GetInt("amount")
		ds, _ := cmd.Flags().GetString("date")
		ex.Category, _ = cmd.Flags().GetString("category")

		if ds != "" {
			ex.Date, err = util.ValidateDateFlag(ds)
			if err != nil {
				return
			}
		}

		e, err := db.OpenDB()
		if err != nil {
			printf.Red("Error: There was an error accessing the database, %v", err)
			return
		}
		defer e.Close()

		err = e.Update(ex)
		if err != nil {
			printf.Red("Error: There was an error when updating expense in the database, %v", err)
			return
		}

		printf.Green("Updated Successfully")

	},
}
