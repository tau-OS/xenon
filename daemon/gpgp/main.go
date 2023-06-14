// gpg processor!!!
package gpgp

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	// gpg "github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/tau-OS/xenon/daemon/storage"
)

var l = log.New(os.Stderr, "[gpgp] ", log.LstdFlags)
var Gpgkey string

func Prep() {
	Gpgkey = storage.Local.GetItem("gpgkey")
	if Gpgkey == "" {
		Gpgkey = setup()
	}

}

func setup() string {
	out, err := exec.Command("gpg", "--list-secret-keys").Output()
	if err != nil {
		l.Fatal("FAIL to list secret keys: " + err.Error())
	}
	rawkeys := strings.Split(string(out), "\n")[2:]
	var keys []string
	for _, line := range rawkeys {
		if strings.HasPrefix(line, " ") {
			// pretty sure it's a key?
			key := strings.TrimSpace(line)
			keys = append(keys, key)
		}
	}
	l.Println("=>> You don't have a GnuPG key chosen. Please select one from below:")
	i := 0
	for _, line := range rawkeys {
		if strings.HasPrefix(line, "sec") {
			i += 1
		} else if i == 0 {
			continue
		}
		if strings.TrimSpace(line) == "" {
			l.Println()
			continue
		}
		l.Printf(" >> %d. "+line, i)
	}
	l.Println()
	l.Println("... The above is the output of `gpg --list-secret-keys`. You can only enter a choice from the above list.")
	l.Println("... If you need to perform any other actions, press Ctrl+C to stop the daemon, then re-run it afterwards.")
	fmt.Print("<<- Enter your choice: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		l.Fatalln("FAIL to read input: " + err.Error())
	}
	i, err = strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		l.Fatalf("FAIL to convert `%s` to int: "+err.Error(), input)
	}
	l.Println("Using " + keys[i-1])
	storage.Local.SetItem("gpgkey", keys[i-1])
	return keys[i-1]
}
