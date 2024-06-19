package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
}

var rootCmd = &cobra.Command{
	Use:   "expense",
	Short: "A CLI expense tracker tool for managing your monthly expenses.",
	Long: `
Track your finances with ease!  
expense helps you add, list, and analyze your spending.  
Run expense help for detailed commands.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
