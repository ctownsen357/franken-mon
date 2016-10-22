package main

import (
	"testing"
)

func TestGetConfig(t *testing.T) {
	var config Configuration = GetConfig() //explicitly declaring the expected type as part of the test
	if config.Timeout <= 0 {
		t.Fatal("Timout must be greater than 0.", config.Timeout)
	}
}
