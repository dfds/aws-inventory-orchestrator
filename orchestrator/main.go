package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dfds/aws-inventory-orchestrator/orchestrator/aws"
	"github.com/dfds/aws-inventory-orchestrator/orchestrator/k8s"
)

func main() {

	now := time.Now().UTC()
	s3Bucket := os.Args[1]

	// display caller identity; useful during the debug phase
	fmt.Println("Assumed role:")
	callerIdentity := aws.GetCallerIdentity()
	fmt.Println(callerIdentity)

	// get accounts to target for inventory
	var includeAccountIds []string
	if len(os.Args[2]) > 0 {
		includeAccountIds = strings.Split(os.Args[2], ",")
	}

	acct, err := aws.OrgAccountList(includeAccountIds)
	if err != nil {
		fmt.Println("%v\n", err)
	} else {
		for _, v := range acct {

			roleArn := fmt.Sprintf("arn:aws:iam::%s:role/managed/inventory", *v.Id)
			s3Path := fmt.Sprintf("s3://%s/AWSRecon/%s/recon_%s.json", s3Bucket, now.Format("2006/01/02"), *v.Id)
			uploadArg := fmt.Sprintf("set -u; FILE=/recon/output.json; [ -f \"$FILE\" ] && { echo Uploading \"$FILE\" to \"%s\"; aws s3 cp --no-progress \"$FILE\" %s; } || { echo Output file \"$FILE\" not found!; exit 404; }", s3Path, s3Path)

			fmt.Printf("Creating job using IAM role \"%s\"\n", roleArn)

			jobSpec := k8s.AssumeJobSpec{
				AccountId:          *v.Id,
				JobName:            "aws-inventory-runner",
				JobNamespace:       "inventory",
				ServiceAccountName: "aws-inventory-runner-sa",
				CredsVolName:       "aws-creds",
				CredsVolPath:       "/aws",
				OutputVolName:      "inventory-output",
				OutputVolPath:      "/recon",
				AssumeName:         "assume",
				AssumeImage:        k8s.GetPodImageName(),
				AssumeCmd:          []string{"./app/runner"},
				AssumeArgs:         []string{roleArn},
				InventoryName:      "inventory",
				InventoryImage:     "darkbitio/aws_recon:latest",
				InventoryCmd:       []string{"aws_recon"},
				InventoryArgs:      []string{"-v", "-r", "global,eu-west-1,eu-central-1"},
				UploadName:         "upload",
				UploadImage:        "amazon/aws-cli:latest",
				UploadCmd:          []string{"/bin/sh", "-c", "--"},
				UploadArgs:         []string{uploadArg},
			}

			k8s.CreateJob(&jobSpec)
		}
	}

}
