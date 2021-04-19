package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_FormatMessage(t *testing.T) {
	c := &Config{
		checkState:         "DOWN",
		checkOutput:        "Stuff is broken!",
		notificationAuthor: "",
		hostName:           "HOST",
		//date:               "2014-01-03 11:23:08 +0000",
		notificationType: "PROBLEM",
	}

	// Host Problem Notification
	assert.Equal(t, "PROBLEM: HOST - DOWN\nStuff is broken!", c.FormatMessage())

	// Service Problem Notification
	c.checkState = "CRITICAL"
	c.serviceName = "SERVICE"

	assert.Equal(t, "PROBLEM: SERVICE @ HOST - CRITICAL\nStuff is broken!", c.FormatMessage())

	// With comment
	c.notificationType = "CUSTOM"
	c.notificationAuthor = "icingaadmin"
	c.comment = "ok for now"

	assert.Equal(t,
		"CUSTOM: SERVICE @ HOST - CRITICAL\n\"ok for now\" by icingaadmin\nStuff is broken!",
		c.FormatMessage())

	// With a long message cut off
	c.checkOutput = "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been" +
		"the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and" +
		"scrambled it to make a type specimen book."

	assert.Equal(t,
		"CUSTOM: SERVICE @ HOST - CRITICAL\n\"ok for now\" by icingaadmin\n"+
			"Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has bee...",
		c.FormatMessage())
}
