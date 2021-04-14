package main

import (
	"bytes"
	"net/http"
	"fmt"
	"github.com/NETWAYS/go-check"
	"io/ioutil"
	"errors"
	"crypto/tls"
	"regexp"
)

const readme = `Notifications via a Brevis.One gateway.
Sends SMS or rings at a given number`

func main() {
	defer check.CatchPanic()

    plugin := check.NewConfig()
    plugin.Name = "notify-brevisOne"
    plugin.Readme = readme
    plugin.Timeout = 30
	plugin.Version = "0.1"

	// targetNumber
	//ring := plugin.FlagSet.BoolP("ring", "r", false, "Ring mode (optional, if not set, send SMS)")
	gateway := plugin.FlagSet.StringP("gatewayIP", "g", "", "IP/Address of the Brevis.one gateway (required)")
	contact := plugin.FlagSet.StringP("contact", "c", "", "Contact name (Group or contact) or number (required)")
	username := plugin.FlagSet.StringP("username", "u", "", "API user name (required)")
	password := plugin.FlagSet.StringP("password", "p", "", "API user password (required)")

	checkState := plugin.FlagSet.IntP("checkresult", "s", -1, "Return code of the host/service check (required)")
	checkOutput := plugin.FlagSet.StringP("output", "o", "", "Output of the host/service check (required)")
	notificationAuthor := plugin.FlagSet.StringP("notificationAuthor", "a", "", "Author of the notification (optional)")
	hostName := plugin.FlagSet.StringP("host", "", "", "Name of the host object (required)")
	//hostDisplayName := plugin.FlagSet.StringP("hostDisplay", "", "", "Display name of the host object (optional)")
	serviceName := plugin.FlagSet.StringP("service", "", "", "Name of the service object (required for service notifications)")
	//serviceDisplayName := plugin.FlagSet.StringP("serviceDisplay", "", "", "Display name of the service object (optional)")
	comment := plugin.FlagSet.StringP("comment", "", "", "Notification comment (optional)")
	date := plugin.FlagSet.StringP("data", "", "", "Notification data")

	notificationType := plugin.FlagSet.StringP("type", "", "", "Notification type (e.g. Problem, Recovery, etc.")


	// Parsing the arguments
	plugin.ParseArguments()

	if *username == "" {
		check.ExitError(errors.New("Username not set"))
	}

	if *password == "" {
		check.ExitError(errors.New("Password not set"))
	}

	if *gateway == "" {
		check.ExitError(errors.New("Gateway IP/Address not set"))
	}

	if *contact == "" {
		check.ExitError(errors.New("Contact not set"))
	}

	if *hostName == "" {
		check.ExitError(errors.New("Host name must be set"))
	}

	if *checkState == -1 {
		check.ExitError(errors.New("Check result (return code) must be set"))
	}

	msg := ""
	if *date != "" {
		msg += *date + "/"
	}
	if *notificationType != "" {
		msg += *notificationType + ": "
	}

	if *serviceName != "" {
		// This is a service notification
		msg += fmt.Sprintf("Srvc:%s on %s is %s", *serviceName, *hostName, check.StatusText(*checkState))
	} else {
		longState := ""
		if *checkState == 0 {
			longState = "OK"
		} else {
			longState = "DOWN"
		}
		msg += fmt.Sprintf("Hst:%s is %s", *hostName, longState)
	}

	msg += "\n" + *checkOutput

	remainingSymbols := 160 - len(msg)
	if remainingSymbols < 0 {
		// Gotta cut it :-|
		msg = msg[0:159]
	}

	if *comment != "" && remainingSymbols >= (len(*comment) + 1) {
		msg += "\n" + *comment

		if *notificationAuthor != "" && ( (159 - len(msg)) < len(*notificationAuthor) ) {
			msg += "\n" + *notificationAuthor
		}
	}

	var tlsConf tls.Config
	tlsConf.InsecureSkipVerify = true

	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tlsConf}}

	var baseUrl string = "https://" + *gateway + "/api/"

	// Get authentication token
	signinCreds := []byte("{\"username\": \"" + *username + "\",\"password\": \"" + *password + "\"}")
	req, err := http.NewRequest("POST", baseUrl + "signin", bytes.NewBuffer(signinCreds))
	if err != nil {
		check.ExitError(err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		check.ExitError(err)
	}

	if resp.StatusCode != 200 {
		check.ExitError(errors.New("Could not create authentication token"))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		check.ExitError(err)
	}
	resp.Body.Close()

	// body now should look like that {"jwt":"AUTH_TOKEN","expireAt":SOME_NUMBER}
	// we need AUTH_TOKEN
	re := regexp.MustCompile("\"jwt\":\"(?P<authToken>[[:graph:]]+?)\"")
	//authPart := re.Find(body)
	authPart := re.FindSubmatch(body)
	if authPart == nil {
		//fmt.Println(string(body[:]))
		check.ExitError(errors.New("No Auth Token in body"))
	}

	authToken := string(authPart[1][:])

	// ============
	// Send message
	// ============

	recipients := `[{"to":"` + *contact + `","target":"number"}]`
	text := `"text":"` + msg + `"`
	provider := `"provider":"sms"`
	providerType := `"type":"default"`
	messageBody := fmt.Sprintf(`{"recipients":%s,%s,%s,%s}` , recipients, text, provider, providerType)

	//fmt.Printf("messageBody: %s\n", messageBody)

	req, err = http.NewRequest("POST", baseUrl + "messages", bytes.NewBuffer([]byte(messageBody)))
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " + authToken)
	//fmt.Println(req)
	if err != nil {
		check.ExitError(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		check.ExitError(err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		check.ExitError(err)
	}
	resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Send message Status code: %d\n", resp.StatusCode)
		fmt.Printf("Send message Return body: %s\n", body)
		check.ExitError(errors.New("Could not send message"))
	}


	//fmt.Printf("Return from send message: %s\n", body)

	check.Exit(check.OK,"done")
}
