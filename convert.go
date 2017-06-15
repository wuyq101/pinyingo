package pinyingo

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/yanyiwu/gojieba"
)

var (
	STYLE_NORMAL       = 1
	STYLE_TONE         = 2
	STYLE_INITIALS     = 3
	STYLE_FIRST_LETTER = 4
	USE_SEGMENT        = true
	NO_SEGMENT         = false
	use_hmm            = true
	DICT_DIR           = path.Join(os.Getenv("GOPATH"), "src/github.com/wuyq101/Go-pinyin/dict")
	DICT_PHRASES       = path.Join(DICT_DIR, "phrases-dict")
	dict               []string
)

func get(index int) string {
	return dict[index]
}

func init() {
	dict = make([]string, 200000)
	loadZi()
}

func loadZi() {
	file := path.Join(os.Getenv("GOPATH"), "src/github.com/wuyq101/Go-pinyin/dict/zi.txt")
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(buf), "\n")
	for _, line := range lines {
		idx := strings.Index(line, "=")
		code := line[0 : idx-1]
		v, _ := strconv.ParseInt(code, 16, 64)
		start := strings.Index(line, "\"")
		end := strings.LastIndex(line, "\"")
		pinyin := line[start+1 : end]
		dict[v] = pinyin
	}
}

var phrasesDict map[string]string
var reg *regexp.Regexp
var INITIALS []string = strings.Split("b,p,m,f,d,t,n,l,g,k,h,j,q,x,r,zh,ch,sh,z,c,s", ",")
var keyString string
var jieba *gojieba.Jieba
var sympolMap = map[string]string{
	"ā": "a1",
	"á": "a2",
	"ǎ": "a3",
	"à": "a4",
	"ē": "e1",
	"é": "e2",
	"ě": "e3",
	"è": "e4",
	"ō": "o1",
	"ó": "o2",
	"ǒ": "o3",
	"ò": "o4",
	"ī": "i1",
	"í": "i2",
	"ǐ": "i3",
	"ì": "i4",
	"ū": "u1",
	"ú": "u2",
	"ǔ": "u3",
	"ù": "u4",
	"ü": "v0",
	"ǘ": "v2",
	"ǚ": "v3",
	"ǜ": "v4",
	"ń": "n2",
	"ň": "n3",
	"": "m2",
}

func init() {
	keyString = getMapKeys()
	reg = regexp.MustCompile("([" + keyString + "])")

	//初始化时将gojieba实例化到内存
	jieba = gojieba.NewJieba()

	//初始化多音字到内存
	initPhrases()
}

func getMapKeys() string {
	keyString := ""
	for key := range sympolMap {
		keyString += key
	}
	return keyString
}

func normalStr(str string) string {
	findRet := reg.FindString(str)
	if findRet == "" {
		return str
	}
	return strings.Replace(str, findRet, string([]byte(sympolMap[findRet])[0]), -1)
}

func firstLetter(str string) string {
	firstLetter := string(str[0])

	if sympolMap[str] != "" {
		firstLetter = string(sympolMap[str][0])
	}
	return firstLetter
}

func initPhrases() {
	f, err := os.Open(DICT_PHRASES)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&phrasesDict); err != nil {
		log.Fatal(err)
	}
}

type options struct {
	style     int
	segment   bool
	heteronym bool
}

func (opt *options) perStr(pinyinStrs string) string {
	switch opt.style {
	case STYLE_INITIALS:
		for i := 0; i < len(INITIALS); i++ {
			if strings.Index(pinyinStrs, INITIALS[i]) == 0 {
				return INITIALS[i]
			}
		}
		return ""
	case STYLE_TONE:
		ret := strings.Split(pinyinStrs, ",")
		return ret[0]
	case STYLE_NORMAL:
		ret := strings.Split(pinyinStrs, ",")
		return normalStr(ret[0])
	case STYLE_FIRST_LETTER:
		ret := strings.Split(pinyinStrs, ",")
		return firstLetter(ret[0])
	}
	return ""
}

func (opt *options) doConvert(strs string) []string {
	//获取字符串的长度
	bytes := []byte(strs)
	pinyinArr := make([]string, 0)
	nohans := ""
	var tempStr string
	var single string
	for len(bytes) > 0 {
		r, w := utf8.DecodeRune(bytes)
		bytes = bytes[w:]
		single = get(int(r))
		// 中文字符判断
		tempStr = string(r)
		if len(single) == 0 {
			nohans += tempStr
		} else {
			if len(nohans) > 0 {
				pinyinArr = append(pinyinArr, nohans)
				nohans = ""
			}
			pinyinArr = append(pinyinArr, opt.perStr(single))
		}
	}
	//处理末尾非中文的字符串
	if len(nohans) > 0 {
		pinyinArr = append(pinyinArr, nohans)
	}
	return pinyinArr
}
func (opt *options) Convert(strs string) []string {
	retArr := make([]string, 0)
	if opt.segment {
		jiebaed := jieba.Cut(strs, use_hmm)
		for _, item := range jiebaed {
			mapValuesStr, exist := phrasesDict[item]
			mapValuesArr := strings.Split(mapValuesStr, ",")
			if exist {
				for _, v := range mapValuesArr {
					retArr = append(retArr, opt.perStr(v))
				}
			} else {
				converted := opt.doConvert(item)
				for _, v := range converted {
					retArr = append(retArr, v)
				}
			}
		}
	} else {
		retArr = opt.doConvert(strs)
	}

	return retArr
}

func NewPy(style int, segment bool) *options {
	return &options{style, segment, false}
}
