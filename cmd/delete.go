package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tars47/go-expense-tracker/db"
	"github.com/tars47/go-expense-tracker/printf"
)

func init() {

	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete ID",
	Short: "Delete an expenses from your tracker",
	Long: `
Permanently remove expenses!  
Use expense delete ID.  
Replace ID with the unique number of the expense you wish to remove.  
Caution: This action cannot be undone.

expense delete 1 2 3 4

This command deletes expense id numbers 1, 2, 3 and 4
	`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		ids := make([]int, 0, 10)
		for _, sid := range args {
			id, err := strconv.Atoi(sid)
			if err != nil || id < 1 {
				printf.Red("Error: ID value : %v is invalid\n", sid)
				return
			}
			ids = append(ids, id)
		}

		e, err := db.OpenDB()
		if err != nil {
			printf.Red("Error: There was an error accessing the database, %v", err)
			return
		}
		defer e.Close()

		rows, err := e.Delete(ids)
		if err != nil {
			printf.Red("Error: There was an error when deleting expenses from the database, %v", err)
			return
		}

		printf.Green("Deleted %v expenses Successfully", rows)
	},
}
