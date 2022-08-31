package cmd

import (
	"encoding/json"
	"fmt"
	b "goshelly-client/basic"
	t "goshelly-client/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"github.com/spf13/cobra"
)

const logURL = "/link/"
const discprompt = `NOTE: Each requested log is made available for a short period of time at the URL returned.
Changing params of this link will NOT lead to other logs. This link is not transferrable, 
only the requesting computer will be able to access the contents of the link generated.
	`

// shwlogCmd represents the shwlog command
var shwlogCmd = &cobra.Command{
	Use:   "shwlog",
	Short: "See logs from your GoShelly runs.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if !b.LoginStatus(GetDom() + statusURL) {
			fmt.Println("No user found.")
			return
		}
		inputIDS, _ := cmd.Flags().GetString("ID")
		genLinks(inputIDS)
	},
}

func genLinks(ids string) {
	fmt.Printf("Fetching log for log ID: ")
	switch ids{
	case "-1":
		fmt.Println("latest")
	default:
		fmt.Println(ids)
	}
	var obj t.Msg
	tempUser := b.GetLoggedUser()
	i, _ := strconv.Atoi(ids)
	user := t.UserLinks{
		EMAIL: tempUser.EMAIL,
		TOKEN: tempUser.TOKEN,
		LOGID: i,
	}
	resp := b.SendPOST(GetDom()+logURL, user)
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(resp.StatusCode, "Could not read response.")
		return
	}

	json.Unmarshal(body, &obj)
	if resp.StatusCode != http.StatusOK {
		fmt.Println(obj.MESSAGE)
		return
	}
	
	obj.MESSAGE = "https://"+GetIP()+obj.MESSAGE
	u, err := url.Parse(obj.MESSAGE)
	if err != nil {
		fmt.Println("Could not parse log link. Err: ",err)
		return
	}
	fmt.Printf("\nYou can find the log here: %+v\n\n", u)
	fmt.Println(discprompt) 
}

func init() {
	rootCmd.AddCommand(shwlogCmd)
	rootCmd.PersistentFlags().String("ID", "-1", "Host the log data for a Goshelly run ID.")
}
