package main

import (
	"github.com/NETWAYS/go-check"
	"testing"
)

func TestConfig_Bind(t *testing.T) {
	plugin := check.NewConfig()
	config := &Config{}
	config.BindArguments(plugin.FlagSet)
	plugin.ParseArguments()

	if !plugin.FlagSet.Parsed() {
		t.Fatalf("expected CLI be parsed, got %v", plugin.FlagSet.Parsed())
	}
}

func TestConfig_Validate(t *testing.T) {
	c := &Config{
		gateway:          "brevisone.local",
		username:         "admin",
		password:         "admin",
		target:           "+491710000",
		targetType:       "number",
		checkState:       "DOWN",
		hostName:         "HOST",
		notificationType: "PROBLEM",
	}

	err := c.Validate()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestConfig_FormatMessage(t *testing.T) {
	c := &Config{
		checkState:  "DOWN",
		checkOutput: "Stuff is broken!",
		author:      "",
		hostName:    "HOST",
		//date:               "2014-01-03 11:23:08 +0000",
		notificationType: "PROBLEM",
	}

	// Host Problem Notification
	expected := "PROBLEM: HOST - DOWN\r\nStuff is broken!"
	if expected != c.FormatMessage() {
		t.Fatalf("expected %v, got %v", expected, c.FormatMessage())
	}

	// Service Problem Notification
	c.checkState = "CRITICAL"
	c.serviceName = "SERVICE"

	expected = "PROBLEM: SERVICE @ HOST - CRITICAL\r\nStuff is broken!"

	if expected != c.FormatMessage() {
		t.Fatalf("expected %v, got %v", expected, c.FormatMessage())
	}

	// With comment
	c.notificationType = "CUSTOM"
	c.author = "icingaadmin"
	c.comment = "ok for now"

	expected = "CUSTOM: SERVICE @ HOST - CRITICAL\r\n\"ok for now\" by icingaadmin\r\nStuff is broken!"

	if expected != c.FormatMessage() {
		t.Fatalf("expected %v, got %v", expected, c.FormatMessage())
	}

	// With a long message cut off
	c.checkOutput = "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been" +
		"the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and" +
		"scrambled it to make a type specimen book."

	expected = "CUSTOM: SERVICE @ HOST - CRITICAL\r\n\"ok for now\" by icingaadmin\r\n" +
		"Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has b..."

	if expected != c.FormatMessage() {
		t.Fatalf("expected %v, got %v", expected, c.FormatMessage())
	}
}
