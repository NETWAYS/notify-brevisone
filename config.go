package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"

	"github.com/spf13/pflag"
)

const SmsLength = 160

type Config struct {
	gateway          string
	target           string
	targetType       string
	ring             bool
	username         string
	password         string
	insecure         bool
	checkState       string
	checkOutput      string
	author           string
	hostName         string
	serviceName      string
	comment          string
	date             string
	notificationType string
	doNotUseTLS      bool
	useLegacyHttpApi bool
}

func (c *Config) BindArguments(fs *pflag.FlagSet) {
	// Basic connection settings
	fs.StringVarP(&c.gateway, "gateway", "g", "", "IP/address of the brevis.one gateway (required)")
	fs.StringVarP(&c.username, "username", "u", "", "API user name (required)")
	fs.StringVarP(&c.password, "password", "p", "", "API user password (required)")
	fs.BoolVar(&c.insecure, "insecure", false,
		"Skip verification of the TLS certificates (is needed for the default self signed certificate, default false)")
	fs.BoolVar(&c.doNotUseTLS, "doNotUseTLS", false, "Do NOT use TLS to connect to the gateway (default false)")
	fs.BoolVar(&c.useLegacyHttpApi, "useLegacyHttpApi", false, "Use old HTTP API (required on older firmware versions, default false)")

	// Where to send the message to
	fs.StringVarP(&c.target, "target", "T", "", "Target contact, group or phone number (required)")
	fs.StringVar(&c.targetType, "target-type", "number", "Target type, one of: number, contact or contactgroup")
	fs.BoolVarP(&c.ring, "ring", "R", false, "Add ring mode (also ring the target after sending SMS)")

	// Notification data
	fs.StringVarP(&c.notificationType, "type", "", "", "Icinga $notification.type$ (required)")
	fs.StringVarP(&c.hostName, "host", "H", "", "Icinga $host.name$ (required)")
	fs.StringVarP(&c.serviceName, "service", "S", "", "Icinga $service.name$ (required for service notifications)")
	fs.StringVarP(&c.checkState, "state", "s", "", "Icinga $host.state$ or $service.state$ (required)")
	fs.StringVarP(&c.checkOutput, "output", "o", "", "Icinga $host.output or $service.output$ (required)")
	fs.StringVarP(&c.comment, "comment", "C", "", "Icinga $notification.comment$ (optional)")
	fs.StringVarP(&c.author, "author", "a", "", "Icinga $notification.author$ (optional)")
	fs.StringVar(&c.date, "date", "", "Icinga $icinga.long_date_time$ (optional)")
}

func (c *Config) Validate() error {
	if c.gateway == "" {
		return errors.New("gateway IP/address not set")
	}

	if c.username == "" {
		return errors.New("username not set")
	}

	if c.password == "" {
		return errors.New("password not set")
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

	if c.notificationType == "" {
		return errors.New("notification type must be set")
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
		msg += fmt.Sprintf("%s @ %s - %s", c.serviceName, c.hostName, c.checkState)
	} else {
		msg += fmt.Sprintf("%s - %s", c.hostName, c.checkState)
	}

	if c.comment != "" {
		msg += fmt.Sprintf("\r\n\"%s\"", c.comment)

		if c.author != "" {
			msg += fmt.Sprintf(` by %s`, c.author)
		}
	}

	if c.checkOutput != "" {
		msg += "\r\n" + c.checkOutput
	}

	// Cut off text longer than a single message
	if len(msg) > SmsLength {
		msg = msg[0:SmsLength-4] + "..."
	}

	return
}

func (c *Config) Run() (err error) {
	// Setup API client
	api := NewApiClient(c.gateway)

	// Update client to allow insecure when requested
	if c.insecure {
		api.Client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // nolint:gosec
			},
		}
	}

	api.UseTls = !c.doNotUseTLS

	if c.useLegacyHttpApi {
		err := api.DoLegacyReqest(c.targetType,
			c.target,
			c.FormatMessage(),
			c.username,
			c.password)
		if err != nil {
			return err
		}
	} else {

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
			return err
		}

		// TODO: check response content

		// Additional request to ring after sending SMS
		if c.ring {
			message.Type = "ring"

			// TODO: check response content
			response, err = api.DoRequest("messages", message)
			if err != nil {
				err = fmt.Errorf("ringing failed: %s - %w", response, err)
				return err
			}
		}
	}

	return
}
