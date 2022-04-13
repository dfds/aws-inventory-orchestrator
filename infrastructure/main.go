package main

import (
	"github.com/dfds/aws-inventory-orchestrator/infrastructure/aws"
	"github.com/dfds/aws-inventory-orchestrator/infrastructure/util"
)

const (
	billingAwsProfile string = "billing-admin"
	runnerAwsProfile  string = "inventory-runner-admin"
)

func main() {

	// Retrieve environment variables
	envMap := util.MapEnvs([]string{"BILLING_ACCOUNT_ID", "BUCKET_NAME", "INVENTORY_ROLE", "OIDC_PROVIDER", "INVENTORY_ORCHESTRATOR_ROLE", "RUNNER_ACCOUNT_ID", "INVENTORY_RUNNER_ROLE"})

	// Create inventory orchestrator role
	orchTrustTokens := util.TemplateTokens{AccountId: envMap["BILLING_ACCOUNT_ID"], Oidc: envMap["OIDC_PROVIDER"], ServiceAccount: "aws-inventory-orchestrator-sa"}
	orchTrustDoc := util.ParseTemplateFile("./policies/orchestrator_trust.json", orchTrustTokens)
	orchPolicyTokens := util.TemplateTokens{}
	orchPolicyDoc := util.ParseTemplateFile("./policies/orchestrator_policy.json", orchPolicyTokens)
	aws.IamCreateRole(billingAwsProfile, envMap["INVENTORY_ORCHESTRATOR_ROLE"], "", orchTrustDoc, orchPolicyDoc, 3600)

	// Create inventory runner role
	runnerTrustTokens := util.TemplateTokens{AccountId: envMap["RUNNER_ACCOUNT_ID"], Oidc: envMap["OIDC_PROVIDER"], ServiceAccount: "aws-inventory-runner-sa"}
	runnerTrustDoc := util.ParseTemplateFile("./policies/runner_trust.json", runnerTrustTokens)
	runnerPolicyTokens := util.TemplateTokens{BucketName: envMap["BUCKET_NAME"], InventoryRole: envMap["INVENTORY_ROLE"]}
	runnerPolicyDoc := util.ParseTemplateFile("./policies/runner_policy.json", runnerPolicyTokens)
	aws.IamCreateRole(runnerAwsProfile, envMap["INVENTORY_RUNNER_ROLE"], "", runnerTrustDoc, runnerPolicyDoc, 3600)

	// Create inventory runner role
	aws.S3CreateBucket(runnerAwsProfile, envMap["BUCKET_NAME"])

}
