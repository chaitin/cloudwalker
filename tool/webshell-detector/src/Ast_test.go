package WebshellDetector

import (
	"testing"
)

func Test_Ast_New(t *testing.T) {
	ast, err := newAstFromGenerator([]byte(`<?php echo f(g()+231*544+555+f(1,2,3+4,5)); ?>`), stdin, stdout)
	if err != nil {
		t.Error(err.Error())
		t.Log(ast)
	}
}

func Test_Ast_GetOpSerial(t *testing.T) {
	ast, err := newAstFromGenerator([]byte(`<?php echo f(g()+231*544+555+f(1,2,3+4,5)); ?>`), stdin, stdout)
	if err != nil {
		t.Error(err.Error())
	}

	det, err := NewDefaultDetector(stdin, stdout)
	if err != nil {
		t.Error(err.Error())
	}

	vec1 := ast.GetOpSerial(&det.hashState).data

	vec2 := make([]float64, len(vec1))
	vec2[0] = 3
	vec2[2] = 1
	vec2[6] = 1
	vec2[20] = 2
	vec2[26] = 1
	vec2[45] = 1
	vec2[195] = 1

	for i := 0; i < len(vec1); i++ {
		if vec1[i] != vec2[i] {
			t.Error("vec1 != vec2")
		}
	}
}

func Test_Ast_GetWords(t *testing.T) {

	ast, err := newAstFromGenerator([]byte(`<?php echo f(g()+231*544+555+f(1,2,3+4,5)); ?>`), stdin, stdout)
	if err != nil {
		return
	}
	res, _ := ast.GetWordsAndCallable()
	vec := res.data

	t.Log(len(vec))

	for i := 0; i < len(vec); i++ {
		if vec[i] != "f" && vec[i] != "g" {
			t.Error("vec1 != vec2")
		}
	}

}

func Test_Ast_IsCallable(t *testing.T) {
	ast1, err := newAstFromGenerator([]byte(`<?php $a="ass"."ert"; $a("phpinfo()"); ?>`), stdin, stdout)
	if err != nil {
		return
	}
	_, res1 := ast1.GetWordsAndCallable()
	if res1 != true {
		t.Error("Missing Callable")
	}

	ast2, err := newAstFromGenerator([]byte(`<?php $a="ass"."ert"; echo $a; ?>`), stdin, stdout)
	if err != nil {
		return
	}
	_, res2 := ast2.GetWordsAndCallable()
	if res2 != false {
		t.Error("Misjudge Callable")
	}
}
