package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type TrustTemplate struct {
	AccountId      string
	Oidc           string
	ServiceAccount string
}

func main() {

	// To create:
	// In Billing/Master
	//  Inventory-Orchestrator IAM role (in billing/master)
	//
	// In DFDS-Security
	//  Inventory-Runner IAM role (in dfds-security)
	//  Inventory S3 bucket (in dfds-security)

	// Retrieve environment variables
	billingAccountId := os.Getenv("BILLING_ACCOUNT_ID")
	if billingAccountId == "" {
		log.Fatal("BILLING_ACCOUNT_ID not specified.")
	}

	runnerAccountId := os.Getenv("RUNNER_ACCOUNT_ID")
	if runnerAccountId == "" {
		log.Fatal("RUNNER_ACCOUNT_ID not specified.")
	}

	oidcProvider := os.Getenv("OIDC_PROVIDER")
	if oidcProvider == "" {
		log.Fatal("OIDC_PROVIDER not specified.")
	}

	bucketName := os.Getenv("BUCKET_NAME")
	if bucketName == "" {
		log.Fatal("BUCKET_NAME not specified.")
	}

	inventoryRole := os.Getenv("INVENTORY_ROLE")
	if inventoryRole == "" {
		log.Fatal("INVENTORY_ROLE not specified.")
	}

	inventoryOrchestratorRole := os.Getenv("INVENTORY_ORCHESTRATOR_ROLE")
	if inventoryOrchestratorRole == "" {
		log.Fatal("INVENTORY_ORCHESTRATOR_ROLE not specified.")
	}

	inventoryRunnerRole := os.Getenv("INVENTORY_RUNNER_ROLE")
	if inventoryRunnerRole == "" {
		log.Fatal("INVENTORY_RUNNER_ROLE not specified.")
	}

	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("billing-admin"), config.WithRegion("eu-west-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := sts.NewFromConfig(cfg)
	identity, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Account: %s, Arn: %s\n", aws.ToString(identity.Account), aws.ToString(identity.Arn))

	// declare buffer for templated data
	buf := new(bytes.Buffer)

	orchTrustTemplate := template.Must(template.ParseFiles("./policies/orchestrator_trust.json"))
	orchestratorTrust := TrustTemplate{AccountId: billingAccountId, Oidc: oidcProvider, ServiceAccount: "aws-inventory-orchestrator-sa"}
	err = orchTrustTemplate.Execute(buf, orchestratorTrust)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	trustPolicy := buf.String()
	fmt.Println("Trust Policy: ", trustPolicy)

	orchPolicyTemplate := template.Must(template.ParseFiles("./policies/orchestrator_policy.json"))
	orchestratorPolicy := TrustTemplate{}
	err = orchPolicyTemplate.Execute(buf, orchestratorPolicy)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	policy := buf.String()
	fmt.Println("Policy: ", policy)

	// func CreateIAMRole(client *iam.Client, name string, description string, policy string, trustPolicy string, maxSessionDuration int32) {
	// aws.CreateIAMRole(client, inventoryOrchestratorRole, "Inventory Orchestrator Role", )

	// cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("security-cloudadmin"), config.WithRegion("eu-west-1"))
	// if err != nil {
	// 	log.Fatalf("unable to load SDK config, %v", err)
	// }

	// client = sts.NewFromConfig(cfg)
	// identity, err = client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("Account: %s, Arn: %s\n", aws.ToString(identity.Account), aws.ToString(identity.Arn))

}
