
output "issuing_lambda_arn" {
  value = aws_lambda_function.km.arn
  description = "The ARN of the keymaster issuing lambda."
}

output "configuration_bucket_name" {
  value = var.config_bucket_enable ? aws_s3_bucket.km_config[0].bucket : ""
  description = "The name of the km configuration bucket. Will be empty if not configured."
}
