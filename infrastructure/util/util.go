package util

import (
	"bytes"
	"fmt"
	"text/template"
)

type TemplateTokens struct {
	AccountId      string
	Oidc           string
	ServiceAccount string
	BucketName     string
	InventoryRole  string
}

func ParseTemplateFile(fileName string, tokens TemplateTokens) string {
	buf := new(bytes.Buffer)
	temp := template.Must(template.ParseFiles(fileName))
	err := temp.Execute(buf, tokens)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return buf.String()
}
