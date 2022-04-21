package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func GetCallerIdentity(profileName string) string {

	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profileName), config.WithRegion("eu-west-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := sts.NewFromConfig(cfg)

	resp, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatalln(err)
	}

	return *resp.Account
}
