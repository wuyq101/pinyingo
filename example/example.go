package main

import (
	"fmt"

	"github.com/wuyq101/pinyingo"
	"github.com/yanyiwu/gojieba"
)

func main() {
	py := pinyingo.NewPy(pinyingo.StyleNormal, pinyingo.UseSegment)
	strs := []string{"苹果电脑", "厦门港务", "厦门钨业", "厦门空港", "厦门信达", "厦门", "厦门国贸", "我爱北京天安门",
		"西藏", "西藏城投", "西藏矿业", "西藏旅游",
		"西藏天路",
		"西藏药业",
		"瑞贝卡",
		"天津磁卡",
		"厦华电子",
		"任子行",
		"金卡智能",
		"会稽山",
		"做一天和尚撞一天钟",
		"重庆百货",
		"奇正藏药",
	}
	jieba := gojieba.NewJieba()
	for _, s := range strs {
		r := py.Convert(s)
		fmt.Println(r)
		a := jieba.Cut(s, true)
		fmt.Println(a)
		wordinfos := jieba.Tokenize(s, gojieba.SearchMode, true)
		fmt.Println(s)
		fmt.Println("Tokenize:(搜索引擎模式)", wordinfos)
	}
	/*
		buf, _ := ioutil.ReadFile("../dict/zi.txt")
		lines := strings.Split(string(buf), "\n")
		for _, line := range lines {
			idx := strings.Index(line, "=")
			code := line[0 : idx-1]
			v, _ := strconv.ParseInt(code, 16, 64)
			start := strings.Index(line, "\"")
			end := strings.LastIndex(line, "\"")
			pinyin := line[start+1 : end]
			r := rune(v)
	*/
	//	fmt.Printf("%x = \"%s\" /* %c */\n", v, pinyin, r)
	//	}
}
