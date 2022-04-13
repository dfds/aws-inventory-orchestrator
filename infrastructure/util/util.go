package util

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"
)

type TemplateTokens struct {
	AccountId      string
	Oidc           string
	ServiceAccount string
	BucketName     string
	InventoryRole  string
}

type EnvVars struct {
	BillingAccountId string
	RunnerAccountId  string
	OidcProvider     string
	BucketName       string
	InventoryRole    string
	OrchestratorRole string
	RunnerRole       string
}

func MapEnvs(envs []string) map[string]string {

	envMap := make(map[string]string)

	for _, v := range envs {
		envMap[v] = os.Getenv(v)
		if envMap[v] == "" {
			log.Fatalf("Environment varaible %s not set.\n", v)
		}
	}

	return envMap
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
