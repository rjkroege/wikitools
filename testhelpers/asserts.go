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
