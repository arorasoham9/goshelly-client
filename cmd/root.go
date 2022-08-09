package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	t "goshelly-client/template"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goshelly-client",
	Short: "",
	Long:  `Araali GoShelly is an open source tool that helps security teams safely test their detect and response readiness (the fire drill for SIEM/SOAR/EDR/NDR/XDR investment) for backdoors. This is typical when supply chain vulnerabilities like remote code execution (RCE) are exploited and represents a doomsday scenario where an attacker has full remote control capabilities based on the backdoor.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}

func GetIP() string {
	var config t.ApiConnIP
	file, err := ioutil.ReadFile("./config/api_conn_config.json")
	if err != nil {
		fmt.Println("Could not read in IP configuration. Err: ", err)
		os.Exit(1)
	}
	err = json.Unmarshal([]byte(file), &config)
	if err != nil {
		fmt.Println("Could not read in configuration. Err: ", err)
		os.Exit(1)
	}
	sDec, _ := base64.StdEncoding.DecodeString(config.IP)
	return string(sDec)
}


func GetDom() string{
	return "http://" + GetIP() + ":9000"
}
