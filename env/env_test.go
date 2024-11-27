package env

import (
	"testing"
)

func TestReadEnv(t *testing.T) {
	tmp := t.TempDir()
	ConfigPath = tmp

	_, err := ReadEnv()

	if err != nil {
		t.Errorf("expected no errors, but one was found \"%s\"", err.Error())
	}
}
