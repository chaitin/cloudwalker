package WebshellDetector

import (
	"testing"
)

func Test_BinData_opSerialStream(t *testing.T) {
	_, err := Asset("static/model-latest/OpSerial.model")
	if err != nil {
		t.Log(err)
		t.Error("FAILED - Test_BinData - opSerialStream")
	}
}
func Test_BinData_processorStream(t *testing.T) {
	_, err := Asset("static/model-latest/Processor.model")
	if err != nil {
		t.Log(err)
		t.Error("FAILED - Test_BinData - processorStream")
	}
}
func Test_BinData_wordStream(t *testing.T) {
	_, err := Asset("static/model-latest/Words.model")
	if err != nil {
		t.Log(err)
		t.Error("FAILED - Test_BinData - wordStream")
	}
}
func Test_BinData_hashStateBytes(t *testing.T) {
	_, err := Asset("static/model-latest/hashState.json")
	if err != nil {
		t.Log(err)
		t.Error("FAILED - Test_BinData - hashStateBytes")
	}
}
func Test_BinData_statStateBytes(t *testing.T) {
	_, err := Asset("static/config/statState.json")
	if err != nil {
		t.Log(err)
		t.Error("FAILED - Test_BinData - statStateBytes")
	}
}
func Test_BinData_regMatcherStream(t *testing.T) {
	_, err := Asset("static/config/rules.conf")
	if err != nil {
		t.Log(err)
		t.Error("FAILED - Test_BinData - regMatcherStream")
	}
}
func Test_BinData_sampleMatcherStream(t *testing.T) {
	_, err := Asset("static/model-latest/SampleHash.txt")
	if err != nil {
		t.Log(err)
		t.Error("FAILED - Test_BinData - sampleMatcherStream")
	}
}
