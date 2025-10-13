package testhelpers

import (
	"testing"
)

func AssertStringMap(t *testing.T, expected map[string]string, actual map[string]string) {
	if len(expected) != len(actual) {
		t.Errorf("expected and actual not the same length")
		return
	}
	for i := range expected {
		_, ok := actual[i]
		if !ok {
			t.Errorf("expected[%s] not in actual", i)
		}
		if expected[i] != actual[i] {
			t.Errorf("expected[%s] %s != actual[%s] %s", i, expected[i], i, actual[i])
			return
		}
	}
}
