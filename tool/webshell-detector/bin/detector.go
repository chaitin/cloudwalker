package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"../php"
	"../src"
)

var outputPtr *string
var htmlPtr *bool
var detector *WebshellDetector.Detector
var countRisk = 0
var countRisk1 = 0
var countRisk2 = 0
var countRisk3 = 0
var countRisk4 = 0
var countRisk5 = 0
var countFile = 0
var t0 = time.Now()
var err error

func walk(path string, info os.FileInfo, _ error) error {

	countFile++

	if strings.ToLower(filepath.Ext(path)) != ".php" &&
		strings.ToLower(filepath.Ext(path)) != ".phpt" &&
		strings.ToLower(filepath.Ext(path)) != ".php3" &&
		strings.ToLower(filepath.Ext(path)) != ".php4" &&
		strings.ToLower(filepath.Ext(path)) != ".php5" &&
		strings.ToLower(filepath.Ext(path)) != ".txt" &&
		strings.ToLower(filepath.Ext(path)) != ".bak" {
		return nil
	}
	src, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	score, err := detector.Predict(src)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	for i := 0; i < 128; i++ {
		fmt.Printf("\b")
	}
	fmt.Printf("\rTesting %-50s / %d risks / Runtime %v", info.Name(), countRisk, time.Since(t0))
	if score > 0 {
		countRisk++
		var risk int
		switch score {
		case 1, 2:
			risk = 1
			countRisk1++
		case 3:
			risk = 2
			countRisk2++
		case 4:
			risk = 3
			countRisk3++
		case 5:
			risk = 4
			countRisk4++
		case 6, 7:
			risk = 5
			countRisk5++
		default:
			risk = 0
		}
		printResult(countFile, path, risk)
	}
	return nil
}

func multiTest(detectPath string) {
	if err != nil {
		log.Fatal(err)
	}
	filepath.Walk(detectPath, walk)
}

func singleTest(path string) {
	if err != nil {
		log.Fatal(err)
	}
	src, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	score, err := detector.Predict(src)
	if err != nil {
		log.Fatal(err)
	}
	if score > 0 {
		var risk int
		switch score {
		case 1, 2:
			risk = 1
		case 3:
			risk = 2
		case 4:
			risk = 3
		case 5:
			risk = 4
		case 6, 7:
			risk = 5
		default:
			risk = 0
		}
		printResult(countFile, path, risk)
	}
}

func printResult(fileIndex int, filePath string, fileRisk int) {
	var res string
	if len(filePath) > 80 && !*htmlPtr {
		filePath = "..." + filePath[len(filePath)-77:len(filePath)]
	}
	if *htmlPtr {
		res = fmt.Sprintf("<tr><td>%08d<td>%s<td>%d\n", fileIndex, filePath, fileRisk)
	} else {
		res = fmt.Sprintf("[+] %08d %-80s Risk:%d\n", fileIndex, filePath, fileRisk)
	}
	if *outputPtr != "" {
		outputFile, err := os.OpenFile(*outputPtr, os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			log.Fatal(err)
		}
		outputFile.WriteString(res)
		outputFile.Close()
	} else {
		for i := 0; i < 128; i++ {
			fmt.Printf("\b")
		}
		fmt.Print(res)
	}
}

func work(files []string) {
	for _, v := range files {
		info, err := os.Stat(v)
		if err != nil {
			log.Fatal(err)
		}
		if info.IsDir() {
			log.Printf("Testing %s...", v)
			truePath, err := filepath.EvalSymlinks(v)
			if err != nil {
				log.Fatal(err)
			}
			multiTest(truePath)
		} else {
			singleTest(v)
		}
	}
}

func main() {
	version := "1.0.0"
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Chaitin CloudWalker Webshell Detector\n[version %s]\n\nusage: %s [options] name ...\n\n", version, os.Args[0])
		flag.PrintDefaults()
	}
	outputPtr = flag.String("output", "", "Export result to output file")
	htmlPtr = flag.Bool("html", false, "Show result as HTML")
	flag.Parse()
	if len(flag.Args()) == 0 {
		flag.Usage()
	}
	files := flag.Args()
	php.Start()
	detector, err = WebshellDetector.NewDefaultDetector(php.Stdin, php.Stdout)
	if err != nil {
		log.Println("Detector kernel cannot load.")
	}
	if *outputPtr != "" {
		outputFile, err := os.OpenFile(*outputPtr, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
		if err != nil {
			log.Fatal(err)
		}
		if *htmlPtr {
			nowStr := time.Now().Format("2006-01-02 15:04:05")
			outputFile.WriteString(`<!DOCTYPE html>
		<html>
			<head>
				<meta charset="UTF-8">
				<title>Report - Chaitin CloudWalker</title>
			</head>
			<style>
				html {
					font-family: Helvetica, Tahoma, Arial, "PingFang SC", "Hiragino Sans GB", "Heiti SC", STXihei, "Microsoft YaHei", SimHei, "WenQuanYi Micro Hei";
				}
				.header {
					vertical-align: middle;
					background: #1b1c1d;
					text-align: center;
					padding: 4px;
					border-radius: 9px 9px 9px 9px;
				}
				img {
					height: 1.2em;
				}
				h1 {
					font-weight: bolder;
				}
				h2 {
					font-weight: bold;
					color: #fff;
					font-weight: lighter;
				} 
				h3 {
					font-weight: lighter;
				} 
				table {
					width: 100%;
					border-collapse: collapse;
				} 
				th,td {
					padding: .65em;
				} 
				th {
					background: #555;
					color: #fff;
				} 
				tbody tr:nth-child(odd) {
					background: #ccc;
				} 
				th:first-child {
					border-radius: 9px 0 0 0;
				} 
				th:last-child {
					border-radius: 0 9px 0 0;
				} 
				tr:last-child td:first-child {
					border-radius: 0 0 0 9px;
				} 
				tr:last-child td:last-child {
					border-radius: 0 0 9px 0;
				}
				tbody tr:hover {
					background: #eee;		 
				}
			</style>
			<body>
				<div class="header">
					<h2><img src="data:image/png;base64,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"/>&nbsp;&nbsp;&nbsp;Chaitin CloudWalker Webshell Detector</h2>
				</div>
				<div align="center">
					<h1>Report</h1>
					<h3>Time: ` + nowStr + `</h3>
					<table style="text-align: center">
						<caption><br><captionBody></captionBody></caption>
						<thead>
							<tr>
								<th>Tested
								<th>Risk:1
								<th>Risk:2
								<th>Risk:3
								<th>Risk:4
								<th>Risk:5
						</thead>
						<tbody>
							<tr>
								<td style="background: #eee" id="t">Waiting
								<td style="background: #fec" id="r1">Waiting
								<td style="background: #edc" id="r2">Waiting
								<td style="background: #dcc" id="r3">Waiting
								<td style="background: #dbc" id="r4">Waiting
								<td style="background: #fbc" id="r5">Waiting
						</tbody>
					</table>
					<table>
						<caption><br><captionBody></captionBody></caption>
						<thead>
							<tr>
								<th>Index
								<th>File
								<th>Risk
						</thead>
						<tbody>`)
			outputFile.Close()
		}
		work(files)
		outputFile, err = os.OpenFile(*outputPtr, os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			log.Fatal(err)
		}
		if *htmlPtr {
			outputFile.WriteString(fmt.Sprintf(`</tbody></table><script>
			document.getElementById("t").innerHTML="%d";
			document.getElementById("r1").innerHTML="%d";
			document.getElementById("r2").innerHTML="%d";
			document.getElementById("r3").innerHTML="%d";
			document.getElementById("r4").innerHTML="%d";
			document.getElementById("r5").innerHTML="%d";
		</script></div></body></html>`, countFile, countRisk1, countRisk2, countRisk3, countRisk4, countRisk5))
			outputFile.Close()
		}
	} else {
		fmt.Println("  _____ _                 ___          __   _ _")
		fmt.Println(" / ____| |               | \\ \\        / /  | | |")
		fmt.Println("| |    | | ___  _   _  __| |\\ \\  /\\  / /_ _| | | _____ _ __   __   ___ ")
		fmt.Println("| |    | |/ _ \\| | | |/ _` | \\ \\/  \\/ / _` | | |/ / _ \\ '__| /_ | / _ \\ ")
		fmt.Println("| |____| | (_) | |_| | (_| |  \\  /\\  / (_| | |   <  __/ |     | | ||_||")
		fmt.Println(" \\_____|_|\\___/ \\__,_|\\__,_|   \\/  \\/ \\__,_|_|_|\\_\\___|_|     |_(_)___/ \n")
		log.Println("Detector started.")

		work(files)

		for i := 0; i < 256; i++ {
			fmt.Printf("\b")
		}
		fmt.Print("\n\n")
		log.Printf("Risk (level1): %d files", countRisk1)
		log.Printf("Risk (level2): %d files", countRisk2)
		log.Printf("Risk (level3): %d files", countRisk3)
		log.Printf("Risk (level4): %d files", countRisk4)
		log.Printf("Risk (level5): %d files", countRisk5)
		log.Printf("Tested: %d files", countFile)

		fmt.Print("\n\n")
		log.Printf("Detector done (%v).\n", time.Since(t0))
	}
}
