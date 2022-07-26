package cmd

import (
	"encoding/json"
	"fmt"
	b "goshelly-client/basic"
	t "goshelly-client/template"
	"io/ioutil"

	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

const statusURL = "/auth/"

// demoCmd represents the demo command
var demoCmd = &cobra.Command{
	Use:   "assess",
	Short: "Creates a reverse shell, few commands are run on your system from an external source.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var newUser t.User
		if !b.LoginStatus(GetDom() + statusURL) {
			newUser.NAME, newUser.EMAIL = b.GetCred()
			newUser.PASSWORD = []byte("default")
			resp := b.SendPOST(GetDom()+signupURL, newUser)
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(resp.StatusCode, "Could not read response.")
				return
			}
			var obj t.Msg
			json.Unmarshal(body, &obj)
			fmt.Println(obj.MESSAGE)
			if resp.StatusCode == http.StatusCreated {
				LoginRun(GetDom()+loginURL, t.LoginUser{
					EMAIL:    newUser.EMAIL,
					PASSWORD: newUser.PASSWORD,
				})
			}
		}
		PORT, _ := cmd.Flags().GetString("PORT")
		if cmd.Flags().Changed("PORT") {
			_, portErr := strconv.ParseInt(PORT, 10, 64)
			if portErr != nil {
				fmt.Printf("PORT Error: Not a number.\n %s", portErr)
				os.Exit(1)
			}
		}
		if !cmd.Flags().Changed("IP") {
			fmt.Println("Flag missing, 'IP'. Defaulting to Araali backdoor.")
		}
		SSLEMAIL := b.GetLoggedUser().EMAIL
		if cmd.Flags().Changed("SSLEMAIL") {
			SSLEMAIL, _ = cmd.Flags().GetString("SSLEMAIL")
		}
		HOST, _ := cmd.Flags().GetString("IP")
		LOGMAX, _ := cmd.Flags().GetInt("LOGMAX")
		b.StartClient(HOST, PORT, SSLEMAIL, LOGMAX)
		shwlogCmd.Run(cmd,[]string{})
	},
}

func init() {
	rootCmd.AddCommand(demoCmd)
	rootCmd.PersistentFlags().String("PORT", "443", "PORT")
	rootCmd.PersistentFlags().String("IP", GetIP(), "Server IP")
	rootCmd.PersistentFlags().String("SSLEMAIL", "", "Email to generate SSL certificate.")
	rootCmd.PersistentFlags().Int("LOGMAX", 50, "Number of log files to keep")
	rootCmd.PersistentFlags().Bool("CFGF", false, "Read config from file.")
	rootCmd.PersistentFlags().Bool("RAW", false, "Just run the demo and return log, no need to auth.")
}
