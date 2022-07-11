/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	b "goshelly-client/basic"
	t "goshelly-client/template"
	"io/ioutil"

	// "net/http"
	// "strings"

	"github.com/spf13/cobra"
)

const listLogURL = "http://localhost:9000/users/logs/"

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the last 5 GoShelly runs on your account.",
	Long: `Returns a list of a only the last 5 GoShelly runs on your account. No data from runs previous to the last 5 is maintained by the GoShelly server.
	If you wish to see the data from the runs earlier to the last 5, check under the 'logs' directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		getLogs()
	},
}

func getLogs() {
	resp := b.SendPOST(listLogURL, b.GetLoggedUser())
	var str t.Msg
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(resp.StatusCode, "Could not read response.")
		return
	}
	json.Unmarshal(body, &str)
	fmt.Println(str.MESSAGE)

}
func init() {
	rootCmd.AddCommand(listCmd)
}
