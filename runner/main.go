package main

import (
	"fmt"
	"os"

	"github.com/dfds/aws-inventory-orchestrator/runner/aws"
)

func main() {

	fmt.Println("This is a job spawned by the Inventory-Orchestrator.")

	tempCreds, err := aws.AssumeRole(os.Args[1])
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(tempCreds)

}
