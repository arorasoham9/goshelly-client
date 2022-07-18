package cmd

import (
	"fmt"
	b "goshelly-client/basic"
	t "goshelly-client/template"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var NAME, EMAIL string

// assessCmd represents the assess command
var assessCmd = &cobra.Command{
	Use:   "assess",
	Short: "Attempt to run a few commands on you computer remotely and verify your attack readiness.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var newUser t.User
		if !b.LoginStatus(GetDom() + statusURL) {
			newUser.NAME, newUser.EMAIL, newUser.PASSWORD = b.GetCredentials(1,3)
			resp := b.SendPOST(GetDom()+signupURL, newUser)
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
			fmt.Println("Flag missing, 'IP'.")
			os.Exit(1)
		}
		SSLEMAIL := b.GetLoggedUser().EMAIL
		if cmd.Flags().Changed("SSLEMAIL") {
			SSLEMAIL, _ = cmd.Flags().GetString("SSLEMAIL")
		}
		HOST, _ := cmd.Flags().GetString("IP")
		LOGMAX, _ := cmd.Flags().GetInt("LOGMAX")
		b.StartClient(HOST, PORT, SSLEMAIL, LOGMAX)

	},
}

func init() {
	rootCmd.AddCommand(assessCmd)
}
