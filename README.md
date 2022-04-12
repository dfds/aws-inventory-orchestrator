# inventory-orchestrator

## To do

- Terraform:
  - Inventory-Orchestrator IAM role (in billing/master)
  - Inventory-Runner IAM role (in dfds-security)
  - Inventory S3 bucket (in dfds-security)
- Kubernetes:
  - Add common and different labels to orchestrator and runner jobs
- Containers:
  - Run as non-root
    - Can specify in Dockerfile?
    - CronJob and Job

### Inventory-Orchestrator and Runner trust relationship

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Federated": "arn:aws:iam::$ACCOUNT_ID:oidc-provider/$OIDC"
            },
            "Action": "sts:AssumeRoleWithWebIdentity",
            "Condition": {
                "StringEquals": {
                    "$OIDC:sub": "system:serviceaccount:inventory:$K8S_SA"
                }
            }
        }
    ]
}
```

Replace $OIDC, $ACCOUNT_ID, $K8S_SA (billing and security respectively).

### Inventory-Runner policies

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": "sts:AssumeRole",
            "Resource": "arn:aws:iam::*:role/inventory"
        }
    ]
}
```

### Inventory inline policy

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "access-analyzer:List*",
                "acm:Describe*",
                "apigateway:GET",
                "application-autoscaling:Describe*",
                "athena:Get*",
                "autoscaling:Describe*",
                "backup:List*",
                "cloudtrail:List*",
                "cloudwatch:Describe*",
                "codebuild:List*",
                "config:Describe*",
                "dms:Describe*",
                "ecr:Get*",
                "eks:Describe*",
                "eks:List*",
                "elasticloadbalancing:Describe*",
                "elasticmapreduce:Get*",
                "glue:Get*",
                "guardduty:Get*",
                "guardduty:List*",
                "iam:GenerateCredentialReport",
                "iam:Get*",
                "kafka:List*",
                "kms:Describe*",
                "kms:Get*",
                "kms:List*",
                "lightsail:Get*",
                "redshift:Describe*",
                "secretsmanager:List*",
                "securityhub:Describe*",
                "servicequotas:List*",
                "shield:Describe*",
                "SNS:Get*",
                "ssm:Describe*",
                "transfer:List*",
                "xray:Get*"
            ],
            "Resource": "*"
        }
    ]
}
```

## Sequence

Diagrams needed:

- Role/trust hierarchy [TD]
- CronJob job spawning and role assumption [sequence]

```mermaid
sequenceDiagram
    participant Orchestrator CronJob
    participant Orchestrator IAM role
    participant Runner Job

    CronJob ->> AWS STS: Assume "Inventory-Orchestrator" IAM role
    Alice ->> Bob: Hello Bob, how are you?
    Bob-->>John: How about you John?
    Bob--x Alice: I am good thanks!
    Bob-x John: I am good thanks!
    Note right of John: Bob thinks a long<br/>long time, so long<br/>that the text does<br/>not fit on a row.

    Bob-->Alice: Checking with John...
    Alice->John: Yes... John, how are you?
```

```mermaid
graph TD
    A["inventory-orchestrator" K8S CronJob] --> B("aws-inventory-orchestrator-sa" K8S SA)
    B --> C["Inventory-Orchestrator" AWS IAM role ]
    C -->|One| D[Laptop]
    C -->|Two| E[iPhone]
    C -->|Three| F[fa:fa-car Car]
```


## Development

*Work in progress.*

Create `./k8s/vars.env`:

```env
ORCHESTRATOR_ROLE_ARN=arn:aws:iam::$BILLING_ACCOUNT_ID:role/Inventory-Orchestrator
RUNNER_ROLE_ARN=arn:aws:iam::$SECURITY_ACCOUNT_ID:role/Inventory-Runner
CRON_SCHEDULE=* * * * 0
```

Run `skaffold dev`.

Trigger CronJob (`kubectl command`)

```
kubectl -n inventory create job aws-inventory-orchestrator-manual --from=cronjob/aws-inventory-orchestrator
```

```
aws_recon -v -r global,eu-west-1,eu-central-1 --s3-bucket raras-inventory:eu-west-1
aws_recon -v -r global,eu-west-1,eu-central-1 | grep "not authorized"
```

Example path:

s3://${BUCKET_NAME}/AWSRecon/${YEAR}/${MONTH}/${DAY}/${ACCOUNT_ID}_aws_recon_1649767394.json.gz
