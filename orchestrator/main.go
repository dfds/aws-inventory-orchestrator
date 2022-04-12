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
				VolumeName:         "aws-creds",
				VolumePath:         "/aws",
				InitName:           "auth",
				InitImage:          k8s.GetPodImageName(),
				InitCmd:            []string{"./app/runner"},
				InitArgs:           []string{roleArn},
				ContainerName:      "inventory",
				ContainerImage:     "darkbitio/aws_recon:latest",
				// ContainerCmd:       []string{"aws_recon"},
				// ContainerArgs:      []string{"-v", "-r", "global,eu-west-1,eu-central-1"},
				ContainerCmd:  []string{"/bin/sh", "-c", "--"},
				ContainerArgs: []string{"sleep 3600"},
			}

			k8s.CreateJob(&jobSpec)
		}
	}

}
