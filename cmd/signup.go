/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	t "goshelly-client/template"
	"time"
	b "goshelly-client/basic"

	"github.com/spf13/cobra"
)

var SIGNUPURL = ""
var NEWUSER t.User


var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Create an account with Araali GoShelly",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("NOTE: A valid email and a password is required to sign up. Passwords cannot be reset.")
		time.Sleep(time.Second*2)
		b.GetCredentials(NEWUSER)
	},
}

func createUser(){



}
func init() {
	rootCmd.AddCommand(signupCmd)

}
