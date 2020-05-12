Usage:

```hcl
module "issuing_lambda" {
  source = "github.com/bsycorp/keymaster/terraform/aws/issuing_lambda"

  # The environment label will be added to all named resources
  env_label   = "myproject-npe"

  # Keymaster configuration file
  configuration = {
     CONFIG: "s3://km-tools-bls-01/km.yaml"
  }

  # List of target roles that the lambda may issue creds for
  target_role_arns = [
   "arn:aws:iam::218296299700:role/test_env_admin"
  ]

  # List of client accounts that may invoke issuing lambda
  client_account_arns = [
   "arn:aws:iam::062921715666:root",   # myproj-dev-01
  ]

  # Enable auto-creation of the configuration bucket
  config_bucket_enable = true
  config_file_upload_enable = true
  config_file_name = "${path.module}/test_api_config.yaml"

  resource_tags = {
    Name         = "baz"
    Created-By   = "you@your.com"
  }
}
```

## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| aws | n/a |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| artifact\_file | Local path to lambda deployment package. Conflicts with artifact\_s3\* | `string` | `null` | no |
| artifact\_s3\_bucket | S3 bucket with existing keymaster deployment artifact (lambda zip file) | `string` | `null` | no |
| artifact\_s3\_key | S3 key with existing keymaster deployment artifcat (lambda zip file) | `string` | `null` | no |
| client\_account\_arns | List of accounts with permission to invoke km issuing api | `list(string)` | `[]` | no |
| config\_bucket\_enable | Create the config bucket | `bool` | `false` | no |
| config\_bucket\_name | Name of bucket to store configuration file | `string` | `""` | no |
| config\_file\_name | Name of local file to upload for km configuration | `string` | `""` | no |
| config\_file\_upload\_enable | Enable uploading a configuration file for km | `bool` | `false` | no |
| configuration | Keymaster configuration (environment variables) | `map(string)` | n/a | yes |
| env\_label | The tag label of the environment km will be deployed into (e.g. btr-place) | `string` | n/a | yes |
| lambda\_function\_name | Lambda function name to create | `string` | `""` | no |
| lambda\_role\_arn | Set this to override the IAM role used by the km issuing lambda | `string` | `""` | no |
| reserved\_concurrent\_executions | Reserved executions for each keymaster lambda | `number` | `-1` | no |
| resource\_tags | Map of tags to apply to all AWS resources | `map(string)` | `{}` | no |
| target\_role\_arns | List of roles which km may issue credentials for | `list(string)` | `[]` | no |
| timeout | Lambda timeout | `number` | `30` | no |

## Outputs

| Name | Description |
|------|-------------|
| configuration\_bucket\_name | The name of the km configuration bucket |
| issuing\_lambda\_arn | The ARN of the keymaster issuing lambda |

