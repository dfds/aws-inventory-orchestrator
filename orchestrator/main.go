package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dfds/aws-inventory-orchestrator/orchestrator/aws"
	"github.com/dfds/aws-inventory-orchestrator/orchestrator/k8s"
)

func main() {

	callerIdentity := aws.GetCallerIdentity()
	fmt.Println(callerIdentity)

	includeAccountIds := strings.Split(os.Args[1], ",")
	acct, err := aws.OrgAccountList(includeAccountIds)

	if err != nil {
		fmt.Println("%v\n", err)
	} else {
		for _, v := range acct {
			fmt.Println(*v.Id)
		}
	}

	k8s.TestFunc()
}
