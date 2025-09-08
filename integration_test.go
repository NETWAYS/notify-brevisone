package main

import (
	"os"
	"testing"
)

func TestConfig_Run_Integration(t *testing.T) {
	var (
		gateway  = os.Getenv("BREVIS_INTEGRATION_GATEWAY")
		username = os.Getenv("BREVIS_INTEGRATION_USERNAME")
		password = os.Getenv("BREVIS_INTEGRATION_PASSWORD")
		target   = os.Getenv("BREVIS_INTEGRATION_TARGET")
		ring     = os.Getenv("BREVIS_INTEGRATION_RING") != ""
	)

	if gateway == "" || username == "" || password == "" || target == "" {
		t.Skip("You need to set BREVIS_INTEGRATION_* environment for integration tests")
		return
	}

	c := &Config{
		gateway:          gateway,
		target:           target,
		targetType:       "number",
		ring:             ring,
		username:         username,
		password:         password,
		insecure:         true,
		checkState:       "WARNING",
		checkOutput:      "test notification for notify-brevisone",
		hostName:         "brevis.one",
		serviceName:      "notify",
		author:           "icingaadmin",
		comment:          "comment",
		date:             "",
		notificationType: "CUSTOM",
	}

	err := c.Run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
