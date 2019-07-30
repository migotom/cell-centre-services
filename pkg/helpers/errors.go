package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// AssertErrors checks and compares error with expected one.
func AssertErrors(t *testing.T, err error, expected string) {
	if expected != "" {
		assert.EqualError(t, err, expected)
	} else {
		if err != nil {
			t.Errorf("not expected err %v", err)
		}
	}
}
