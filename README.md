# SORACOM CLI Lambda Function

An AWS Lambda Function to execute [SORACOM CLI](https://github.com/soracom/soracom-cli).

## How to deploy

### Prerequisite

1. Installing Docker.
2. Installing [AWS SAM CLI](https://github.com/aws/aws-sam-cli).

### Step 1: Create your SORACOM AuthKey

Generate a new SORACOM AuthKey and save it to AWS Secrets Manager.

```sh
aws secretsmanager create-secret \
  --name "soracom-cli/profile/default" \
  --secret-string file://authkey.json
```

Contents of `authkey.json`:

```json
{
  "AUTH_KEY_ID": "keyId-xxxxxxxxxx",
  "AUTH_KEY": "secret-xxxxxxxxxx"
}
```

Copy ARN of the created secret. You will use ARN in step 3.

### Step 2: Build the function

```sh
sam build
```

### Step 3: Deploy the function to your AWS account

```sh
sam deploy --guided
```

Follow the on-screen prompts.

- Input ARN of your secret for `SoracomAuthKeySecretName` parameter.
- You can use default or custom values for other parameters.

## How to invoke

If you deployed your function with default function name (`soracom-cli`):

```json
aws lambda invoke --function-name soracom-cli \
  --payload file://payload.json
```

Contents of `payload.json`:

```json
{
  "command": "sims list"
}
```

For commands that consumes request body:

```json
aws lambda invoke --function-name soracom-cli \
  --payload file://payload-with-body.json
```

Contents of `payload-with-body.json`:

```json
{
  "command": "auth",
  "body": {
    "operatorId": "OPXXXXXXXXXX",
    "password": "p@$$w0rd",
    "userName": "SORACOMAPI"
  }
}
```
