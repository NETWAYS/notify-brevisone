package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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

	assert.NoError(t, c.Validate())
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
	assert.Equal(t, "PROBLEM: HOST - DOWN\r\nStuff is broken!", c.FormatMessage())

	// Service Problem Notification
	c.checkState = "CRITICAL"
	c.serviceName = "SERVICE"

	assert.Equal(t, "PROBLEM: SERVICE @ HOST - CRITICAL\r\nStuff is broken!", c.FormatMessage())

	// With comment
	c.notificationType = "CUSTOM"
	c.author = "icingaadmin"
	c.comment = "ok for now"

	assert.Equal(t,
		"CUSTOM: SERVICE @ HOST - CRITICAL\r\n\"ok for now\" by icingaadmin\r\nStuff is broken!",
		c.FormatMessage())

	// With a long message cut off
	c.checkOutput = "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been" +
		"the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and" +
		"scrambled it to make a type specimen book."

	assert.Equal(t,
		"CUSTOM: SERVICE @ HOST - CRITICAL\r\n\"ok for now\" by icingaadmin\r\n"+
			"Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has b...",
		c.FormatMessage())
}
