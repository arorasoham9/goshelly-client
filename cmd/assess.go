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
var COMPLETION_STATUS, RAW bool
var SSLEMAIL,HOST, PORT string
// demoCmd represents the demo command
var demoCmd = &cobra.Command{
	Use:   "assess",
	Short: "Creates a reverse shell and runs few CLI commands on your system from an external source.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var newUser t.User
		if cmd.Flags().Changed("IP") {
			RAW = true
			SSLEMAIL,_ = cmd.Flags().GetString("SSLEMAIL")
		}else{
			RAW = false
			fmt.Println("Flag missing, 'IP'. Defaulting to Araali backdoor.")
			if cmd.Flags().Changed("PORT"){
				fmt.Println("PORT cannot be changed if IP is NOT changed.")
				return
			}
			// b.InitRest()
			if !b.LoginStatus(GetDom() + statusURL) {
			newUser.EMAIL = b.GetCred()
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
				SSLEMAIL = newUser.EMAIL
			}else {
				return
			}
		} 
		}
	
		PORT, _ = cmd.Flags().GetString("PORT")
		if cmd.Flags().Changed("PORT") {
			_, portErr := strconv.ParseInt(PORT, 10, 64)
			if portErr != nil {
				fmt.Printf("PORT Error: Not a number.\n %s", portErr)
				os.Exit(1)
			}
		}
		
		
		HOST, _ = cmd.Flags().GetString("IP")

		LOGMAX, _ := cmd.Flags().GetInt("LOGMAX")
		CFGF, _ := cmd.Flags().GetBool("CFGF")
		COMPLETION_STATUS = b.StartClient(HOST, PORT, SSLEMAIL, LOGMAX,RAW, CFGF)
		if !COMPLETION_STATUS{
			fmt.Println("GoShelly Failed.")
			return 
		}
		shwlogCmd.Run(cmd,[]string{})
	},
}

func init() {
	rootCmd.AddCommand(demoCmd)
	rootCmd.PersistentFlags().String("PORT", "443", "PORT")
	rootCmd.PersistentFlags().String("IP", GetIP(), "Server IP") // replace GetIP() with dns
	rootCmd.PersistentFlags().String("SSLEMAIL", "default@default.com", "Email to generate SSL certificate")
	rootCmd.PersistentFlags().Int("LOGMAX", 50, "Number of log files to keep")
	rootCmd.PersistentFlags().Bool("CFGF", false, "Read config from file.")
	// rootCmd.PersistentFlags().Bool("RAW", false, " Just run the demo and return log, no need to auth.")
}
