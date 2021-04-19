package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/spf13/pflag"
	"net/http"
)

type Config struct {
	gateway            string
	target             string
	targetType         string
	ring               bool
	username           string
	password           string
	skipVerify         bool
	checkState         string
	checkOutput        string
	notificationAuthor string
	hostName           string
	serviceName        string
	comment            string
	date               string
	notificationType   string
}

func (c *Config) BindArguments(fs *pflag.FlagSet) {
	fs.StringVarP(&c.gateway, "gateway", "g", "", "IP/Address of the Brevis.one gateway (required)")
	fs.StringVarP(&c.target, "target", "", "", "Contact name (Group or contact) or phone number (required)")
	// TODO:Detect target type automatically?
	fs.StringVarP(&c.targetType, "targetType", "", "number",
		"Type of the contact, may be one of: number, contact or contactgroup")
	fs.BoolVarP(&c.ring, "ring", "r", false, "Ring mode (optional, if not set, send SMS)")
	fs.StringVarP(&c.username, "username", "u", "", "API user name (required)")
	fs.StringVarP(&c.password, "password", "p", "", "API user password (required)")
	fs.BoolVarP(&c.skipVerify, "skipTlsVerify", "", false,
		"Skip verification of the TLS certificates (is needed for the default self signed certificate)")

	// Message configuration
	fs.StringVar(&c.checkState, "checkresult", "s", "Return code of the host/service check (required)")
	fs.StringVarP(&c.checkOutput, "output", "o", "", "Output of the host/service check (required)")
	fs.StringVarP(&c.notificationAuthor, "notificationAuthor", "a", "", "Author of the notification (optional)")
	fs.StringVarP(&c.hostName, "host", "", "", "Name of the host object (required)")
	//hostDisplayName := fs.StringP("hostDisplay", "", "", "Display name of the host object (optional)")
	fs.StringVarP(&c.serviceName, "service", "", "", "Name of the service object (required for service notifications)")
	//serviceDisplayName := fs.StringP("serviceDisplay", "", "", "Display name of the service object (optional)")
	fs.StringVarP(&c.comment, "comment", "", "", "Notification comment (optional)")
	fs.StringVarP(&c.date, "date", "", "", "Notification date")
	fs.StringVarP(&c.notificationType, "type", "", "", "Notification type (e.g. Problem, Recovery, etc.")
}

func (c *Config) Validate() error {
	if c.username == "" {
		return errors.New("username not set")
	}

	if c.password == "" {
		return errors.New("password not set")
	}

	if c.gateway == "" {
		return errors.New("gateway IP/address not set")
	}

	if c.target == "" {
		return errors.New("target not set")
	}

	if c.targetType != "number" && c.targetType != "contact" && c.targetType != "contactgroup" {
		return fmt.Errorf("not a valid targetType: %s", c.targetType)
	}

	if c.hostName == "" {
		return errors.New("hostName must be set")
	}

	if c.checkState == "" {
		return errors.New("check result state (return code) must be set")
	}

	return nil
}

func (c *Config) FormatMessage() (msg string) {
	if c.date != "" {
		msg += c.date + "/"
	}

	if c.notificationType != "" {
		msg += c.notificationType + ": "
	}

	if c.serviceName != "" {
		// This is a service notification
		msg += fmt.Sprintf("Srvc:%s @ %s - %s", c.serviceName, c.hostName, c.checkState)
	} else {
		msg += fmt.Sprintf("Hst:%s - %s", c.hostName, c.checkState)
	}

	msg += " - " + c.checkOutput

	remainingSymbols := 160 - len(msg)
	if remainingSymbols < 0 {
		// Gotta cut it :-|
		msg = msg[0:159]
	}

	if c.comment != "" && remainingSymbols >= (len(c.comment)+1) {
		msg += "\n" + c.comment

		if c.notificationAuthor != "" && ((159 - len(msg)) < len(c.notificationAuthor)) {
			msg += "\n" + c.notificationAuthor
		}
	}

	return
}

func (c *Config) Run() (err error) {
	// Setup API client
	api := NewApiClient(c.gateway)

	// Update client to allow insecure when requested
	if c.skipVerify {
		api.Client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // nolint:gosec
			},
		}
	}

	// Login to API
	err = api.Login(c.username, c.password)
	if err != nil {
		return
	}

	// Sending message via API
	message := Message{
		Recipients: []Recipient{{
			To:     c.target,
			Target: c.targetType,
		}},
		Text:     c.FormatMessage(),
		Provider: "sms",
		Type:     "default",
	}

	response, err := api.DoRequest("messages", message)
	if err != nil {
		err = fmt.Errorf("sending message failed: %s - %w", response, err)
		return
	}

	// TODO: check response content

	// Additional request to ring after sending SMS
	if c.ring {
		message.Type = "ring"

		// TODO: check response content
		response, err = api.DoRequest("messages", message)
		if err != nil {
			err = fmt.Errorf("ringing failed: %s - %w", response, err)
			return
		}
	}

	return
}
