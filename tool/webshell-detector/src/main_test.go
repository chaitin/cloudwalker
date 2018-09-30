package WebshellDetector

import (
	"os"
	"testing"

	"../php"
)

var stdin *os.File
var stdout *os.File

func TestMain(m *testing.M) {
	php.Start()
	stdin = php.Stdin
	stdout = php.Stdout

	code := m.Run()

	os.Exit(code)
}
