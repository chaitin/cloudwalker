package WebshellDetector

import (
	"testing"

	"../php"
)

func Test_Processor_NewProcessor(t *testing.T) {
	detector, _ := NewDefaultDetector(stdin, stdout)
	src := []byte(`<?php echo f(g()+231*544+555+f(1,2,3+4,5)); ?>`)
	ast, err := newAstFromGenerator(src, stdin, stdout)
	if err != nil {
		t.Error("FAILED - Test_Processor_NewProcessor - AST")
	}
	stat := newTextStat(src)
	processor := newProcessor(detector, ast, &stat)
	if &processor == nil {
		t.Error("FAILED - Test_Processor_NewProcessor")
	}
}

func Test_Processor_NewProcessorFromSrc(t *testing.T) {
	detector, _ := NewDefaultDetector(stdin, stdout)
	src := []byte(`<?php echo f(g()+231*544+555+f(1,2,3+4,5)); ?>`)
	var err error
	stdin := php.Stdin
	stdout := php.Stdout
	processor, err := newProcessorFromSrc(detector, src, stdin, stdout)
	if err != nil || &processor == nil {
		t.Error("FAILED - Test_Processor_NewProcessorFromSrc")
	}
}
func Test_Processor_Predict(t *testing.T) {
	detector, _ := NewDefaultDetector(stdin, stdout)
	src := []byte(`<?php echo f(g()+231*544+555+f(1,2,3+4,5)); ?>`)
	stdin := php.Stdin
	stdout := php.Stdout
	processor, _ := newProcessorFromSrc(detector, src, stdin, stdout)
	result := processor.Predict()
	if &result == nil {
		t.Error("FAILED - Test_Processor_NewProcessor - Predict")
	}
}
