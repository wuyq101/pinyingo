package pinyingo

import (
	"encoding/json"
	"fmt"
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

// 常量
const (
	StyleNormal      = 1
	StyleTone        = 2
	StyleInitials    = 3
	StyleFirstLetter = 4
	UseSegment       = true
	NoSegment        = false
	useHmm           = true
)

var (
	dict        []string
	phrasesDict map[string]string
	reg         *regexp.Regexp
	initials    = strings.Split("b,p,m,f,d,t,n,l,g,k,h,j,q,x,r,zh,ch,sh,z,c,s", ",")
	keyString   string
	jieba       *gojieba.Jieba
	sympolMap   = map[string]string{
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
)

func init() {
	keyString = getMapKeys()
	reg = regexp.MustCompile("([" + keyString + "])")

	//初始化时将gojieba实例化到内存
	// jieba = gojieba.NewJieba()

	//初始化多音字到内存
	file := path.Join(os.Getenv("GOPATH"), "src/github.com/wuyq101/pinyingo/dict/phrases.txt")
	loadPhrases(file)
	//初始化字表
	dict = make([]string, 180000)
	file = path.Join(os.Getenv("GOPATH"), "src/github.com/wuyq101/pinyingo/dict/zi.txt")
	loadZi(file)
}

// Init 用户指定多音词表和汉字库表文件
func Init(confDir string) {
	//清空默认，重新加载
	for i := 0; i < len(dict); i++ {
		dict[i] = ""
	}

	jieba = gojieba.NewJieba(confDir+"/jieba.dict.utf8", confDir+"/hmm_model.utf8", confDir+"/user.dict.utf8", confDir+"/idf.utf8", confDir+"/stop_words.utf8")

	loadZi(confDir + "/zi.txt")
	phrasesDict = make(map[string]string)
	loadPhrases(confDir + "/phrases.txt")
}

func get(index int) string {
	if index >= 0 && index < 180000 {
		return dict[index]
	}
	return ""
}

func exist(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

func loadZi(file string) {
	if !exist(file) {
		fmt.Fprintf(os.Stderr, "character file %s does not exist\n", file)
		return
	}
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

func loadPhrases(file string) {
	if !exist(file) {
		fmt.Fprintf(os.Stderr, "phrases file %s does not exist\n", file)
		return
	}
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&phrasesDict); err != nil {
		log.Fatal(err)
	}
}

// Py 拼音对象
type Py struct {
	style     int
	segment   bool
	heteronym bool
}

func (py *Py) perStr(pinyinStrs string) string {
	switch py.style {
	case StyleInitials:
		for i := 0; i < len(initials); i++ {
			if strings.Index(pinyinStrs, initials[i]) == 0 {
				return initials[i]
			}
		}
		return ""
	case StyleTone:
		ret := strings.Split(pinyinStrs, ",")
		return ret[0]
	case StyleNormal:
		ret := strings.Split(pinyinStrs, ",")
		return normalStr(ret[0])
	case StyleFirstLetter:
		ret := strings.Split(pinyinStrs, ",")
		return firstLetter(ret[0])
	}
	return ""
}

func (py *Py) doConvert(strs string) []string {
	//获取字符串的长度
	bytes := []byte(strs)
	pinyinArr := make([]string, 0)
	nohans := ""
	var tempStr string
	var single string
	for len(bytes) > 0 {
		r, size := utf8.DecodeRune(bytes)
		bytes = bytes[size:]
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
			pinyinArr = append(pinyinArr, py.perStr(single))
		}
	}
	//处理末尾非中文的字符串
	if len(nohans) > 0 {
		pinyinArr = append(pinyinArr, nohans)
	}
	return pinyinArr
}

// Convert 根据配置返回拼音字符串
func (py *Py) Convert(strs string) []string {
	if !py.segment {
		return py.doConvert(strs)
	}
	retArr := make([]string, 0)
	jiebaed := jieba.Cut(strs, useHmm)
	for _, item := range jiebaed {
		mapValuesStr, exist := phrasesDict[item]
		mapValuesArr := strings.Split(mapValuesStr, ",")
		if exist {
			for _, v := range mapValuesArr {
				retArr = append(retArr, py.perStr(v))
			}
		} else {
			converted := py.doConvert(item)
			for _, v := range converted {
				retArr = append(retArr, v)
			}
		}
	}
	return retArr
}

// NewPy 返回一个拼音配置对象
func NewPy(style int, segment bool) *Py {
	return &Py{style, segment, false}
}
