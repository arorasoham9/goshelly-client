package cmd

import (
	"fmt"
	b "goshelly-client/basic"
	t "goshelly-client/template"

	"github.com/spf13/cobra"
)

var loginUser t.LoginUser

const loginURL = "/login/"

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login into your GoShelly account.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if b.LoginStatus(GetDom()+statusURL) {
			fmt.Println("Already logged in as: ", b.GetLoggedUser().EMAIL)
			return
		}
		 loginUser.EMAIL, loginUser.PASSWORD = b.GetCredentials(0,5)
		obj := LoginRun(GetDom()+loginURL,loginUser)
		fmt.Println(obj.MESSAGE)
	},
}

func LoginRun(url string, user t.LoginUser) t.LogSuccess {
	resp := b.SendPOST(url, user)
	check, obj := b.SaveLoginResult(resp, user.EMAIL)
	if !check{
		fmt.Println("Unable to run GoShelly. User created, token not stored.")
	}
	return obj
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
