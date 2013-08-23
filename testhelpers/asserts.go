package testhelpers

import (
    "testing"
)

func AssertString(t *testing.T, expected string, actual string) {
    if (expected != actual) {
        t.Errorf("expected %s != actual %s", expected, actual)
    }
}

func AssertInt(t *testing.T, expected int, actual int) {
    if (expected != actual) {
        t.Errorf("expected %s != actual %s", expected, actual)
    }
}

func AssertStringArray(t *testing.T, expected []string, actual []string) {
	if len(expected) != len(actual) {
		t.Errorf("expected and actual not the same length")
		return
	}
	for i, _ := range(expected) {
		if expected[i] != actual[i] {
			t.Errorf("expected[%d] %s != actual[%d] %s", i, expected[i], i, actual[i])
			return
		}
	}
}

func AssertStringMap(t* testing.T, expected map[string]string, actual map[string]string) {
	if len(expected) != len(actual) {
		t.Errorf("expected and actual not the same length")
		return
	}
	for i, _ := range(expected) {
		_, ok := actual[i]
		if !ok {
			t.Errorf("expected[%d] not in actual", i)
		}
		if expected[i] != actual[i] {
			t.Errorf("expected[%d] %s != actual[%d] %s", i, expected[i], i, actual[i])
			return
		}
	}
}
