package WebshellDetector

import (
	"bytes"
	"testing"
)

func Test_PhpAstGen_GetData(t *testing.T) {
	phpAstGen := newPhpAstGenerator(stdin, stdout)
	src := []byte(`<?php echo 'hello world'; ?>`)
	std := []byte(`{"status":"successed","ast":{"kind":132,"flags":0,"lineno":1,"children":[{"kind":282,"flags":0,"lineno":1,"children":{"expr":"hello world"}}]}}` + "\n")
	result, err := phpAstGen.GetData(src)
	if bytes.Compare(result, std) != 0 || err != nil {
		t.Log(bytes.Compare(result, std))
		t.Log(err)
		t.Log(string(result))
		t.Log(string(std))
		t.Error("FAILED - Test_PhpAstServer_GetData")
	}
}
