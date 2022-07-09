package cmd

import (
	"fmt"
	b "goshelly-client/basic"
	t "goshelly-client/template"

	"github.com/spf13/cobra"
)

var loginUser t.LoginUser
const loginURL = "http://localhost:9000/users/login/"

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login into your GoShelly account.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_, loginUser.EMAIL, loginUser.PASSWORD = b.GetCredentials(0)
		msg, tkn := b.SendPOST(loginURL, loginUser)
		fmt.Println(msg)
		b.SaveLoginResult(tkn, loginUser.EMAIL)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
