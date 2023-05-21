package auth

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/tau-OS/xenon/daemon/storage"
)

var JWT_CLAIMS = jwt.MapClaims{}
var l = log.New(os.Stderr, "[auth] ", log.LstdFlags)

func failTo(msg string, err error) {
	if err != nil {
		l.Fatalln("FAIL to " + msg + ": " + err.Error())
	}
}

func LogIn() {
	accessToken := storage.GetKey("token") // access token

	if accessToken == "" {
		initToken()
		return
	}

	ack(accessToken)
}

func ack(accessToken string) {
	req, err := http.NewRequest("GET", "https://sync.fyralabs.com/ack", nil)
	failTo("acknoledge session", err)

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, JWT_CLAIMS).SignedString(accessToken)
	failTo("create JWT", err)

	req.Header.Add("Authentication", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	failTo("connect to server", err)

	if resp.StatusCode != 200 {
		l.Printf("Acknoledgement failed: Status code %d\n", resp.StatusCode)
		l.Println("Resetting token...")
		storage.SetKey("token", "")
		LogIn()
	}
}

const initTokenPrompt = `
┌─────────────────────────────────────────────┐
│ You are not signed in. On a browser, go to: │
│                                             │
│ ==>  https://sync.fyralabs.com/sign-in  <== │
│                                             │
│ Done? Paste the authentication token below. │
└─────────────────────────────────────────────┘`

func initToken() {
	l.Println(initTokenPrompt)
	l.Print("Press ENTER afterwards: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	failTo("read input", err)
	inputs := strings.Split(input, ".")
	token := inputs[0]
	ack(token)
	// if we get to here, the token is valid :3
	storage.SetKey("token", token)
}
