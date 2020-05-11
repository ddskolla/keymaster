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
| invoke\_arn | n/a |

