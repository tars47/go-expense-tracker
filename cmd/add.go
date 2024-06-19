package cmd

import (
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tars47/go-expense-tracker/db"
	"github.com/tars47/go-expense-tracker/printf"

	"github.com/tars47/go-expense-tracker/util"
)

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().IntP(
		"amount",
		"a",
		0,
		"specify the amount of the expense item",
	)
	addCmd.Flags().StringP(
		"date",
		"d",
		time.Now().Format("02/01/2006"),
		"specify a date of the expense item",
	)
	addCmd.Flags().StringP(
		"category",
		"c",
		"general",
		"specify a category of the expense item",
	)
	addCmd.MarkFlagRequired("amount")
}

var addCmd = &cobra.Command{
	Use:   "add ITEM -amount AMOUNT [-date DATE] [-category CATEGORY]",
	Short: "Add new expenses to your tracker",
	Long: `
Add new expenses to your tracker!  
Use expense add ITEM -amount AMOUNT -date DATE -category CATEGORY.  
Replace item, amount, dates & categories with your details.
date and category is optional, defaults to current date and general category
	
Examples:

expense add lunch at mtr -a 500 -c food
expense add lunch at mtr -a 500 -d 12 -c food
expense add lunch at mtr --amount 500 --category food
expense add lunch at mtr --amount 500 --date 12 --category food
expense add lunch at mtr --amount 500 --date 12/06/2024 --category food

all the above commands will produce the same results
	`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		i := strings.Join(args, " ")
		if i == "" {
			printf.Red("Error: ITEM should not be empty")
			return
		}
		a, _ := cmd.Flags().GetInt("amount")
		ds, _ := cmd.Flags().GetString("date")
		c, _ := cmd.Flags().GetString("category")

		if a < 1 {
			printf.Red("Error: AMOUNT value should be greater than 0")
			return
		}

		d, err := util.ValidateDateFlag(ds)
		if err != nil {
			return
		}

		exp := db.Expense{
			Item:     i,
			Amount:   a,
			Date:     d,
			Category: c,
		}

		e, err := db.OpenDB()
		if err != nil {
			printf.Red("Error: There was an error accessing the database, %v", err)
			return
		}
		defer e.Close()

		if b, _ := e.GetBudget(int(d.Month()), d.Year()); b.Budget == 0 {
			ds := d.Format("02/01/2006")
			my := ds[3:]
			printf.Red("Error: Please set the budget for %v before adding expenses ", my)
			printf.Blue("expense set budget AMOUNT -m %v", my)
			return
		}

		if err = e.Insert(exp); err != nil {
			printf.Red("Error: There was an error when adding expense to the database, %v", err)
			return
		}

		printf.Green("Added Successfully")
	},
}
