package app_test

import (
	"os"
	"testing"

	"github.com/mosteligible/go-logreader/client/app"
)

var a = app.NewApp()

func TestApp(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
