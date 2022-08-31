package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset GoShelly",
	Long:  `Removes logged user configuration from your computer.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Resetting...")
		os.Remove("./config/token-config.json")
		fmt.Printf("Done.")
	},
}


func init() {
	rootCmd.AddCommand(resetCmd)
}
