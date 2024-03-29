package FingerRules

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
)

type Rule struct {
	Finger []struct {
		Name    string   `json:"name"`
		Keyword []string `json:"keyword"`
	}
}

func Matchkeyword(filepath string, respBody string) (string, error) {
	//读取指纹规则文件
	rules, err := ioutil.ReadFile(filepath)
	//fmt.Println("成功读取json文件")
	if err != nil {
		return "", err
	}

	//解析json规则文件
	var rule Rule
	err = json.Unmarshal(rules, &rule)
	if err != nil {
		return "", err
	}

	//匹配关键词
	b := false
	for _, finger := range rule.Finger {
		for _, keyword := range finger.Keyword {
			b, err = matchKeyword(keyword, respBody)
			if b {
				return keyword, err
			} else {
				break
			}
		}
	}
	return "", err
}

func matchKeyword(keyword string, content string) (bool, error) {
	re, err := regexp.Compile(keyword)
	if err != nil {
		return false, err
	}

	return re.MatchString(content), nil
}
