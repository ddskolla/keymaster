
variable "env_label" {
  description = "The tag label of the environment km will be deployed into (e.g. btr-place)"
  type = string
}

variable "resource_tags" {
  description = "Map of tags to apply to all AWS resources"
  type = map(string)
  default = {}
}

variable "artifact_file" {
  description = "Local path to lambda deployment package. Conflicts with artifact_s3*"
  type = string
  default = null
}

variable "artifact_s3_bucket" {
  description = "S3 bucket with existing keymaster deployment artifact (lambda zip file)"
  type = string
  default = null
}

variable "artifact_s3_key" {
  description = "S3 key with existing keymaster deployment artifcat (lambda zip file)"
  type = string
  default = null
}

variable "client_account_arns" {
  description = "List of accounts with permission to invoke km issuing api"
  type = list(string)
  default = []
}

variable "target_role_arns" {
  description = "List of roles which km may issue credentials for"
  type = list(string)
  default = []
}

variable "lambda_function_name" {
  description = "Lambda function name to create"
  type = string
  default = ""
}

variable "lambda_role_arn" {
  description = "Set this to override the IAM role used by the km issuing lambda"
  type = string
  default = ""
}

variable "configuration" {
  description = "Keymaster configuration (environment variables)"
  type = map(string)
}

variable "reserved_concurrent_executions" {
  description = "Reserved executions for each keymaster lambda"
  type = number
  default = -1
}

variable "timeout" {
  description = "Lambda timeout"
  type = number
  default = 30
}

variable "config_bucket_enable" {
  description = "Create the config bucket"
  type = bool
  default = false
}

variable "config_file_upload_enable" {
  description = "Enable uploading a configuration file for km"
  type = bool
  default = false
}

variable "config_bucket_name" {
  description = "Name of bucket to store configuration file"
  type = string
  default = ""
}

variable "config_file_name" {
  description = "Name of local file to upload for km configuration"
  type = string
  default = ""
}
// TODO: kms key ARN for env var config
