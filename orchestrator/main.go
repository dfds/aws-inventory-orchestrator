package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dfds/aws-inventory-orchestrator/orchestrator/aws"
	"github.com/dfds/aws-inventory-orchestrator/orchestrator/k8s"
)

func main() {

	// display caller identity; useful during the debug phase
	callerIdentity := aws.GetCallerIdentity()
	fmt.Println(callerIdentity)

	imageTag := "513891c-dirty"
	jobName := "inventory-runner"
	jobNamespace := "inventory"
	cmd := "./app/runner"

	image := fmt.Sprintf("dfdsdk/aws-inventory-orchestrator:%s", imageTag)

	// get accounts to target for inventory
	includeAccountIds := strings.Split(os.Args[1], ",")
	acct, err := aws.OrgAccountList(includeAccountIds)

	if err != nil {
		fmt.Println("%v\n", err)
	} else {
		for _, v := range acct {
			fmt.Println(*v.Id)

			roleArn := fmt.Sprintf("arn:aws:iam::%s:role/inventory", *v.Id)
			k8s.CreateJob(&jobName, &jobNamespace, &image, &cmd, &roleArn)
		}
	}

}
