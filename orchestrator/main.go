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

	// get accounts to target for inventory
	includeAccountIds := strings.Split(os.Args[1], ",")
	acct, err := aws.OrgAccountList(includeAccountIds)

	if err != nil {
		fmt.Println("%v\n", err)
	} else {
		for _, v := range acct {
			fmt.Println(*v.Id)

			roleArn := fmt.Sprintf("arn:aws:iam::%s:role/managed/inventory", *v.Id)

			jobSpec := k8s.AssumeJobSpec{
				JobName:            "aws-inventory-runner",
				JobNamespace:       "inventory",
				ServiceAccountName: "aws-inventory-runner-sa",
				InitName:           "auth",
				InitImage:          k8s.GetPodImageName(),
				InitCmd:            []string{"./app/runner"},
				InitArgs:           []string{roleArn},
				ContainerName:      "inventory",
				ContainerImage:     "amazon/aws-cli:latest",
				ContainerCmd:       []string{"/bin/bash", "-c", "--"},
				ContainerArgs:      []string{"aws sts get-caller-identity; sleep 3600"},
				// ContainerName:  "auth",
				// ContainerImage: k8s.GetPodImageName(),
				// ContainerCmd:   []string{"/bin/sh", "-c", "--"},
				// ContainerArgs:  []string{fmt.Sprintf("./app/runner %s; sleep 180", roleArn)},
			}

			k8s.CreateJob(&jobSpec)
		}
	}

}
