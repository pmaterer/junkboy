package junkboy

import (
	"bytes"
	"testing"
)

func assertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("Unexpected error: \n"+
			"%+v", err)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func assertEqual(t *testing.T, expected, got interface{}) {
	t.Helper()

	if expected != got {
		t.Fatalf("Not equal: \n"+
			"expected: %+v\n"+
			"got: %+v", expected, got)
	}
}

func assertBytesEqual(t *testing.T, expected, actual []byte) {
	t.Helper()

	n := bytes.Compare(expected, actual)
	if n != 0 {
		t.Fatalf("Not equal: \n"+
			"expected: % x\n"+
			"actual: % x", expected, actual)
	}
}
