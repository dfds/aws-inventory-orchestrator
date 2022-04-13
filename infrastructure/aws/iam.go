package aws

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func InitAWS() {}

func CreateIAMRole(client *iam.Client, name string, description string, trustPolicy string, policy string, maxSessionDuration int32) {

	// define input for the role creation
	var input *iam.CreateRoleInput = &iam.CreateRoleInput{
		MaxSessionDuration:       &maxSessionDuration,
		AssumeRolePolicyDocument: &trustPolicy,
		// Path:                     &path,
		RoleName:    &name,
		Description: &description,
	}

	// try to create the required role
	_, err := client.CreateRole(context.TODO(), input)

	// in the case of an error
	if err != nil {
		var eae *types.EntityAlreadyExistsException
		if errors.As(err, &eae) {

			// if the role already existed then at least ensure the AssumeRolePolicyDocument is updated
			_, err = client.UpdateAssumeRolePolicy(context.TODO(), &iam.UpdateAssumeRolePolicyInput{PolicyDocument: &trustPolicy, RoleName: &name})

			// display errors if any occurred
			if err != nil {
				fmt.Println("Error during update of role trust policy.")
			}
		} else {
			fmt.Println("Error during creation of role.")
		}
	}

}
