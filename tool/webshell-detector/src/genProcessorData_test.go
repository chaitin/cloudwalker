package WebshellDetector

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"testing"
// )

// var detP *Detector

// var isWebshellP = 0
// var resultP = strings.Builder{}

// func walkFunc_genProcessorData(path string, info os.FileInfo, _ error) error {

// 	if filepath.Ext(path) != ".php" {
// 		return nil
// 	}

// 	fmt.Printf("%-50v [%v]\n", info.Name(), path)
// 	src, _ := ioutil.ReadFile(path)
// 	proc, err := newProcessorFromSrc(detP, src)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return nil
// 	}

// 	resultP.WriteString(fmt.Sprint(isWebshellP, " "))
// 	for i, e := range proc.stat.GetVector() {
// 		resultP.WriteString(fmt.Sprintf("%v:%v ", i+1, e))
// 	}
// 	resultP.WriteString(fmt.Sprintf("9:%v ", proc.opSerial.Predict(proc.detector.opSerialModel)))
// 	resultP.WriteString(fmt.Sprintf("10:%v ", proc.words.Predict(proc.detector.wordsModel)))
// 	resultP.WriteString("\n")

// 	return nil
// }

// func Test_genProcessorData(t *testing.T) {
// 	detP, _ = NewDefaultDetector(stdin, stdout)
// 	isWebshellP = 0
// 	filepath.Walk("../sample/frame", walkFunc_genProcessorData)
// 	isWebshellP = 1
// 	filepath.Walk("../sample/webshell", walkFunc_genProcessorData)
// 	isWebshellP = 0
// 	filepath.Walk("../sample/predict/frame", walkFunc_genProcessorData)
// 	isWebshellP = 1
// 	filepath.Walk("../sample/predict/webshell", walkFunc_genProcessorData)

// 	ioutil.WriteFile("../tools/cross-validation tool/MetaData/processorData.txt", []byte(resultP.String()), os.ModePerm)
// }
