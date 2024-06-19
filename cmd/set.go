package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tars47/go-expense-tracker/cmd/set"
)

func init() {

	rootCmd.AddCommand(setCmd)
	setCmd.AddCommand(set.BudgetCmd)
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Sets a value",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
