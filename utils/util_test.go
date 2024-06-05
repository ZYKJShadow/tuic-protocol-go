package utils

import "testing"

func TestGenerateCert(t *testing.T) {
	err := GenerateCert("localhost", "127.0.0.1")
	if err != nil {
		t.Errorf("GenerateCert() error = %v", err)
		return
	}
}
