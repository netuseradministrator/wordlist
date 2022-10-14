package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"io"
	"net/http"
	"strconv"
	"strings"
)

type cloze_data struct {
	Cloze    interface{} `json:"cloze"`
	Options  interface{} `json:"options"`
	Syllable interface{} `json:"syllable"`
	Tips     interface{} `json:"tips"`
}
type wordliststruct struct {
	Word            interface{} `json:"word"`
	Accent          interface{} `json:"accent"`
	Mean_cn         string      `json:"mean_Cn"`
	Mean_en         interface{} `json:"mean_En"`
	Sentence        interface{} `json:"sentence"`
	Sentence_trans  interface{} `json:"sentence_trans"`
	Sentence_phrase interface{} `json:"sentence_phrase"`
	Word_etyma      interface{} `json:"word_etyma"`
	Cloze_data      cloze_data  `json:"cloze_data"`
}

var (
	word        string
	result      string
	filepath    string
	wordlist    wordliststruct
	trans_cn    string
	api_address = []string{"https://cdn.jsdelivr.net/gh/lyc8503/baicizhan-word-meaning-API/data/list.json",
		"https://cdn.jsdelivr.net/gh/lyc8503/baicizhan-word-meaning-API/data/words/[WORD].json"}
)

func zhToUnicode(raw []byte) ([]byte, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}
func translate() {
	fmt.Scanln(&word)
	fmt.Printf("查%s是吧\n", word)
	url := fmt.Sprintf("https://cdn.jsdelivr.net/gh/lyc8503/baicizhan-word-meaning-API/data/words/%s.json", word)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36")
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if err != nil {
		println(err)
	}
	body, _ := io.ReadAll(res.Body)
	sText := string(body)
	textQuoted := strconv.QuoteToASCII(sText)
	textUnquoted := textQuoted[1 : len(textQuoted)-1]
	v, _ := zhToUnicode([]byte(textUnquoted))
	//fmt.Println(string(v))
	result = strings.ReplaceAll(string(v), "\\", "")
}
func jsonformat() {
	//println(result)
	var output interface{}
	json.Unmarshal([]byte(result), &output)
	result, err := json.Marshal(output)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(string(result))
	err = json.Unmarshal(result, &wordlist)
	if err != nil {
		fmt.Println(err)
		return
	}
	trans_cn = wordlist.Mean_cn

}
func record() {
	filepath = "D:\\学习资料\\typora\\1.md"
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
	}
	file.WriteString("\n" + word + "  " + trans_cn)
	file.Close()
}
func main() {
	for true {
		translate()
		jsonformat()
		fmt.Printf("这个b单词%v的汉译是\n%v\n", word, trans_cn)
		record()
		time.Sleep(10 * time.Millisecond)
		word, wordlist, result = "", wordliststruct{}, ""
	}
}
