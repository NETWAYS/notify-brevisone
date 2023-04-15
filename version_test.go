package main

import (
	"testing"
)

func TestBuildVersion(t *testing.T) {

	var expected string
	expected = "development\ncommit: HEAD\ndate: latest"

	if expected != buildVersion() {
		t.Error("\nActual: ", buildVersion(), "\nExpected: ", expected)
	}

}
