
locals {
  timeout = 30
  handler = "api-linux-x64"
  runtime = "go1.x"
  lambda_config = {
    //      keymaster_version    = var.keymaster_version
    //      kube_env_name        = var.instance_domain
    //      kube_userca_bucket   = aws_s3_bucket.keymaster-secrets.bucket
    //      kube_sshca_bucket    = aws_s3_bucket.keymaster-secrets.bucket
    //      kube_sshca_key       = aws_s3_bucket_object.sshca-key.key
    //      saml_idp_certificate = file("${path.module}/${var.keymaster_saml_idp_certificate_file}")
    //      valid_for_seconds    = var.keymaster_valid_for_seconds
    //      username_pattern     = var.keymaster_username_pattern
    //      service_config       = var.keymaster_service_config
  }
  handler_functions = ["config", "direct_saml_auth", "direct_oidc_auth", "workflow_start", "workflow_auth"]
}

resource "aws_lambda_function" "km" {
  function_name = var.lambda_function_name
  handler       = local.handler
  runtime       = local.runtime
  filename      = var.artifact_file
  s3_bucket     = var.artifact_s3_bucket
  s3_key        = var.artifact_s3_key

  // TODO: conditional
  //role          = var.lambda_role_arn
  role    = aws_iam_role.km.arn
  timeout = var.timeout

  dynamic "environment" {
    for_each = var.configuration[*]
    content {
      variables = environment.value
    }
  }

  reserved_concurrent_executions = var.reserved_concurrent_executions
  tags                           = merge({}, var.resource_tags)
}
