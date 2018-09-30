package WebshellDetector

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"
)

/*
WebshellDetector - Refactor version 1
Date	0814
Author	Cyrus
Intro	Some regular expression match for PHP file
*/

type regItem struct {
	ordered bool
	regItem []*regexp.Regexp
	strItem [][]byte
}

type regMatcher struct {
	regData []regItem
}

var presets = []regItem{
	regItem{false, []*regexp.Regexp{regexp.MustCompile(`function\_exists\s*\(\s*[\'|\"](popen|exec|proc\_open|system|passthru)+[\'|\"]\s*\)`)}, [][]byte{}},
	regItem{false, []*regexp.Regexp{regexp.MustCompile(`(exec|shell\_exec|system|passthru)+\s*\(\s*\$\_(\w+)\[(.*)\]\s*\)`)}, [][]byte{}},
	regItem{false, []*regexp.Regexp{regexp.MustCompile(`((udp|tcp)\:\/\/(.*)\;)+`)}, [][]byte{}},
	regItem{false, []*regexp.Regexp{regexp.MustCompile(`preg\_replace\s*\((.*)\/e(.*)\,\s*\$\_(.*)\,(.*)\)`)}, [][]byte{}},
	regItem{false, []*regexp.Regexp{regexp.MustCompile(`preg\_replace\s*\((.*)\(base64\_decode\(\$`)}, [][]byte{}},
	regItem{false, []*regexp.Regexp{regexp.MustCompile(`(eval|assert|include|require|include\_once|require\_once)+\s*\(\s*(base64\_decode|str\_rot13|gz(\w+)|file\_(\w+)\_contents|(.*)php\:\/\/input)+`)}, [][]byte{}},
	regItem{false, []*regexp.Regexp{regexp.MustCompile(`(eval|assert|include|require|include\_once|require\_once|array\_map|array\_walk)+\s*\(\s*\$\_(GET|POST|REQUEST|COOKIE|SERVER|SESSION)+\[(.*)\]\s*\)`)}, [][]byte{}},
	regItem{false, []*regexp.Regexp{regexp.MustCompile(`eval\s*\(\s*\(\s*\$\$(\w+)`)}, [][]byte{}},
	regItem{false, []*regexp.Regexp{regexp.MustCompile(`(include|require|include\_once|require\_once)+\s*\(\s*[\'|\"](\w+)\.(jpg|gif|ico|bmp|png|txt|zip|rar|htm|css|js)+[\'|\"]\s*\)`)}, [][]byte{}},
	regItem{false, []*regexp.Regexp{regexp.MustCompile(`\$\_(\w+)(.*)(eval|assert|include|require|include\_once|require\_once)+\s*\(\s*\$(\w+)\s*\)`)}, [][]byte{}},
	regItem{false, []*regexp.Regexp{regexp.MustCompile(`\(\s*\$\_FILES\[(.*)\]\[(.*)\]\s*\,\s*\$\_(GET|POST|REQUEST|FILES)+\[(.*)\]\[(.*)\]\s*\)`)}, [][]byte{}},
	// regItem{false, []*regexp.Regexp{regexp.MustCompile(`(fopen|fwrite|fputs|file\_put\_contents)+\s*\((.*)\$\_(GET|POST|REQUEST|COOKIE|SERVER)+\[(.*)\](.*)\)`)}, [][]byte{}},
	regItem{false, []*regexp.Regexp{regexp.MustCompile(`echo\s*curl\_exec\s*\(\s*\$(\w+)\s*\)`)}, [][]byte{}},
	regItem{false, []*regexp.Regexp{regexp.MustCompile(`new com\s*\(\s*[\'|\"]shell(.*)[\'|\"]\s*\)`)}, [][]byte{}},
	regItem{false, []*regexp.Regexp{regexp.MustCompile(`\$(.*)\s*\((.*)\/e(.*)\,\s*\$\_(.*)\,(.*)\)`)}, [][]byte{}},
	regItem{false, []*regexp.Regexp{regexp.MustCompile(`\$\_\=(.*)\$\_`)}, [][]byte{}},
	regItem{false, []*regexp.Regexp{regexp.MustCompile(`\$\_(GET|POST|REQUEST|COOKIE|SERVER)+\[(.*)\]\(\s*\$(.*)\)`)}, [][]byte{}},
	regItem{false, []*regexp.Regexp{regexp.MustCompile(`\$(\w+)\s*\(\s*\$\_(GET|POST|REQUEST|COOKIE|SERVER)+\[(.*)\]\s*\)`)}, [][]byte{}},
	regItem{false, []*regexp.Regexp{regexp.MustCompile(`\$(\w+)\s*\(\s*\$\{(.*)\}`)}, [][]byte{}},
	regItem{false, []*regexp.Regexp{regexp.MustCompile(`\$(\w+)\s*\(\s*chr\(\d+\)`)}, [][]byte{}},
}

func defaultRegMatcher() *regMatcher {
	var self regMatcher
	self.regData = presets
	return &self
}

func newRegMatcher(file io.Reader) (*regMatcher, error) {
	reader := bufio.NewReader(file)
	var self regMatcher
	self.regData = presets
	errNum := 0
	var regGroup []*regexp.Regexp
	var strGroup [][]byte
	inGroup := false
	oederedGroup := false
	lineCount := 0
	for line, _, err := reader.ReadLine(); err == nil; line, _, err = reader.ReadLine() {
		lineCount++
		str := strings.Trim(string(line), " ")
		if str == "" || (strings.Index(str, "##") == 0) {
			continue
		}
		if strings.Index(strings.ToLower(str), "#start#") == 0 {
			if inGroup {
				errStr := fmt.Sprintf("RegMatcher Error: except #End but found #Start in %s (Line %d)", "bindata", lineCount)
				return nil, errors.New(errStr)
			}
			inGroup = true
			oederedGroup = false
			regGroup = make([]*regexp.Regexp, 0)
			strGroup = make([][]byte, 0)
			continue
		}
		if strings.Index(strings.ToLower(str), "#end#") == 0 {
			if !inGroup {
				errStr := fmt.Sprintf("RegMatcher Error: except #Start but found #End in %s (Line %d)", "bindata", lineCount)
				return nil, errors.New(errStr)
			}
			inGroup = false
			self.regData = append(self.regData, regItem{oederedGroup, regGroup, strGroup})
			continue
		}
		if strings.Index(strings.ToLower(str), "#ordered#") == 0 {
			if !inGroup {
				errStr := fmt.Sprintf("RegMatcher Error: [Ordered] label need regexp group, found #ordered# in %s (Line %d)", "bindata", lineCount)
				return nil, errors.New(errStr)
			}
			oederedGroup = true
			continue
		}
		if inGroup {
			if strings.Index(strings.ToLower(str), "#pre#") == 0 {
				preStr := []byte(strings.TrimSpace(str[5:len(str)]))
				strGroup = append(strGroup, preStr)
			} else {
				if reg, err := regexp.Compile(str); err == nil {
					regGroup = append(regGroup, reg)
				} else {
					errNum++
					fmt.Printf("%v\n", str)
				}
			}
		} else {
			if strings.Index(strings.ToLower(str), "#pre#") == 0 {
				preStr := []byte(strings.TrimSpace(str[5:len(str)]))
				self.regData = append(self.regData, regItem{false, []*regexp.Regexp{}, [][]byte{preStr}})
			} else {
				if reg, err := regexp.Compile(str); err == nil {
					self.regData = append(self.regData, regItem{false, []*regexp.Regexp{reg}, [][]byte{}})
				} else {
					errNum++
					fmt.Printf("%v\n", str)
				}
			}
		}
	}
	if inGroup {
		errStr := fmt.Sprintf("RegMatcher Error: could not find #End in %s", "bindata")
		return nil, errors.New(errStr)
	}
	if errNum > 0 {
		fmt.Printf("RegExp compiler: %d failed, ignored\n", errNum)
	}
	return &self, nil
}

func pretreat(src []byte) []byte {

	getFirstNChars := func(src []byte, k int, n int) string {
		if n < 0 {
			if k+n > 0 {
				return fmt.Sprintf("%s", src[k+n:k])
			} else {
				return fmt.Sprintf("%s", src[0:k])
			}
		} else {
			if k+n <= len(src) {
				return fmt.Sprintf("%s", src[k:k+n])
			} else {
				return fmt.Sprintf("%s", src[k:len(src)])
			}
		}
	}

	buf := bytes.NewBufferString("")
	sgnPhp := false
	sgnValue := false
	sgnStr := false
	sgnList := false
	sgnArray := false
	sgnCommentLine := false
	sgnCommentBlock := false
	tmpValue := bytes.NewBufferString("")
	listLevel := 0
	arrayLevel := 0

	for k, v := range src {

		if getFirstNChars(src, k, 2) == "<?" && !sgnPhp {
			sgnPhp = true
		}
		if !sgnPhp {
			continue
		}
		if (getFirstNChars(src, k, 2) == "//" || v == '#') && !sgnStr {
			sgnCommentLine = true
		}
		if (getFirstNChars(src, k, 2) == "?>" || v == '\n') && sgnCommentLine {
			sgnCommentLine = false
		}
		if getFirstNChars(src, k, 2) == "/*" && !sgnStr {
			sgnCommentBlock = true
		}
		if (getFirstNChars(src, k, 2) == "?>" || getFirstNChars(src, k, 2) == "*/") && !sgnStr {
			sgnCommentBlock = false
		}
		if sgnCommentLine || sgnCommentBlock {
			continue
		}
		if v == '$' && !sgnStr && !sgnList {
			sgnValue = true
			tmpValue.Reset()
			buf.WriteByte(v)
			continue
		}
		if sgnValue {
			if (src[k-1] == '$' &&
				(unicode.IsUpper(rune(v)) ||
					unicode.IsLower(rune(v)) ||
					v == '_')) ||
				(src[k-1] != '$' &&
					(unicode.IsUpper(rune(v)) ||
						unicode.IsLower(rune(v)) ||
						v == '_' ||
						unicode.IsDigit(rune(v)))) {
				tmpValue.WriteByte(v)
			} else {
				if tmpValue.Len() > 0 {
					tmpValueByte, _ := tmpValue.ReadBytes(byte(0))
					tmpValueStr := fmt.Sprintf("%s", tmpValueByte)
					if tmpValueStr != "_SERVER" &&
						tmpValueStr != "_GET" &&
						tmpValueStr != "_POST" &&
						tmpValueStr != "_COOKIE" &&
						tmpValueStr != "_FILES" &&
						tmpValueStr != "_ENV" &&
						tmpValueStr != "_REQUEST" &&
						tmpValueStr != "_SESSION" &&
						tmpValueStr != "GLOBALS" {
						buf.WriteByte(byte('v'))
					} else {
						buf.Write(tmpValueByte)
					}
				}
				sgnValue = false
			}
		}

		if v == '"' && !sgnStr {
			sgnStr = true
			if !sgnList {
				buf.WriteByte(v)
			}
			continue
		}
		if v == '"' && src[k-1] != '\\' && sgnStr {
			sgnStr = false
			if !sgnList {
				buf.WriteByte(v)
			}
			continue
		}
		if sgnStr && !sgnList {
			continue
		}
		if v == '[' && !sgnList && !sgnStr {
			sgnList = true
			listLevel = 1
			buf.WriteByte(v)
			continue
		}
		if v == '[' && sgnList && !sgnStr {
			listLevel++
			continue
		}
		if v == ']' && sgnList && !sgnStr {
			listLevel--
			if listLevel == 0 {
				sgnList = false
				buf.WriteByte(v)
				continue
			}
		}
		if v == '(' && !sgnStr && !sgnList {
			sgnArray = true
			arrayLevel++
			buf.WriteByte(v)
			continue
		}
		if v == ')' && !sgnStr && !sgnList {
			arrayLevel--
			if arrayLevel == 0 {
				sgnArray = false
				buf.WriteByte(v)
				continue
			} else {
				buf.WriteByte(v)
				continue
			}
		}
		if sgnArray {
			if v == ',' {
				buf.WriteByte(v)
			}
			continue
		}
		if getFirstNChars(src, k, -2) == "?>" {
			sgnPhp = false
			sgnValue = false
			sgnStr = false
			sgnList = false
			sgnArray = false
			sgnCommentLine = false
			sgnCommentBlock = false
			if v == '\n' {
				buf.WriteByte(v)
			}
			continue
		}
		if !sgnValue && !sgnList {
			if v != ' ' || getFirstNChars(src, k, -5) == "<?php" || getFirstNChars(src, k, -2) == "<?" {
				buf.WriteByte(v)
			}
		}
	}
	res, err := buf.ReadBytes(byte(0))
	if err != io.EOF {
		// fmt.Printf("%v", err)
		return src
	}
	return res
}

func (self regMatcher) IsMatched(code []byte) int {
	codePretreated := pretreat(code)
	res := 0
	for _, regGroup := range self.regData {
		regGroupFlag := true
		regGroupPosi := -1
		for _, reg := range regGroup.regItem {
			loc := reg.FindStringIndex(string(code))
			if loc == nil {
				regGroupFlag = false
				break
			} else {
				if loc[0] < regGroupPosi && regGroup.ordered {
					regGroupFlag = false
					break
				} else {
					regGroupPosi = loc[0]
				}
			}
		}
		regGroupPosi = -1
		for _, reg := range regGroup.strItem {
			loc := bytes.Index(codePretreated, reg)
			if loc == -1 {
				regGroupFlag = false
				break
			} else {
				if loc < regGroupPosi && regGroup.ordered {
					regGroupFlag = false
					break
				} else {
					regGroupPosi = loc
				}
			}
		}
		if regGroupFlag {
			// fmt.Printf("%v\n", k)
			// fmt.Printf("Ordered = %v\n", regGroup.ordered)
			// for _, i := range regGroup.regItem {
			// 	fmt.Printf("%v\n", i)
			// }
			// for _, i := range regGroup.strItem {
			// 	fmt.Printf("#pre# %s\n", i)
			// }
			// fmt.Println()
			res++
			return res // for speed uuuuup, or run all test
		}
	}
	return res
}
