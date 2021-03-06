package main

import (
	"flag"
	"log"

	"github.com/dfds/aws-inventory-orchestrator/infrastructure/aws"
	"github.com/dfds/aws-inventory-orchestrator/infrastructure/util"
)

func main() {

	cliBucketKeyPrefix := "aws/iam/inventory-role/"

	// Parse arguments
	billingAwsProfile := flag.String("billing-aws-profile", "", "Name of AWS profile with admin rights to Billing account.")
	securityAwsProfile := flag.String("security-aws-profile", "", "Name of AWS profile with admin rights to Security account.")
	cliBucketName := flag.String("cli-bucket-name", "", "S3 bucket for uploading \"inventory\" role files (most exist).")
	inventoryBucketName := flag.String("inventory-bucket-name", "", "S3 bucket for inventory reports (will be deployed in Security account).")
	inventoryRole := flag.String("inventory-role", "", "Name of the \"inventory\" IAM role (that needs to be deployed in all accounts).")
	orchRole := flag.String("orchestrator-role", "", "Name of the \"Inventory-Orchestrator\" IAM role (that will be deployed in the Billing account).")
	runnerRole := flag.String("runner-role", "", "Name of the \"Inventory-Runner\" IAM role (that will be deployed in the Security account).")
	oidcProvider := flag.String("oidc-provider-prod", "", "The ID of the OIDC provider for Production use.")

	flag.Parse()

	// Validate all arguments were passed
	if *billingAwsProfile == "" {
		log.Fatal("Required flag \"billing-aws-profile\" not set or empty.")
	}
	if *securityAwsProfile == "" {
		log.Fatal("Required flag \"security-aws-profile\" not set or empty.")
	}
	if *cliBucketName == "" {
		log.Fatal("Required flag \"cli-bucket-name\" not set or empty.")
	}
	if *inventoryBucketName == "" {
		log.Fatal("Required flag \"inventory-bucket-name\" not set or empty.")
	}
	if *inventoryRole == "" {
		log.Fatal("Required flag \"inventory-role\" not set or empty.")
	}
	if *orchRole == "" {
		log.Fatal("Required flag \"orchestrator-role\" not set or empty.")
	}
	if *runnerRole == "" {
		log.Fatal("Required flag \"runner-role\" not set or empty.")
	}
	if *oidcProvider == "" {
		log.Fatal("Required flag \"oidc-provider-prod\" not set or empty.")
	}

	// Get account IDs
	billingAccountId := aws.GetCallerIdentity(*billingAwsProfile)
	securityAccountId := aws.GetCallerIdentity(*securityAwsProfile)

	// Define template tokens for replacement
	emptyTemplateTokens := util.TemplateTokens{}
	fallbackTemplateTokens := util.TemplateTokens{AccountId: securityAccountId}

	/* ORCHESTRATOR ROLES */

	// Create PROD inventory orchestrator role
	orchTrustTokens := util.TemplateTokens{AccountId: billingAccountId, Oidc: *oidcProvider, ServiceAccount: "aws-inventory-orchestrator-sa"}
	orchTrustDoc := util.ParseTemplateFile("./policies/orchestrator_trust.json", orchTrustTokens)
	orchPolicyDoc := util.ParseTemplateFile("./policies/orchestrator_policy.json", emptyTemplateTokens)
	aws.IamCreateRole(*billingAwsProfile, *orchRole, "", orchTrustDoc, orchPolicyDoc, 3600)

	// Create TEST inventory orchestrator role (with no trust policy, as it will be managed manually)
	orchRoleNameTest := *orchRole + "-Test"
	orchTrustDocTest := util.ParseTemplateFile("./policies/fallback_trust.json", fallbackTemplateTokens)
	aws.IamCreateRole(*billingAwsProfile, orchRoleNameTest, "", orchTrustDocTest, orchPolicyDoc, 3600)

	/* RUNNER ROLES */

	// Create PROD inventory runner role
	runnerTrustTokens := util.TemplateTokens{AccountId: securityAccountId, Oidc: *oidcProvider, ServiceAccount: "aws-inventory-runner-sa"}
	runnerTrustDoc := util.ParseTemplateFile("./policies/runner_trust.json", runnerTrustTokens)
	runnerPolicyTokens := util.TemplateTokens{BucketName: *inventoryBucketName, InventoryRole: *inventoryRole}
	runnerPolicyDoc := util.ParseTemplateFile("./policies/runner_policy.json", runnerPolicyTokens)
	aws.IamCreateRole(*securityAwsProfile, *runnerRole, "", runnerTrustDoc, runnerPolicyDoc, 3600)

	// Create TEST inventory runner role (with no trust policy, as it will be managed manually)
	runnerRoleNameTest := *runnerRole + "-Test"
	runnerTrustDocTest := util.ParseTemplateFile("./policies/fallback_trust.json", fallbackTemplateTokens)
	aws.IamCreateRole(*securityAwsProfile, runnerRoleNameTest, "", runnerTrustDocTest, runnerPolicyDoc, 3600)

	// Create inventory runner role
	aws.S3CreateBucket(*securityAwsProfile, *inventoryBucketName)

	/* UPLOAD CLI FILES TO S3 BUCKET */

	inventoryTrustTokens := util.TemplateTokens{AccountId: securityAccountId}
	inventoryTrustDoc := util.ParseTemplateFile("./policies/inventory_trust.json", inventoryTrustTokens)
	inventoryPolicyDoc := util.ParseTemplateFile("./policies/inventory_policy.json", emptyTemplateTokens)
	inventoryPropertiesDoc := util.ParseTemplateFile("./policies/inventory_properties.json", emptyTemplateTokens)
	aws.UploadStringToS3File(*securityAwsProfile, *cliBucketName, cliBucketKeyPrefix+"trust.json", inventoryTrustDoc)
	aws.UploadStringToS3File(*securityAwsProfile, *cliBucketName, cliBucketKeyPrefix+"policy.json", inventoryPolicyDoc)
	aws.UploadStringToS3File(*securityAwsProfile, *cliBucketName, cliBucketKeyPrefix+"properties.json", inventoryPropertiesDoc)
}
