
variable "resource_tags" {
  description = "Map of tags to apply to all AWS resources"
  type = map(string)
}

variable "artifact_s3_bucket" {
  description = "S3 bucket with existing keymaster deployment artifact (lambda zip file)"
  type = string
}

variable "artifact_s3_key" {
  description = "S3 key with existing keymaster deployment artifcat (lambda zip file)"
  type = string
}

variable "lambda_function_name" {
  description = "Lambda function name to create"
  type = string
  default = "km2"
}

variable "lambda_role_arn" {
  description = "Role for keymaster lambda"
  type = string
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

// TODO: kms key ARN for env var config
