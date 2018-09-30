package WebshellDetector

import (
	"testing"
)

func Test_OpSerial_Predict(t *testing.T) {
	detector, _ := NewDefaultDetector(stdin, stdout)
	opserial1 := opSerial{[]float64{1, 2, 3}}
	result1 := opserial1.Predict(detector.opSerialModel)
	if result1 > 0 {
		t.Log(result1)
		t.Error("FAILED - Test_OpSerial_Predict - data1")
	}
	opserial2 := opSerial{[]float64{1, 2, 3, 4}}
	result2 := opserial2.Predict(detector.opSerialModel)
	if result2 > 0 {
		t.Log(result2)
		t.Error("FAILED - Test_OpSerial_Predict - data2")
	}
}
