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
		if b.LoginStatus(statusURL) {
			fmt.Println("Already logged in as: ", b.GetLoggedUser().EMAIL)
			return
		}
		_, loginUser.EMAIL, loginUser.PASSWORD = b.GetCredentials(0)
		LoginRun(loginURL,loginUser)
	},
}

func LoginRun(url string, user t.LoginUser) {
	resp := b.SendPOST(url, user)
	b.SaveLoginResult(resp, user.EMAIL)

}
func init() {
	rootCmd.AddCommand(loginCmd)
}
