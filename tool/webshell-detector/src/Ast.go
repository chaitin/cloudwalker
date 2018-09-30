package WebshellDetector

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"strings"
)

/*
WebshellDetector - Refactor version 1
Date	0814
Author	Twice
Intro	operations for PHP AST
*/

type ast struct {
	root interface{}
}

func newAst(data []byte) (*ast, error) {
	var astData interface{}
	if err := json.Unmarshal(data, &astData); err != nil {
		return nil, err
	}

	if astMap, ok := astData.(map[string]interface{}); ok {
		if astNode, ok := astMap["ast"]; ok {
			return &ast{transformAstNode(astNode)}, nil
		} else if reason, ok := astMap["reason"]; ok {
			if reasonStr, ok := reason.(string); ok {
				return nil, errors.New(reasonStr)
			}
		}
	}

	return &ast{nil}, nil
}

func transformAstNode(ast interface{}) (result interface{}) {
	switch value := ast.(type) {
	case float64:
		result = value
	case string:
		result = value
	case nil:
		result = value
	case []interface{}:
		resArray := make([]interface{}, 0)
		for _, v := range value {
			resArray = append(resArray, transformAstNode(v))
		}
		result = resArray
	case map[string]interface{}:
		if _, ok := value["kind"]; ok {
			result = astNode{
				int(value["kind"].(float64)),
				int(value["flags"].(float64)),
				int(value["lineno"].(float64)),
				transformAstNode(value["children"])}
		} else {
			resMap := make(map[string]interface{})
			for k, v := range value {
				resMap[k] = transformAstNode(v)
			}
			result = resMap
		}
	}

	return
}

func newAstFromGenerator(src []byte, stdin *os.File, stdout *os.File) (*ast, error) {
	server := newPhpAstGenerator(stdin, stdout)
	data, err := server.GetData(src)
	if err != nil {
		return nil, err
	}
	return newAst(data)
}

type opQueueNode struct {
	Key    string
	Value  interface{}
	Layer  int
	Father *astNode
}

func getOpSerial(ast interface{}) (result [][]int) {

	var queue []opQueueNode

	nowSerial := make([]int, 0)

	queue = append(queue, opQueueNode{"root", ast, 0, nil})

	for len(queue) != 0 {
		node := queue[0]
		switch value := node.Value.(type) {
		case astNode:
			nowSerial = append(nowSerial, value.Kind)
			//fmt.Printf("kind:%4v(%4v)\t", value.Kind, node.Layer)
			queue = append(queue, opQueueNode{"children", value.Children, node.Layer, &value})
		case string:
		case float64:
		case nil:
			if node.Key == "separator" {
				if len(nowSerial) != 0 {
					//fmt.Printf("\n")
					finishSerial := make([]int, 1)
					finishSerial[0] = node.Father.Kind
					finishSerial = append(finishSerial, nowSerial...)

					result = append(result, finishSerial)
					nowSerial = make([]int, 0)
				}
			}
		case []interface{}:
			for i, v := range value {
				queue = append(queue, opQueueNode{string(i), v, node.Layer + 1, node.Father})
			}
			if node.Key == "children" {
				queue = append(queue, opQueueNode{"separator", nil, node.Layer + 1, node.Father})
			}
		case map[string]interface{}:
			keys := make([]string, 0)
			for key, _ := range value {
				keys = append(keys, key)
			}
			sort.Strings(keys)

			for _, k := range keys {
				queue = append(queue, opQueueNode{k, value[k], node.Layer + 1, node.Father})
			}
			if node.Key == "children" {
				queue = append(queue, opQueueNode{"separator", nil, node.Layer + 1, node.Father})
			}
		}
		queue = queue[1:]
	}

	return
}

func cleanOpSerial(data [][]int, maxLen int) [][]int {

	blockCompare := func(data []int, i int, length int) bool {
		if i+2*length > len(data) {
			return false
		}
		for k := 0; k < length; k++ {
			if data[i+k] != data[i+length+k] {
				return false
			}
		}
		return true
	}

	var result [][]int

	for _, v := range data {
		tmp := v
		for i := maxLen; i >= 1; i-- {
			for j := 0; j < len(tmp); {
				if blockCompare(tmp, j, i) {
					tmp = append(tmp[:j], tmp[j+i:]...)
				} else {
					j++
				}
			}
		}
		result = append(result, tmp)
	}

	return result
}

func cleanOpSerialRepeatedly(serials *[][]int, cleanTimes, cleanLength int) {
	for i := 1; i <= cleanTimes; i++ {
		*serials = cleanOpSerial(*serials, cleanLength)
	}
}

type arrayHashState struct {
	Table   map[string]int
	DeTable map[int]string
}

func newArrayHashState() arrayHashState {
	return arrayHashState{make(map[string]int), make(map[int]string)}
}

func serialToString(serial []int) string {
	builder := strings.Builder{}
	for _, value := range serial {
		builder.WriteByte(byte(value / 100))
		builder.WriteByte(byte(value % 100))
	}

	return builder.String()
}

func stringToSerial(str string) (result []int) {
	for i := 0; i < len(str); i += 2 {
		result = append(result, int([]byte(str)[i])*100+int([]byte(str)[i+1]))
	}

	return
}

func (state *arrayHashState) Hash(span []int) int {
	if value, ok := state.Table[serialToString(span)]; ok {
		return value
	} else {
		return -1
	}
}

func (state *arrayHashState) Find(hash int) []int {
	return stringToSerial(state.DeTable[hash])
}

func (state *arrayHashState) Load(bytes []byte) error {
	return json.Unmarshal(bytes, state)
}

func (ast ast) GetOpSerial(state *arrayHashState) opSerial {
	serials := getOpSerial(ast.root)
	cleanOpSerialRepeatedly(&serials, 10, 5)

	vector := make([]float64, len(state.Table))
	for _, s := range serials {
		hash := state.Hash(s)
		if hash != -1 {
			vector[hash]++
		}
	}

	return opSerial{vector}
}

type iterateAstState struct {
	Key             string
	KindPredict     func(int) bool
	FlagsPredict    func(int) bool
	ChildrenPredict func(interface{}) bool
}

func predictTrue(int) bool {
	return true
}

func predictTrueInterface(interface{}) bool {
	return true
}

func (state iterateAstState) Iterate(ast interface{}) (count int) {
	count = 0

	switch value := ast.(type) {
	case astNode:
		if state.KindPredict(value.Kind) && state.FlagsPredict(value.Flag) && state.ChildrenPredict(value.Children) {
			count += 1
		}
		count += state.Iterate(value.Children)
	case float64:
	case string:
	case []interface{}:
		for i, v := range value {
			state.Key = string(i)
			count += state.Iterate(v)
		}
	case map[string]interface{}:
		for k, v := range value {
			state.Key = k
			count += state.Iterate(v)
		}
	}

	return
}

func checkNameAndKind(ast interface{}, nameChecker func(string) bool, kindChecker func(int) bool) int {
	state := iterateAstState{
		"",
		kindChecker,
		predictTrue,
		func(ast interface{}) bool {
			if nameMap, ok := ast.(map[string]interface{}); ok {
				if name, ok := nameMap["name"]; ok {
					if nameStr, ok := name.(string); ok {
						return nameChecker(nameStr)
					}
				}
			}
			return false
		}}

	return state.Iterate(ast)
}

func (ast ast) GetWordsAndCallable() (words, bool) {
	result := false
	vector := make([]string, 0)
	checkNameAndKind(ast.root, func(s string) bool {
		vector = append(vector, s)
		return true
	}, func(k int) bool {
		if k == 269 || k == 265 || k == 515 || k == 768 || k == 769 {
			result = true
		}
		return true
	})

	return words{vector}, result
}
