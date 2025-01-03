package pinyin

import (
	"fmt"
	"testing"
)

func TestPinyin(t *testing.T) {
	fmt.Println(NewPinyin("你好").Convert())
	fmt.Println(NewPinyin("你好", " ").WithTone().Convert())
	fmt.Println(NewPinyin("你好", "-").FirstUpper().Convert())
	fmt.Println(NewPinyin("你好", "-").WithTone().FirstUpper().Convert())
}
