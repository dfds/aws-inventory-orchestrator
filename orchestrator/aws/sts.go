package aws

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

func GetCallerIdentity() string {

	sess, err := session.NewSession(&aws.Config{Region: aws.String("eu-west-1")})
	if err != nil {
		log.Fatalln(err)
	}

	stsCli := sts.New(sess)

	resp, err := stsCli.GetCallerIdentity(nil)
	if err != nil {
		log.Fatalln(err)
	}

	//fmt.Println(resp.String())
	return resp.String()
}
