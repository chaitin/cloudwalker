package WebshellDetector

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

/*
WebshellDetector - Refactor version 2
Date   0929
Author Cyrus
Intro  Create a PHP-Ast generator for build Ast, from STDIO.
*/

type phpAstServer struct {
	stdin  *os.File
	stdout *os.File
}

func newPhpAstGenerator(stdin *os.File, stdout *os.File) phpAstServer {
	return phpAstServer{stdin, stdout}
}

func (self *phpAstServer) GetData(src []byte) ([]byte, error) {
	srcLen := len(src)
	if n, err := self.stdin.WriteString(strconv.Itoa(srcLen) + "\n"); err != nil || n <= 1 {
		log.Fatal("Write stdin failed")
	}
	if _, err := self.stdin.Write(src); err != nil {
		log.Fatal("Write stdin failed")
	}

	connReader := bufio.NewReaderSize(self.stdout, 4096) // Get recv data length from server
	resultLenStr, _, err := connReader.ReadLine()
	if err != nil {
		return nil, err
	}

	resultLen, err := strconv.Atoi(string(resultLenStr)) // convert []byte to int
	if err != nil {
		return nil, err
	}

	var resultTextByte = make([]byte, resultLen)
	var buf byte
	for i := 0; i < resultLen; i++ {
		buf, err = connReader.ReadByte()
		if err != nil {
			return nil, err
		}
		resultTextByte[i] = buf
	}

	return resultTextByte, nil
}
