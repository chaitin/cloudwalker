package WebshellDetector

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"testing"
// )

// var detO *Detector

// var isWebshellO = 0
// var resultO = strings.Builder{}

// func walkFunc_genOpSerialData(path string, info os.FileInfo, _ error) error {

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

// 	resultO.WriteString(fmt.Sprint(isWebshellO, " "))
// 	for i, e := range ast.GetOpSerial(&detO.hashState).data {
// 		if e != 0 {
// 			resultO.WriteString(fmt.Sprintf("%v:%v ", i+1, e))
// 		}
// 	}
// 	resultO.WriteString("\n")

// 	return nil
// }

// func Test_genOpSerialData(t *testing.T) {
// 	detO, _ = NewDefaultDetector(stdin, stdout)
// 	isWebshellO = 0
// 	filepath.Walk("../sample/frame", walkFunc_genOpSerialData)
// 	isWebshellO = 1
// 	filepath.Walk("../sample/webshell", walkFunc_genOpSerialData)
// 	isWebshellO = 0
// 	filepath.Walk("../sample/predict/frame", walkFunc_genOpSerialData)
// 	isWebshellO = 1
// 	filepath.Walk("../sample/predict/webshell", walkFunc_genOpSerialData)

// 	ioutil.WriteFile("../tools/cross-validation tool/MetaData/opSerialData.txt", []byte(resultO.String()), os.ModePerm)
// }
