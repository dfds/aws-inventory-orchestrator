package main

import (
	"fmt"

	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/sts"
	"github.com/dfds/aws-inventory-orchestrator/orchestrator/aws"
)

func main() {

	// sess, err := session.NewSession(&aws.Config{Region: aws.String("eu-central-1")})
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// stsCli := sts.New(sess)

	// resp, err := stsCli.GetCallerIdentity(nil)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Println(resp.String())

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
