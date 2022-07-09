package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete account.",
	Long: `Delete's all existence of the user's account and data from the GoShelly Server.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("NOTE: Running delete will delete your account permanently.")
		fmt.Println("All data associated with your account will also be removed without the option of being restored later.")

	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
