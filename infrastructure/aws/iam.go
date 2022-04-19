package aws

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func IamNewClient(profileName string) *iam.Client {

	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profileName), config.WithRegion("eu-west-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := iam.NewFromConfig(cfg)

	return client

}

func IamCreateRole(awsProfile string, name string, description string, trustPolicy string, policy string, maxSessionDuration int32) {

	ctx := context.TODO()

	// create new client
	client := IamNewClient(awsProfile)

	// define input for the role creation
	var input *iam.CreateRoleInput = &iam.CreateRoleInput{
		RoleName:                 &name,
		Description:              &description,
		AssumeRolePolicyDocument: &trustPolicy,
		MaxSessionDuration:       &maxSessionDuration,
	}

	// try to create the required role
	_, err := client.CreateRole(ctx, input)

	// in the case of an error
	if err != nil {
		var eae *types.EntityAlreadyExistsException
		if errors.As(err, &eae) {

			if strings.HasSuffix(name, "-Test") {
				fmt.Printf("Not setting trust policy for role %s (probably expected)\n", name)
			} else {
				// if the role already existed then at least ensure the AssumeRolePolicyDocument is updated
				_, err = client.UpdateAssumeRolePolicy(ctx, &iam.UpdateAssumeRolePolicyInput{PolicyDocument: &trustPolicy, RoleName: &name})

				// display errors if any occurred
				if err != nil {
					fmt.Printf("Error updating trust policy for role %s\n", name)
				}
			}

		} else {
			fmt.Printf("Error creating role %s:\n%v\n", name, err)
		}
	} else {
		fmt.Printf("Created role %s\n", name)
	}

	// Attach inline policies
	IamPutRolePolicy(client, name, "ListOrgAccounts", policy)

}

func IamPutRolePolicy(client *iam.Client, roleName string, policyName string, policy string) {

	// define input for the policy put request
	var input *iam.PutRolePolicyInput = &iam.PutRolePolicyInput{
		RoleName:       &roleName,
		PolicyName:     &policyName,
		PolicyDocument: &policy,
	}

	// put inline policy
	_, err := client.PutRolePolicy(context.TODO(), input)
	if err != nil {
		fmt.Println(" There was a problem whilst trying to create the inline policy.")
		fmt.Printf(" The error was: %v\n", err)
	}

}
