package aws

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	BucketCannedACLPrivate     types.BucketCannedACL          = "private"
	BucketLocationConstraintEu types.BucketLocationConstraint = "EU"
	ServerSideEncryptionAes256 types.ServerSideEncryption     = "AES256"
	// ServerSideEncryptionAwsKms types.ServerSideEncryption     = "aws:kms"
)

func S3NewClient(profileName string) *s3.Client {

	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profileName), config.WithRegion("eu-west-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := s3.NewFromConfig(cfg)

	return client

}

func S3CreateBucket(awsProfile string, name string) {

	ctx := context.TODO()
	_ = ctx

	client := S3NewClient(awsProfile)

	// s3 bucket input
	var bucketInput *s3.CreateBucketInput = &s3.CreateBucketInput{
		Bucket: &name,
		ACL:    BucketCannedACLPrivate,
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: BucketLocationConstraintEu,
		},
	}

	// Create bucket
	_, err := client.CreateBucket(ctx, bucketInput)
	if err != nil {
		var eae *types.BucketAlreadyOwnedByYou
		if errors.As(err, &eae) {
			fmt.Printf("Bucket %s already exists (and is owned by you)\n", name)
		} else {
			fmt.Printf("Failed to create %s bucket\n%s\n", name, err)
		}
	}

	// This is broken, needs to happen
	// See https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/s3-example-default-server-side-encryption.html (v1 though)

	// // Bucket encryption input
	// defEnc := &types.ServerSideEncryptionByDefault{SSEAlgorithm: ServerSideEncryptionAes256}
	// encRule := &types.ServerSideEncryptionRule{
	// 	ApplyServerSideEncryptionByDefault: defEnc,
	// 	// BucketKeyEnabled: true,
	// }
	// encRules := []*types.ServerSideEncryptionRule{encRule}
	// var encryptionInput *s3.PutBucketEncryptionInput = &s3.PutBucketEncryptionInput{
	// 	Bucket: &name,
	// 	ServerSideEncryptionConfiguration: &types.ServerSideEncryptionConfiguration{
	// 		Rules: encRules,
	// 	},
	// }

	// // Encrypt bucket objects by default
	// _, err = client.PutBucketEncryption(ctx, encryptionInput)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// Public access block input
	var accessInput *s3.PutPublicAccessBlockInput = &s3.PutPublicAccessBlockInput{
		Bucket: &name,
		PublicAccessBlockConfiguration: &types.PublicAccessBlockConfiguration{
			BlockPublicAcls:       true,
			BlockPublicPolicy:     true,
			IgnorePublicAcls:      true,
			RestrictPublicBuckets: true,
		},
	}

	// block public access
	_, err = client.PutPublicAccessBlock(ctx, accessInput)
	if err != nil {
		fmt.Println(err)
	}

}
