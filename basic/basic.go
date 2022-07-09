package basic

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	t "goshelly-client/template"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/mail"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func logClean(dir string) {
	files, _ := ioutil.ReadDir(dir)
	if len(files) < CONFIG.MAXLOGSTORE {
		return
	}

	var newestFile string
	var oldestTime = math.Inf(1)
	for _, f := range files {

		fi, err := os.Stat(dir + f.Name())
		if err != nil {
			fmt.Println(err)
		}
		currTime := float64(fi.ModTime().Unix())
		if currTime < oldestTime {
			oldestTime = currTime
			newestFile = f.Name()
		}
	}
	os.Remove(dir + newestFile)
}

// file upl /downl functions, if needed
func uploadFile(conn *tls.Conn, path string) {
	// open file to upload
	fi, err := os.Open(path)
	handleError(err)
	defer fi.Close()
	// upload
	_, err = io.Copy(conn, fi)
	handleError(err)
}

func returnLog() {
	bytearr, err := ioutil.ReadFile(CONFIG.LOGNAME)
	if err != nil {
		fmt.Println("Could not get logs.")
		return
	}
	fmt.Println(string(bytearr))

}

func downloadFile(conn *tls.Conn, path string) {
	// create new file to hold response
	fo, err := os.Create(path)
	handleError(err)
	defer fo.Close()

	handleError(err)
	defer conn.Close()

	_, err = io.Copy(fo, conn)
	handleError(err)
}

func execInput(input string) string {
	// Remove the newline character.
	input = strings.TrimSuffix(input, "\n")

	cmd, err := exec.Command("bash", "-c", input).Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(cmd[:])
}

func validateMailAddress(address string) {
	_, err := mail.ParseAddress(address)
	if err != nil {
		CONFIG.CLIENTLOG.Println("Invalid Email Address. Proceeding anyway.")

		return
	}
	CONFIG.CLIENTLOG.Println("Email Verified. True.")
}

func setReadDeadLine(conn *tls.Conn) {
	err := conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		CONFIG.CLIENTLOG.Panic("SetReadDeadline failed:", err)
	}
}

func setWriteDeadLine(conn *tls.Conn) {
	err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		CONFIG.CLIENTLOG.Panic("SetWriteDeadline failed:", err)
	}
}

func dialReDial(serviceID string, config *tls.Config) *tls.Conn {
	reDial := 0
	for ok := true; ok; ok = reDial < 5 {
		conn, err := tls.Dial("tcp", serviceID, config)
		reDial++
		if err != nil {
			CONFIG.CLIENTLOG.Println("Error: ", err)
			CONFIG.CLIENTLOG.Println("Could not establish connection. Retrying in 5 seconds....")
			time.Sleep(time.Second * 5)
			continue
		}
		CONFIG.CLIENTLOG.Println("Connected to: ", strings.Split(conn.RemoteAddr().String(), ":")[0])
		state := conn.ConnectionState()
		for _, v := range state.PeerCertificates {
			CONFIG.CLIENTLOG.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
			CONFIG.CLIENTLOG.Println(v.Subject)
		}

		CONFIG.CLIENTLOG.Println("client: handshake: ", state.HandshakeComplete)
		return conn

	}
	CONFIG.CLIENTLOG.Println("Timout. Could not reach server. Exiting....")
	os.Exit(1)
	return nil //will never reach this
}
func LoginStatus() bool {
	var user t.LoggedUser
	val, ok := os.LookupEnv("GOSHELLY_ACCESS_TOKEN")

	checkTrue("Access Token exists. Checking validity.", "Token does not exist. Try logging in again.", ok || val=="" )
	if val == "" { //might be redunddant but good to have
		fmt.Println("Invalid Token.")
		return false
	}
	user.ACCESSTOKEN = val
	val, ok = os.LookupEnv("GOSHELLY_ACCESS_EMAIL")
	checkTrue("Identity exists.", "Identity does not exist. ", ok)
	if val == "" || !ok { //might be redunddant but good to have

		fmt.Println("Enter the email associated with your GoShelly Account.")
		fmt.Scanf("%s", &user.EMAIL)
	}
	msg, _ := SendPOST("/users/auth/", user)

	return msg == "Credentials=Valid." //don't like this as im comparing message sent back in json. Comparing status code
	//is better but SendPOST function does not allow sending back resp *http.Response obj as it is being used at multiple places.
}
func getJsonBodyLogin(resp *http.Response) (string, string) {
	var msg t.LogSuccess
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(resp.StatusCode, "Could not read response.")
		return "", ""
	}
	json.Unmarshal(body, &msg)
	fmt.Println("")
	return msg.MESSAGE, msg.TOKEN
}
func GetCredentials(mode int) (string, string, []byte) {
	NAME, EMAIL := "", ""
	switch mode {
	case 1:
		fmt.Printf("Enter your name: ")
		fmt.Scanf("%s", &NAME)
	}

	temp := true
	for ok := true; ok; ok = temp {

		fmt.Printf("Enter email address: ")
		fmt.Scanf("%s", &EMAIL)

		if !validateEMailAddress(EMAIL) {
			fmt.Println("Incorrect email address. Try again.")
			continue
		}
		temp = false
	}
	fmt.Printf("Enter a password: ")
	tmpPass, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("Cannot read from terminal. Try again later.")
		os.Exit(1)
	}
	fmt.Printf("\n.....\n")
	return NAME, EMAIL, tmpPass

}

func genCert() {

	CONFIG.CLIENTLOG.Println("Generating SSL Certificate.")
	validateMailAddress(CONFIG.SSLEMAIL)
	_, err := exec.Command("/bin/bash", "./scripts/certGen.sh", CONFIG.SSLEMAIL).Output()

	if err != nil {
		CONFIG.CLIENTLOG.Printf("Error generating SSL Certificate: %s\n", err)
		os.Exit(1)
	}
}

var CONFIG t.Config

func SendPOST(POSTURL string, user interface{}) (string, string) {
	body, _ := json.Marshal(user)
	resp, err := http.Post(POSTURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Service offline.")
		os.Exit(0)
	}
	return getJsonBodyLogin(resp)
}

func checkTrue(promptTrue, promptFalse string, check bool) {
	if !check {
		fmt.Println(promptFalse)
	} else {
		fmt.Println(promptTrue)
	}
}

func SaveLoginResult(TOKEN, EMAIL string) {

	switch TOKEN {
	case "":
		return
	default:
		fmt.Println("Warning. Your access token and identiy for this session will be stored as an environment variable.")
		os.Setenv("GOSHELLY_ACCESS_TOKEN", TOKEN)
		_, ok := os.LookupEnv("GOSHELLY_ACCESS_TOKEN")
		checkTrue("Token Check=True", "Token failed to save. Try logging in again.", ok)

		os.Setenv("GOSHELLY_ACCESS_EMAIL", EMAIL)
		_, ok = os.LookupEnv("GOSHELLY_ACCESS_EMAIL")
		checkTrue("Email Check=True", "Identity failed to save, you may be prompted to enter your Email at runtime.", ok)
	}

}

func StartClient(HOST string, PORT string, SSLEMAIL string, logmax int) {

	CONFIG.HOST = HOST
	CONFIG.PORT = PORT
	CONFIG.SSLEMAIL = PORT
	CONFIG.MAXLOGSTORE = logmax
	CONFIG.LOGNAME = "./logs/" + "GoShelly" + "-" + time.Now().Format(time.RFC1123) + ".log"
	os.MkdirAll("./logs/", os.ModePerm)
	clientfile, err := os.OpenFile(CONFIG.LOGNAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Client log open error: %s. No logs for this session available. ", err)
		CONFIG.CLIENTLOG = log.New(os.Stdout, "", log.LstdFlags)
	} else {
		CONFIG.CLIENTLOG = log.New(clientfile, "", log.LstdFlags)
		defer clientfile.Close()
	}
	genCert()

	cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
	if err != nil {
		CONFIG.CLIENTLOG.Println("Could not load SSL Certificate. Exiting...")
		return
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn := dialReDial(CONFIG.HOST+":"+CONFIG.PORT, &config)
	defer conn.Close()

	for {
		buffer := make([]byte, 1024)
		setReadDeadLine(conn)
		_, err := conn.Read(buffer)
		if err != nil {
			CONFIG.CLIENTLOG.Println("Checking status.")
			if err == io.EOF {
				CONFIG.CLIENTLOG.Println("All commands ran successfully. Returning exit success.")
				logClean("./logs/")
				fmt.Println("Exit Success.")
				returnLog()
				os.Exit(0)
			}
		}
		sDec, _ := base64.StdEncoding.DecodeString(string(buffer[:]))
		CONFIG.CLIENTLOG.Println("EXECUTE: ", string(sDec))
		resp := execInput(string(sDec))
		time.Sleep(time.Second)
		encodedResp := base64.StdEncoding.EncodeToString([]byte(resp))
		CONFIG.CLIENTLOG.Println("RES:\n", resp)
		setWriteDeadLine(conn)
		_, err = conn.Write([]byte(encodedResp))
		if err != nil {
			CONFIG.CLIENTLOG.Println("Write Error. Exiting. Internal error or server disconnected. Exiting...")
			return
		}
		time.Sleep(time.Second)
		buffer = nil
	}
}
func validateEMailAddress(address string) bool {
	_, err := mail.ParseAddress(address)
	return err == nil
}
