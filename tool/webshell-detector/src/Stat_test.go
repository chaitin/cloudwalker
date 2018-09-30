package WebshellDetector

import (
	"testing"
)

func Test_Stat_NewStat(t *testing.T) {
	src := []byte("<?php\n$a = 'hello world';\necho %a;\n?>")
	stat := newTextStat(src)
	std := textStat{19, 0.8724939396583133, 5, 57.94256397693845, 48.64864864864865, 16.666666666666664, 0.5, 4.0796786498829745}
	if stat != std {
		t.Error("FAILED - Test_Stat_NewStat")
	}
}

func Test_Stat_GetVector(t *testing.T) {
	src := []byte("<?php\n$a = 'hello world';\necho %a;\n?>")
	statVector := newTextStat(src).GetVector()
	std := []float64{19, 0.8724939396583133, 5, 57.94256397693845, 48.64864864864865, 16.666666666666664, 0.5, 4.0796786498829745}
	for i, _ := range statVector {
		if statVector[i] != std[i] {
			t.Error("FAILED - Test_Stat_NewStat")
		}
	}
}

func Test_Stat_IsAbnormal(t *testing.T) {
	detector, _ := NewDefaultDetector(stdin, stdout)

	src1 := []byte("<?php\n$a = 'hello world';\necho %a;\n?>")
	statResult1 := newTextStat(src1).IsAbnormal(detector.statState)
	if statResult1 != false {
		t.Error("FAILED - Test_Stat_IsAbnormal - data1")
	}

	src2 := []byte("aaaaaaaaaaaaaaaaaaaaa1aaaaaaaaaaaaaaaaaaaa")
	statResult2 := newTextStat(src2).IsAbnormal(detector.statState)
	if statResult2 != true {
		t.Error("FAILED - Test_Stat_IsAbnormal - data2")
	}
}
