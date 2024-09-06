package utils

import (
	"log"

	"github.com/abadojack/whatlanggo"
)

type TransLang struct {
	From string
	To   string
}

func AutoCheckLang(origin *string) TransLang {
	info := whatlanggo.Detect(*origin)
	log.Printf("Detected language: %v", info.Lang)
	if lang := info.Lang; lang == whatlanggo.Cmn {
		return TransLang{From: "zh_CN", To: "en_US"}
	}
	return TransLang{From: "en_US", To: "zh_CN"}
}
