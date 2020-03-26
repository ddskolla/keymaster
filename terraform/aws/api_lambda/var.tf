
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

variable "lambda_function_name" {
  description = "Lambda function name to create"
  type = string
  default = "km2"
}

variable "lambda_role_arn" {
  description = "Role for keymaster lambda"
  type = string
  default = null
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

// TODO: kms key ARN for env var config
