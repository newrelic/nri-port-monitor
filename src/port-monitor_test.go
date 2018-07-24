package main

import (
	"testing"
)

func TestPopulateInventory(t *testing.T) {
	// Insert here the logic for your tests
	actual := 2
	expected := 2
	if actual != expected {
		t.Errorf("PopulateInventory was incorrect, got: %d, expected: %d", actual, expected)
	}
}

func TestPopulateMetrics(t *testing.T) {
	// Insert here the logic for your tests
	actual := "foo"
	expected := "foo"
	if actual != expected {
		t.Errorf("PopulateMetrics was incorrect, got: %s, expected: %s", actual, expected)
	}
}
