package cmd

import (
	"fmt"
	b "goshelly-client/basic"
	"strings"

	"github.com/spf13/cobra"
)

var deleteURL = "/delete/"

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete account.",
	Long:  `Delete's all existence of the user's account data from the GoShelly Server.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !b.LoginStatus(GetDom()+statusURL) {
			fmt.Println("No user found.")
			return
		}
		var temp string
		var confirm bool
		fmt.Printf("NOTE: Running delete will delete all previous logs.")
		fmt.Println("All data associated with your account will be removed permanently.")
		fmt.Printf("Are you sure you would like to delete your account? (Y/N) --> ")
		fmt.Scanf("%s", &temp)
		temp = strings.ToLower(temp)
		switch temp {
		case "y":
			confirm = true
		case "n":
			confirm = false
		default:
			return
		}

		b.DeleteUser(confirm, GetDom()+deleteURL)
		
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
