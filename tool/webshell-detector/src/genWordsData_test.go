package WebshellDetector

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"path/filepath"
// 	"testing"
// )

// var isWebshellW = 0
// var resultW struct {
// 	Data   [][]string
// 	Expect []int
// }

// func walkFunc_genWordsData(path string, info os.FileInfo, _ error) error {

// 	if filepath.Ext(path) != ".php" {
// 		return nil
// 	}

// 	fmt.Printf("%-50v [%v]\n", info.Name(), path)
// 	src, _ := ioutil.ReadFile(path)
// 	ast, err := newAstFromServer(src)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return nil
// 	}

// 	words := ast.GetWords()
// 	resultW.Data = append(resultW.Data, words.data)
// 	resultW.Expect = append(resultW.Expect, isWebshellW)

// 	return nil
// }

// func Test_genWordsData(t *testing.T) {
// 	isWebshellW = 0
// 	filepath.Walk("../sample/frame", walkFunc_genWordsData)
// 	isWebshellW = 1
// 	filepath.Walk("../sample/webshell", walkFunc_genWordsData)
// 	isWebshellW = 0
// 	filepath.Walk("../sample/predict/frame", walkFunc_genWordsData)
// 	isWebshellW = 1
// 	filepath.Walk("../sample/predict/webshell", walkFunc_genWordsData)

// 	data, err := json.Marshal(resultW)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}

// 	ioutil.WriteFile("../tools/cross-validation tool/MetaData/wordsData.txt", data, os.ModePerm)
// }
