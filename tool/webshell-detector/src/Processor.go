package WebshellDetector

import "os"

/*
WebshellDetector - Refactor version 1
Date	0814
Author	Cyrus, Twice
Intro	Process 8 textStat vectors and 2 ML vectors, final score
*/

type processor struct {
	detector *Detector
	opSerial *opSerial
	words    *words
	stat     *textStat
	callable bool
}

func newProcessor(detector *Detector, ast *ast, stat *textStat) *processor {
	opSerial := ast.GetOpSerial(&detector.hashState)
	words, callable := ast.GetWordsAndCallable()

	var self processor
	self.detector = detector
	self.opSerial = &opSerial
	self.words = &words
	self.stat = stat
	self.callable = callable

	return &self
}

func newProcessorFromSrc(detector *Detector, src []byte, stdin *os.File, stdout *os.File) (*processor, error) {
	ast, err := newAstFromGenerator(src, stdin, stdout)
	if err != nil {
		return nil, err
	}
	stat := newTextStat(src)
	return newProcessor(detector, ast, &stat), nil
}

func (self processor) Predict() float64 {
	var inputData = make(map[int]float64, 11)
	for k, v := range self.stat.GetVector() {
		inputData[k+1] = v // model start at 1
	}
	inputData[9] = self.opSerial.Predict(self.detector.opSerialModel)
	inputData[10] = self.words.Predict(self.detector.wordsModel)
	_, result := self.detector.processModel.PredictValues(inputData)
	return result[0]
}
