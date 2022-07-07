/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	t "goshelly-client/template"
	"github.com/spf13/cobra"
	b "goshelly-client/basic"
)

var LOGINUSER t.User
var LOGINURL = ""




var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login into your GoShelly account.",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		b.GetCredentials(LOGINUSER)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
