package main

import (
	"fmt"

	"github.com/dfds/aws-inventory-orchestrator/orchestrator/aws"
)

func main() {

	callerIdentity := aws.GetCallerIdentity()
	fmt.Println(callerIdentity)

	var includeAccountIds []string
	acct, err := aws.OrgAccountList(includeAccountIds)

	if err != nil {
		fmt.Println("%v\n", err)
	} else {
		for _, v := range acct {
			fmt.Println(*v.Id)
		}
	}

}
