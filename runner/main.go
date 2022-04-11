package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dfds/aws-inventory-orchestrator/runner/aws"
)

func main() {

	assumedCreds, err := aws.AssumeRole(os.Args[1])
	if err != nil {
		panic(err.Error())
	}

	awsCredsString := fmt.Sprintf("[default]\naws_access_key_id = %s\naws_secret_access_key = %s\naws_session_token = %s\n", *assumedCreds.AccessKeyId, *assumedCreds.SecretAccessKey, *assumedCreds.SessionToken)

	byteSlice := []byte(awsCredsString)             // convert string to byte slice
	ioutil.WriteFile("/aws/creds", byteSlice, 0644) // the 0644 is octal representation of the filemode

}
