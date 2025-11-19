# aws-securityhub-integration-slack-go

AWS Lambda function that sends **AWS Security Hub v2** findings to Slack via EventBridge. **Security Hub v2** uses OCSF (Open Cybersecurity Schema Framework) format and centralizes findings from GuardDuty, Inspector, Macie, IAM Access Analyzer, and Security Hub CSPM.

> **Note:** This is for Security Hub v2 only. Not compatible with the original AWS Security Hub (now called Security Hub CSPM).

## Features

* **multi-service support** – handles findings from GuardDuty, Inspector, Macie, IAM Access Analyzer, and Security Hub CSPM
* **ocsf format** – native support for Security Hub v2 OCSF schema
* **eventbridge trigger** – findings invoke the Lambda function directly
* **rich messages** – displays source service, severity, category, region, account, resource details, and remediation links
* **config-driven** – all behavior controlled by environment variables
* **severity filtering** – EventBridge rules can filter by severity (Critical/High only)

---

## Deployment

### Prerequisites

* AWS account with **AWS Security Hub v2** enabled in at least one region
  * **Important:** This must be Security Hub v2, not the original Security Hub (now Security Hub CSPM)
  * Security Hub v2 uses OCSF format and has product ARNs like `arn:aws:securityhub:region::productv2/aws/guardduty`
* At least one integrated security service enabled (GuardDuty, Inspector, Macie, IAM Access Analyzer, or Security Hub CSPM)
* Slack App with a Bot Token (`chat:write` scope) installed in your workspace
* Go ≥ 1.24
* AWS CLI configured for the deployment account

### Steps

```bash
git clone https://github.com/cruxstack/aws-securityhub-integration-slack-go.git
cd aws-securityhub-integration-slack-go

# build static Linux binary for lambda
mkdir -p dist
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -C cmd/lambda -o ../../dist/bootstrap

# package
cd dist && zip deployment.zip bootstrap && cd ..
```

## Required Environment Variables

| name                        | example                                    | purpose                                                      |
| --------------------------- | ------------------------------------------ | ------------------------------------------------------------ |
| `APP_SLACK_TOKEN`           | `xoxb-…`                                   | slack bot token (store in secrets manager)                   |
| `APP_SLACK_CHANNEL`         | `C000XXXXXXX`                              | channel id to post findings                                  |

## Optional Environment Variables

| name                           | example                                    | purpose                                                      | default                           |
| ------------------------------ | ------------------------------------------ | ------------------------------------------------------------ | --------------------------------- |
| `APP_DEBUG_ENABLED`            | `true`                                     | verbose logging & event dump                                 | `false`                           |
| `APP_AWS_CONSOLE_URL`          | `https://console.aws.amazon.com`           | base AWS console URL                                         | `https://console.aws.amazon.com`  |
| `APP_AWS_ACCESS_PORTAL_URL`    | `https://myorg.awsapps.com/start`          | AWS access portal URL (for federated access)                 | _(none - direct console links)_   |
| `APP_AWS_ACCESS_ROLE_NAME`     | `SecurityAuditor`                          | IAM role name for access portal                              | _(none - direct console links)_   |
| `APP_AWS_SECURITYHUBV2_REGION` | `us-east-1`                                | AWS region for centralized SecurityHub v2 if applicable      | _(none - direct console links)_   |

## Create Lambda Function

1. **IAM role**
   * `AWSLambdaBasicExecutionRole` managed policy
   * no additional AWS API permissions are required
2. **Lambda config**
   * Runtime: `al2023provided.al2023` (provided.al2 also works)
   * Handler: `bootstrap`
   * Architecture: `x86_64` or `arm64`
   * Upload `deployment.zip`
   * Set environment variables above
3. **EventBridge rule**
    ```json
    {
      "source": ["aws.securityhub"],
      "detail-type": ["Findings Imported V2"]
    }
   ```
   Optional: Filter by severity (recommended for high-volume environments):
    ```json
    {
      "source": ["aws.securityhub"],
      "detail-type": ["Findings Imported V2"],
      "detail": {
        "findings": {
          "severity": ["Critical", "High"]
        }
      }
    }
   ```
   Or filter by specific source services:
    ```json
    {
      "source": ["aws.securityhub"],
      "detail-type": ["Findings Imported V2"],
      "detail": {
        "findings": {
          "metadata": {
            "product": {
              "name": ["GuardDuty", "Inspector"]
            }
          }
        }
      }
    }
   ```
   Target: the Lambda function.
4. **Slack App**
   * Add `chat:write` and `chat:write.public`
   * Custom bot avatar: upload AWS Security Hub logo in the Slack App *App Icon*
     section.


## Local Development

### Test with Samples

```bash
cp .env.example .env # edit the values
go run -C cmd/sample .
```

The sample runner reads OCSF-formatted Security Hub v2 findings from `fixtures/samples.json`, wraps them in EventBridge events, and posts to Slack exactly as the live Lambda would.

