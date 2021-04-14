package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/NETWAYS/go-check"
	"io/ioutil"
	"net/http"
	"regexp"
)

const readme = `Notifications via a Brevis.One gateway.
Sends SMS or rings at a given number`

type Config struct {
	gateway	string
	target	string
	targetType string
	ring bool
	username string
	password string
	skipVerify bool
	checkState string
	checkOutput string
	notificationAuthor string
	hostName string
	serviceName string
	comment string
	date string
	notificationType string
}

func main() {
	defer check.CatchPanic()

	plugin := check.NewConfig()
	plugin.Name = "notify-brevisOne"
	plugin.Readme = readme
	plugin.Timeout = 30
	plugin.Version = "0.1"

	config := &Config{}

	// BrevisOne related configuration
	plugin.FlagSet.StringVarP(&config.gateway, "gateway", "g", "", "IP/Address of the Brevis.one gateway (required)")
	plugin.FlagSet.StringVarP(&config.target, "target", "", "", "Contact name (Group or contact) or phone number (required)")
	// TODO:Detect target type automatically?
	plugin.FlagSet.StringVarP(&config.targetType, "targetType", "", "number", "Type of the contact, may be one of: number, contact or contactgroup")
	plugin.FlagSet.BoolVarP(&config.ring, "ring", "r", false, "Ring mode (optional, if not set, send SMS)")
	plugin.FlagSet.StringVarP(&config.username, "username", "u", "", "API user name (required)")
	plugin.FlagSet.StringVarP(&config.password, "password", "p", "", "API user password (required)")
	plugin.FlagSet.BoolVarP(&config.skipVerify, "skipTlsVerify", "", false, "Skip verification of the TLS certificates (is needed for the default self signed certificate of the Brevis.One device")

	// Message configuration
	plugin.FlagSet.StringVar(&config.checkState, "checkresult", "s",  "Return code of the host/service check (required)")
	plugin.FlagSet.StringVarP(&config.checkOutput, "output", "o", "", "Output of the host/service check (required)")
	plugin.FlagSet.StringVarP(&config.notificationAuthor, "notificationAuthor", "a", "", "Author of the notification (optional)")
	plugin.FlagSet.StringVarP(&config.hostName, "host", "", "", "Name of the host object (required)")
	//hostDisplayName := plugin.FlagSet.StringP("hostDisplay", "", "", "Display name of the host object (optional)")
	plugin.FlagSet.StringVarP(&config.serviceName, "service", "", "", "Name of the service object (required for service notifications)")
	//serviceDisplayName := plugin.FlagSet.StringP("serviceDisplay", "", "", "Display name of the service object (optional)")
	plugin.FlagSet.StringVarP(&config.comment, "comment", "", "", "Notification comment (optional)")
	plugin.FlagSet.StringVarP(&config.date, "date", "", "", "Notification date")
	plugin.FlagSet.StringVarP(&config.notificationType, "type", "", "", "Notification type (e.g. Problem, Recovery, etc.")

	// Parsing the arguments
	plugin.ParseArguments()

	if config.username == "" {
		check.ExitError(errors.New("Username not set"))
	}

	if config.password == "" {
		check.ExitError(errors.New("Password not set"))
	}

	if config.gateway == "" {
		check.ExitError(errors.New("Gateway IP/Address not set"))
	}

	if config.target == "" {
		check.ExitError(errors.New("Contact not set"))
	}

	if config.targetType != "number" && config.targetType != "contact" && config.targetType != "contactgroup" {
		check.ExitError(errors.New("Not a valid targetType"))
	}

	if config.hostName == "" {
		check.ExitError(errors.New("Host name must be set"))
	}

	if config.checkState == "" {
		check.ExitError(errors.New("Check result (return code) must be set"))
	}

	msg := ""
	if config.date != "" {
		msg += config.date + "/"
	}
	if config.notificationType != "" {
		msg += config.notificationType + ": "
	}

	if config.serviceName != "" {
		// This is a service notification
		msg += fmt.Sprintf("Srvc:%s @ %s - %s", config.serviceName, config.hostName, config.checkState)
	} else {
		msg += fmt.Sprintf("Hst:%s - %s", config.hostName, config.checkState)
	}

	msg += " - " + config.checkOutput

	remainingSymbols := 160 - len(msg)
	if remainingSymbols < 0 {
		// Gotta cut it :-|
		msg = msg[0:159]
	}

	if config.comment != "" && remainingSymbols >= (len(config.comment)+1) {
		msg += "\n" + config.comment

		if config.notificationAuthor != "" && ((159 - len(msg)) < len(config.notificationAuthor)) {
			msg += "\n" + config.notificationAuthor
		}
	}

	var tlsConf tls.Config
	if config.skipVerify {
		tlsConf.InsecureSkipVerify = true
	}

	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tlsConf}}

	var baseUrl string = "https://" + config.gateway + "/api/"

	// ========================
	// Get authentication token
	// ========================
	signinCreds := []byte("{\"username\": \"" + config.username + "\",\"password\": \"" + config.password + "\"}")
	req, err := http.NewRequest("POST", baseUrl+"signin", bytes.NewBuffer(signinCreds))
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

	recipients := `[{"to":"` + config.target + `","target":"` + config.targetType + `"}]`
	text := `"text":"` + msg + `"`
	provider := `"provider":"sms"`
	providerType := `"type":"default"`
	messageBody := fmt.Sprintf(`{"recipients":%s,%s,%s,%s}`, recipients, text, provider, providerType)

	//fmt.Printf("messageBody: %s\n", messageBody)

	req, err = http.NewRequest("POST", baseUrl+"messages", bytes.NewBuffer([]byte(messageBody)))
	if err != nil {
		check.ExitError(err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authToken)

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

	if config.ring {
		// Additional request to ring
		providerType = `"type":"ring"`
		messageBody = fmt.Sprintf(`{"recipients":%s,%s,%s,%s}`, recipients, text, provider, providerType)

		req, err = http.NewRequest("POST", baseUrl+"messages", bytes.NewBuffer([]byte(messageBody)))
		if err != nil {
			check.ExitError(err)
		}

		req.Header.Add("accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+authToken)

		resp, err = client.Do(req)
		if err != nil {
			check.ExitError(err)
		}

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			check.ExitError(err)
		}
		if resp.StatusCode != 200 {
			fmt.Printf("Send message Status code: %d\n", resp.StatusCode)
			fmt.Printf("Send message Return body: %s\n", body)
			check.ExitError(errors.New("Failed to ring"))
		}
		resp.Body.Close()
	}

	//fmt.Printf("Return from send message: %s\n", body)

	check.Exit(check.OK, "done")
}
