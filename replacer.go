package main

import (
	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
)

type Replacer interface {
	Replace(content []byte) []byte
}

func (r ReplacePair) Replace(content []byte) []byte {
	re := pcre.MustCompile(r.matchBy, 0)
	return re.ReplaceAll(content, []byte(r.replaceWith), 0)
}

func (replacePairArr ReplacePairs) Replace(content []byte) []byte {
	newContent := content
	var r Replacer
	for _, replacePair := range replacePairArr {
		r = replacePair
		newContent = r.Replace(newContent)
	}
	return newContent
}

type ReplacePair struct {
	matchBy, replaceWith string
}

type ReplacePairs []ReplacePair
