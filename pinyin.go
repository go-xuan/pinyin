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
)

type Pinyin struct {
	withTone   bool
	firstUpper bool
	chinese    string
	split      string
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
		if i, err := strconv.ParseInt(k, 16, 32); err != nil {
			continue
		} else {
			pinyinDic[rune(i)] = v
		}
	}
}

// NewPinyin 默认带声调，拼音之间空格分隔
func NewPinyin(cn string, split ...string) *Pinyin {
	var py = &Pinyin{chinese: cn}
	if len(split) > 0 && len(split[0]) > 0 {
		py.split = split[0]
	}
	return py
}

func (p *Pinyin) WithTone() *Pinyin {
	p.withTone = true
	return p
}

func (p *Pinyin) FirstUpper() *Pinyin {
	p.firstUpper = true
	return p
}

func (p *Pinyin) Convert() string {
	runes := []rune(p.chinese)
	pinyins := make([]string, 0)
	var pinyin string
	for i, cnRune := range runes {
		if _, ok := pinyinDic[cnRune]; !ok {
			pinyin += string(cnRune)
			if i == len(runes)-1 {
				pinyins = append(pinyins, pinyin)
			}
			continue
		}
		if len(pinyin) > 0 {
			pinyins = append(pinyins, pinyin)
			pinyin = ""
		}
		if py := convertPinyin(cnRune, p.withTone, p.firstUpper); len(py) > 0 {
			pinyins = append(pinyins, py)
		}
	}
	return strings.Join(pinyins, p.split)
}

// 转成拼音
func convertPinyin(cn rune, withTone, firstUpper bool) string {
	pinyin := pinyinDic[cn]
	if pinyin != "" {
		if !withTone {
			runes, i := make([]rune, utf8.RuneCountInString(pinyin)), 0
			for _, t := range pinyin {
				if tone, ok := toneMap[t]; ok {
					runes[i] = tone
				} else {
					runes[i] = t
				}
				i++
			}
			pinyin = string(runes)
		}
		if firstUpper {
			if sr := []rune(pinyin); sr[0] > 32 {
				sr[0] = sr[0] - 32
				pinyin = string(sr)
			}
		}
	}
	return pinyin
}
