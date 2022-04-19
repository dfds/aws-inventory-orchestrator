$ErrorActionPreference = "Stop"
$VarsFile = "./k8s/vars.env"

"BILLING_ACCOUNT_ID=$(Get-Item "Env:/BILLING_ACCOUNT_ID" | Select -Expand Value)
INVENTORY_ORCHESTRATOR_ROLE=$(Get-Item "Env:/INVENTORY_ORCHESTRATOR_ROLE" | Select -Expand Value)
RUNNER_ACCOUNT_ID=$(Get-Item "Env:/RUNNER_ACCOUNT_ID" | Select -Expand Value)
INVENTORY_RUNNER_ROLE=$(Get-Item "Env:/INVENTORY_RUNNER_ROLE" | Select -Expand Value)
BUCKET_NAME=$(Get-Item "Env:/BUCKET_NAME" | Select -Expand Value)
CRON_SCHEDULE=$(Get-Item "Env:/CRON_SCHEDULE" | Select -Expand Value)
INCLUDE_ACCOUNTS=$(Get-Item "Env:/INCLUDE_ACCOUNTS" | Select -Expand Value)
INVENTORY_ROLE=$(Get-Item "Env:/INVENTORY_ROLE" | Select -Expand Value)
OIDC_PROVIDER=$(Get-Item "Env:/OIDC_PROVIDER" | Select -Expand Value)" | Out-File $VarsFile