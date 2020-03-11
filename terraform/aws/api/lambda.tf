
locals {
  // Download and upload from a url
  // artifact_s3_bucket
  // artifact_s3_key
}

resource "aws_lambda_function" "keymaster" {
  function_name = var.lambda_function_name
  handler       = "api-linux-x64"
  runtime       = "go1.x"
  s3_bucket     = var.artifact_s3_bucket
  s3_key        = var.artifact_s3_key
  role          = var.lambda_role_arn
  timeout       = 30
  environment {
    variables = {
//      keymaster_version    = var.keymaster_version
//      kube_env_name        = var.instance_domain
//      kube_userca_bucket   = aws_s3_bucket.keymaster-secrets.bucket
//      kube_sshca_bucket    = aws_s3_bucket.keymaster-secrets.bucket
//      kube_sshca_key       = aws_s3_bucket_object.sshca-key.key
//      saml_idp_certificate = file("${path.module}/${var.keymaster_saml_idp_certificate_file}")
//      ip_whitelist         = join(",", var.keymaster_ip_whitelist)
//      valid_for_seconds    = var.keymaster_valid_for_seconds
//      username_pattern     = var.keymaster_username_pattern
//      service_config       = var.keymaster_service_config
    }
  }
  tags          = merge({}, var.resource_tags)
}
