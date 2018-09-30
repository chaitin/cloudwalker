package WebshellDetector

import "testing"

func Test_Words_Predict(t *testing.T) {
	detector, _ := NewDefaultDetector(stdin, stdout)
	words1 := words{[]string{"f", "g", "hello", "world"}}
	result1 := words1.Predict(detector.wordsModel)
	if result1 >= 0.5 {
		t.Log(result1)
		t.Error("data1")
	}
	words2 := words{[]string{"create_function", "assert"}}
	result2 := words2.Predict(detector.wordsModel)
	if result2 <= 0.5 {
		t.Log(result2)
		t.Error("data2")
	}
}
