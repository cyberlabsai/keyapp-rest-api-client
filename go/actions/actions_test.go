package actions

import (
	"testing"
)

func TestValid(t *testing.T) {
	empty := &ActionToken{}

	if empty.Valid() {
		t.Errorf("Empty action token must be invalid")
	}

	actionToken := &ActionToken{
		Token: "fake-token",
	}

	if !actionToken.Valid() {
		t.Errorf("Action token with a filled token field must be valid")
	}
}
