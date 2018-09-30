package WebshellDetector

import (
	"github.com/CyrusF/libsvm-go"
)

/*
WebshellDetector - Refactor version 1
Date	0814
Author	Cyrus
Intro	Use operation serial array
        Convert array to map for model prediction
*/

type opSerial struct {
	data []float64
}

func (self opSerial) Predict(model *libSvm.Model) float64 {
	var inputData = make(map[int]float64)
	for k, v := range self.data {
		if v != 0 {
			inputData[k+1] = v // model start at 1
		}
	}
	_, result := model.PredictValues(inputData)
	return result[0]
}
