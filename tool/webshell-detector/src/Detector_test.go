package WebshellDetector

import (
	"fmt"
	"testing"
)

func Test_Detector_NewDetector(t *testing.T) {
	detector, err := NewDefaultDetector(stdin, stdout)
	if err != nil {
		t.Error(err.Error())
	}
	if detector.opSerialModel.NrClass() == 0 {
		t.Error("FAILED - Test_Detector_NewDetector - opSerialModel")
	}
	if len(detector.wordsModel.LearningResults) == 0 {
		t.Error("FAILED - Test_Detector_NewDetector - wordsModel")
	}
	if detector.processModel.NrClass() == 0 {
		t.Error("FAILED - Test_Detector_NewDetector - processModel")
	}
	if len(detector.regMatcher.regData) == 0 {
		t.Error("FAILED - Test_Detector_NewDetector - regMatcher")
	}
	if len(detector.hashState.Table) == 0 {
		t.Error("FAILED - Test_Detector_NewDetector - hashState")
	}
	if fmt.Sprintf("%v", detector.statState) != "{{NaN 0.1 NaN NaN 10 NaN 0.001 NaN} {2048 NaN 1024 NaN NaN NaN NaN NaN}}" {
		t.Logf("%v", detector.statState)
		t.Error("FAILED - Test_Detector_NewDetector - statState")
	}
}

func Test_Detector_Predict(t *testing.T) {
	detector, _ := NewDefaultDetector(stdin, stdout)
	src := []byte(`<?php echo f(g()+231*544+555+f(1,2,3+4,5)); ?>`)
	result, err := detector.Predict(src)
	if err != nil {
		t.Error(err.Error())
	}
	if result < 0 {
		t.Error("FAILED - Test_Detector_Predict")
	}
}
