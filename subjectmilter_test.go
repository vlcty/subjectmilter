package main

import (
	"testing"
)

func TestContainsBadString(t *testing.T) {
	badstrings := []string{
		"Fuckbuddy",
		"señor",
		"Höhle"}

	// Teststring -> Expected result
	testcases := map[string]bool{
		"=?UTF-8?b?SG9sYSwgc2XDsW9yISBBIGdyZWF0IGdpZnQgYXdhaXRzIHlvdSE=?=": true,
		"„Höhle der Löwen“ System macht Deutsche Bürger reich!":            true,
		"Fuckbuddy gesucht": true,
		"señor\nChang":      true,
		"señ\nor":           true,
		"fuckbuddy gesucht": false}

	milter := &MyFilter{badstrings: badstrings}

	for testcase, expected := range testcases {
		if result := milter.ContainsBadString(testcase); result != expected {
			t.Errorf("Teststring: %q. Expected %t but resulted in %t", testcase, expected, result)
		}
	}
}
