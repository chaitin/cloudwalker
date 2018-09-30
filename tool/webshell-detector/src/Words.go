package WebshellDetector

import "github.com/CyrusF/go-bayesian"

/*
WebshellDetector - Refactor version 1
Date	0814
Author	Twice
Intro	predict webshell by words using Naive Bayes
*/

type words struct {
	data []string
}

func (words words) Predict(model *bayesian.Classifier) float64 {
	score, _, _ := model.Classify(words.data...)
	return score["webshell"] / (score["webshell"] + score["normal"])
}
