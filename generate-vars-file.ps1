$ErrorActionPreference = "Stop"
$VarsFile = "./k8s/vars.env"

"BILLING_ACCOUNT_ID=$(Get-Item "Env:/BILLING_ACCOUNT_ID" | Select -Expand Value)
SECURITY_ACCOUNT_ID=$(Get-Item "Env:/SECURITY_ACCOUNT_ID" | Select -Expand Value)
ORCHESTRATOR_ROLE=$(Get-Item "Env:/ORCHESTRATOR_ROLE" | Select -Expand Value)
RUNNER_ROLE=$(Get-Item "Env:/RUNNER_ROLE" | Select -Expand Value)
BUCKET_NAME=$(Get-Item "Env:/BUCKET_NAME" | Select -Expand Value)
CRON_SCHEDULE=$(Get-Item "Env:/CRON_SCHEDULE" | Select -Expand Value)
INCLUDE_ACCOUNTS=$(Get-Item "Env:/INCLUDE_ACCOUNTS" | Select -Expand Value)" | Out-File $VarsFile