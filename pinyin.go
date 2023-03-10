package pinyin

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	tones = [][]rune{
		{'ā', 'ē', 'ī', 'ō', 'ū', 'ǖ', 'Ā', 'Ē', 'Ī', 'Ō', 'Ū', 'Ǖ'},
		{'á', 'é', 'í', 'ó', 'ú', 'ǘ', 'Á', 'É', 'Í', 'Ó', 'Ú', 'Ǘ'},
		{'ǎ', 'ě', 'ǐ', 'ǒ', 'ǔ', 'ǚ', 'Ǎ', 'Ě', 'Ǐ', 'Ǒ', 'Ǔ', 'Ǚ'},
		{'à', 'è', 'ì', 'ò', 'ù', 'ǜ', 'À', 'È', 'Ì', 'Ò', 'Ù', 'Ǜ'},
	}
	noTones = []rune{'a', 'e', 'i', 'o', 'u', 'v', 'A', 'E', 'I', 'O', 'U', 'V'}
)

var (
	// 声调字母映射，例：ā-a
	toneMap map[rune]rune
	// 拼音字典
	pinyinDic map[rune]string
	hasInit   bool
)

type Fmt uint

const (
	WithTone   Fmt = iota // 带声调 例：quán
	NoTone                // 不带声调，例：quan
	FirstUpper            // 首字母大写，例：Quan
)

type Pinyin struct {
	chinese string
	fmt     Fmt
	split   string
}

func init() {
	toneMap = make(map[rune]rune)
	pinyinDic = make(map[rune]string)
	for _, runes := range tones {
		for i, tone := range runes {
			toneMap[tone] = noTones[i]
		}
	}
	for k, v := range dictionary {
		i, err := strconv.ParseInt(k, 16, 32)
		if err != nil {
			continue
		}
		pinyinDic[rune(i)] = v
	}
	hasInit = true
}

// 默认带声调，拼音之间空格分隔
func NewPinyin(s string) *Pinyin {
	return &Pinyin{s, WithTone, " "}
}

func (py *Pinyin) Fmt(fmt Fmt) *Pinyin {
	py.fmt = fmt
	return py
}

func (py *Pinyin) Convert() string {
	if !hasInit {
		return ""
	}
	cnRunes := []rune(py.chinese)
	pinyins := make([]string, 0)
	var temp string
	for i, cnRune := range cnRunes {
		_, ok := pinyinDic[cnRune]
		if !ok {
			temp += string(cnRune)
			if i == len(cnRunes)-1 {
				pinyins = append(pinyins, temp)
			}
			continue
		}
		pinyin := convertPinyin(cnRune, py.fmt)
		if len(temp) > 0 {
			pinyins = append(pinyins, temp)
			temp = ""
		}
		if len(pinyin) > 0 {
			pinyins = append(pinyins, pinyin)
		}
	}
	return strings.Join(pinyins, py.split)
}

// 转成拼音
func convertPinyin(cn rune, fmt Fmt) string {
	py := pinyinDic[cn]
	if py == "" {
		return py
	}
	switch fmt {
	case NoTone:
		output := make([]rune, utf8.RuneCountInString(py))
		count := 0
		for _, t := range py {
			neutral, ok := toneMap[t]
			if ok {
				output[count] = neutral
			} else {
				output[count] = t
			}
			count++
		}
		return string(output)
	case FirstUpper:
		sr := []rune(py)
		if sr[0] > 32 {
			sr[0] = sr[0] - 32
		}
		return string(sr)
	default:
		return py
	}
}
