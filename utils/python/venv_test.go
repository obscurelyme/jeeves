package python

import "testing"

func TestFormatPythonVersion(t *testing.T) {
	officialVersion := formatPythonVersion("Python 3.10.12")

	if officialVersion != "python3.10" {
		t.Errorf("Expected python version to be %s, got %s", "python3.10", officialVersion)
	}
}

func TestFormatPythonVersion2(t *testing.T) {
	officialVersion := formatPythonVersion("Python 3.10")

	if officialVersion != "python3.10" {
		t.Errorf("Expected python version to be %s, got %s", "python3.10", officialVersion)
	}
}
