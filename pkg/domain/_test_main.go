package domain

import (
	"os"
	"testing"
)

func TestMain(t *testing.M) {
	code := t.Run()
	os.Exit(code)
}
