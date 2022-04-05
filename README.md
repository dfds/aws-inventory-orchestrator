# inventory-orchestrator

## Development

*Work in progress.*

Create `./k8s/vars.env`:

```env
ROLE_ARN=<ARN of the AWS IAM role (typically 'Inventory-Orchestrator') for the workload to assume>
CRON_SCHEDULE=<Cron Schedule>
```

Suggested `CRON_SCHEDULE`:

- Prod: `* * * * 0`
- Dev: `* * * * *`

Run `skaffold dev`.