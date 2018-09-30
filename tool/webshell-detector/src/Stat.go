package WebshellDetector

import (
	"encoding/json"
	"math"
	"regexp"
	"strings"

	"github.com/grd/stat"
)

/*
WebshellDetector - Refactor version 1
Date	0814
Author	Cyrus
Intro	Statastic model for PHP file
*/

type textStat struct {
	LM  float64 // int64 in fact
	LVC float64
	WM  float64 // int64 in fact
	WVC float64
	SR  float64
	TR  float64
	SPL float64
	IE  float64
}

type arrayStatState struct {
	MinStat textStat
	MaxStat textStat
}

var nan = math.NaN()
var nanStat = textStat{nan, nan, nan, nan, nan, nan, nan, nan}

func newArrayStatState() arrayStatState {

	return arrayStatState{nanStat, nanStat}
}

func (self *arrayStatState) Load(bytes []byte) error {
	return json.Unmarshal(bytes, self)
}

func outOfRange(x float64, min float64, max float64) bool {
	if min != math.NaN() {
		if x < min {
			return true
		}
	}
	if max != math.NaN() {
		if x > max {
			return true
		}
	}
	return false
}

func statLines(src string) []int64 {
	// get all lines
	var result []int64
	splitResult := strings.Split(src, "\n")
	for _, v := range splitResult {
		result = append(result, int64(len(v)))
	}
	return result
}

func lineMax(src string) int64 {
	lines := stat.IntSlice(statLines(src))
	if len(lines) > 0 {
		result, _ := stat.Max(lines)
		return int64(result)
	} else {
		return 0
	}

}

func lineVariationCoefficient(src string) float64 {
	lines := stat.IntSlice(statLines(src))
	return math.Sqrt(stat.Variance(lines)) / stat.Mean(lines)
}

func statWords(src string) []int64 {
	// get all words
	wordReg, _ := regexp.Compile(`[a-zA-Z0-9]*`)
	var result []int64
	regRusult := wordReg.FindAllStringIndex(src, -1)
	for _, index := range regRusult {
		if index[1]-index[0] > 0 {
			result = append(result, int64(index[1]-index[0]))

		}
	}
	return result
}

func wordMax(src string) int64 {
	words := stat.IntSlice(statWords(src))
	if len(words) > 0 {
		result, _ := stat.Max(words)
		return int64(result)
	} else {
		return 0
	}
}

func wordVariationCoefficient(src string) float64 {
	words := stat.IntSlice(statWords(src))
	return math.Sqrt(stat.Variance(words)) / stat.Mean(words) * 100
}

func symbolRatio(src string) float64 {
	symbolReg, _ := regexp.Compile(`[^a-zA-Z0-9]`)
	symbolNumber := len(symbolReg.FindAllString(src, -1))
	if len(src) > 0 {
		return float64(symbolNumber) / float64(len(src)) * 100
	} else {
		return 1
	}
}

func tagRatio(src string) float64 {
	tagReg, _ := regexp.Compile(`<[\x00-\xFF]*?>`)
	tagNumber := len(tagReg.FindAllString(src, -1))
	lenWords := len(statWords(src))
	if lenWords > 0 {
		return float64(tagNumber) / float64(lenWords) * 100
	} else {
		return 0
	}
}

func statementPerLine(src string) float64 {
	statementReg, _ := regexp.Compile(`;`)
	statementNumber := len(statementReg.FindAllString(src, -1))
	lenLines := len(statLines(src))
	if lenLines > 0 {
		return float64(statementNumber) / float64(lenLines)
	} else {
		return 0
	}
}

func infomationEntropy(src string) float64 {
	var lst []float64
	chrs := 0.00
	for i := 0; i < 256; i++ {
		lst = append(lst, 0)
	}
	for _, chr := range src {
		if 0 <= chr && chr < 256 && chr != '\n' {
			lst[chr]++
			chrs++
		}
	}
	var result float64
	for i := 0; i < 256; i++ {
		if lst[i] > 0 {
			result -= lst[i] / chrs * math.Log2(lst[i]/chrs)
		}
	}
	return result
}

func newTextStat(src []byte) textStat {
	var self textStat
	str_src := string(src)
	self.LM = float64(lineMax(str_src))
	self.LVC = lineVariationCoefficient(str_src)
	self.WM = float64(wordMax(str_src))
	self.WVC = wordVariationCoefficient(str_src)
	self.SR = symbolRatio(str_src)
	self.TR = tagRatio(str_src)
	self.SPL = statementPerLine(str_src)
	self.IE = infomationEntropy(str_src)
	return self
}

func (self textStat) GetVector() []float64 {
	return []float64{self.LM, self.LVC, self.WM, self.WVC, self.SR, self.TR, self.SPL, self.IE}
}

func (self textStat) IsAbnormal(std arrayStatState) bool {
	return outOfRange(self.LM, std.MinStat.LM, std.MaxStat.LM) ||
		outOfRange(self.LVC, std.MinStat.LVC, std.MaxStat.LVC) ||
		outOfRange(self.WM, std.MinStat.WM, std.MaxStat.WM) ||
		outOfRange(self.WVC, std.MinStat.WVC, std.MaxStat.WVC) ||
		outOfRange(self.SR, std.MinStat.SR, std.MaxStat.SR) ||
		outOfRange(self.TR, std.MinStat.TR, std.MaxStat.TR) ||
		outOfRange(self.SPL, std.MinStat.SPL, std.MaxStat.SPL) ||
		outOfRange(self.IE, std.MinStat.IE, std.MaxStat.IE)
}
