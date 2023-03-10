package FingerRules

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strings"
)

type Rule struct {
	Matchers struct {
		Key []string `json:"keyword"`
	} `json:"matchers"`
}

func Matchkeyword(filepath string, linkBody string) (bool, error) {
	//读取指纹规则文件
	rules, err := ioutil.ReadFile(filepath)
	if err != nil {
		return false, err
	}
	//解析json规则文件
	var rule Rule
	err = json.Unmarshal(rules, &rule)
	if err != nil {
		return false, err
	}

	//匹配关键词
	for _, keyword := range rule.Matchers.Key {
		b, err := matchKeyword(keyword, linkBody)
		if err != nil {
			panic(err)
		}
	}
	return b, err
}

func matchKeyword(keyword string, content string) (bool, error) {
	re, err := regexp.Compile(keyword)
	if err != nil {
		return false, err
	}

	return re.MatchString(content), nil
}
