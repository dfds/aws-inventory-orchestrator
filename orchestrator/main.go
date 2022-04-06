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

	jobName := "inventory-runner"
	jobNamespace := "inventory"
	//image := "dfdsdk/aws-inventory-orchestrator:latest"
	//cmd := "./runner"
	image := "ubuntu:latest"
	cmd := "ls"

	// kubernetes test which will read pod data
	k8s.TestFunc()

	// kuberenetes test which will try to spawn a new job
	k8s.CreateJob(&jobName, &jobNamespace, &image, &cmd)
}
