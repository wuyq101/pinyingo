package pinyingo

import (
	"testing"
)

func TestInit(t *testing.T) {
	Init("dict/phrases.txt", "dict/zi.txt")
}

func TestConvert(t *testing.T) {
	pinyinMap := map[string][]string{
		"hello中国": []string{"hello", "zhōng", "guó"},
		"中国hello": []string{"zhōng", "guó", "hello"},
		"中国":      []string{"zhōng", "guó"},
		"123qwe":  []string{"123qwe"},
	}
	pinyinMap1 := map[string][]string{
		"hello中国": []string{"hello", "zh", "g"},
		"中国":      []string{"zh", "g"},
		"123qwe":  []string{"123qwe"},
	}
	pinyinMap2 := map[string][]string{
		"hello中国": []string{"hello", "zhong", "guo"},
		"中国":      []string{"zhong", "guo"},
		"123qwe":  []string{"123qwe"},
	}

	pinyinMap3 := map[string][]string{
		"重阳": []string{"chóng", "yáng"},
		"重点": []string{"zhòng", "diǎn"},
	}

	pinyinMap4 := map[string][]string{
		"你好啊": []string{"n", "h", "a"},
	}

	for k, v := range pinyinMap {
		py := NewPy(StyleTone, NoSegment)
		converted := py.Convert(k)
		for i := 0; i < len(converted); i++ {
			if converted[i] != v[i] {
				t.Errorf("%s is not equal %s", converted, v)
			}
		}
	}

	for k, v := range pinyinMap1 {
		py := NewPy(StyleInitials, NoSegment)
		converted := py.Convert(k)
		for i := 0; i < len(converted); i++ {
			if converted[i] != v[i] {
				t.Errorf("%s is not equal %s", converted, v)
			}
		}
	}

	for k, v := range pinyinMap2 {
		py := NewPy(StyleNormal, NoSegment)
		converted := py.Convert(k)
		for i := 0; i < len(converted); i++ {
			if converted[i] != v[i] {
				t.Errorf("%s is not equal %s", converted, v)
			}
		}
	}

	for k, v := range pinyinMap3 {
		py := NewPy(StyleTone, UseSegment)
		converted := py.Convert(k)
		for i := 0; i < len(converted); i++ {
			if converted[i] != v[i] {
				t.Errorf("%s is not equal %s", converted, v)
			}
		}
	}

	for k, v := range pinyinMap4 {
		py := NewPy(StyleFirstLetter, UseSegment)
		converted := py.Convert(k)
		for i := 0; i < len(converted); i++ {
			if converted[i] != v[i] {
				t.Errorf("%s is not equal %s", converted, v)
			}
		}
	}
}
