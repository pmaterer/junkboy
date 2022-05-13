package junkboy

import "testing"

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

func assertEqual(t *testing.T, expected interface{}, got interface{}) {
	t.Helper()
	if expected != got {
		t.Fatalf("Not equal: \n"+
			"expected: %+v\n"+
			"got: %+v", expected, got)
	}
}
