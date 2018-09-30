package WebshellDetector

/*
WebshellDetector - Refactor version 1
Date	0814
Author	Twice
Intro	A node structure in PHP AST
*/

type astNode struct {
	Kind int
	Flag int
	LineNo int
	Children interface{}
}
