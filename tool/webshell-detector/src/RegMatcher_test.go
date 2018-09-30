package WebshellDetector

import (
	"bytes"
	"log"
	"testing"
)

func Test_RegMatcher_DefaultRegMatcher(t *testing.T) {
	regMatcher := defaultRegMatcher()
	if len(regMatcher.regData) == 0 {
		t.Error("FAILED - Test_RegMatcher_DefaultRegMatcher")
	}
}
func Test_RegMatcher_NewRegMatcher(t *testing.T) {
	regMatcherStream, err := Asset("static/config/rules.conf")
	if err != nil {
		log.Fatal(err)
	}
	regMatcherStreamReader := bytes.NewReader(regMatcherStream)
	regMatcher, err := newRegMatcher(regMatcherStreamReader)
	if len(regMatcher.regData) == 0 || err != nil {
		t.Error("FAILED - Test_RegMatcher_NewRegMatcher")
	}
}

func Test_RegMatcher_IsMatched(t *testing.T) {
	regMatcher := defaultRegMatcher()
	if regMatcher.IsMatched([]byte("aaaaaaaaaaaa")) > 0 {
		t.Error("FAILED - Test_RegMatcher_IsMatched - misjudge")
	}
	if regMatcher.IsMatched([]byte("eval($_GET[\"A\"])")) == 0 {
		t.Error("FAILED - Test_RegMatcher_IsMatched - missing")
	}
}
