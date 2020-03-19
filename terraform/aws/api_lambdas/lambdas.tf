
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
  tracing_config = {

  }
}

// Using for_each would be nicer but copy/paste works across a wide
// range of tf versions

resource "aws_lambda_function" "km_config" {
  function_name = "${var.lambda_function_name}_config"
  handler       = local.handler
  runtime       = local.runtime
  s3_bucket     = var.artifact_s3_bucket
  s3_key        = var.artifact_s3_key
  role          = var.lambda_role_arn
  timeout       = local.timeout
  environment {
    variables = merge(local.lambda_config, {
      "_HANDLER" = "config"
    })
  }
  tags          = merge({}, var.resource_tags)
  reserved_concurrent_executions = var.reserved_concurrent_executions
}

resource "aws_lambda_function" "km_direct_saml_auth" {
  function_name = "${var.lambda_function_name}_direct_saml_auth"
  handler       = local.handler
  runtime       = local.runtime
  s3_bucket     = var.artifact_s3_bucket
  s3_key        = var.artifact_s3_key
  role          = var.lambda_role_arn
  timeout       = local.timeout
  environment {
    variables = merge(local.lambda_config, {
      "_HANDLER" = "direct_saml_auth"
    })
  }
  tags          = merge({}, var.resource_tags)
  reserved_concurrent_executions = var.reserved_concurrent_executions
}

resource "aws_lambda_function" "km_direct_oidc_auth" {
  function_name = "${var.lambda_function_name}_direct_oidc_auth"
  handler       = local.handler
  runtime       = local.runtime
  s3_bucket     = var.artifact_s3_bucket
  s3_key        = var.artifact_s3_key
  role          = var.lambda_role_arn
  timeout       = local.timeout
  environment {
    variables = merge(local.lambda_config, {
      "_HANDLER" = "direct_oidc_auth"
    })
  }
  tags          = merge({}, var.resource_tags)
  reserved_concurrent_executions = var.reserved_concurrent_executions
}

resource "aws_lambda_function" "km_workflow_start" {
  function_name = "${var.lambda_function_name}_workflow_start"
  handler       = local.handler
  runtime       = local.runtime
  s3_bucket     = var.artifact_s3_bucket
  s3_key        = var.artifact_s3_key
  role          = var.lambda_role_arn
  timeout       = local.timeout
  environment {
    variables = merge(local.lambda_config, {
      "_HANDLER" = "workflow_start"
    })
  }
  tags          = merge({}, var.resource_tags)
  reserved_concurrent_executions = var.reserved_concurrent_executions
}

resource "aws_lambda_function" "km_workflow_auth" {
  function_name = "${var.lambda_function_name}_workflow_auth"
  handler       = local.handler
  runtime       = local.runtime
  s3_bucket     = var.artifact_s3_bucket
  s3_key        = var.artifact_s3_key
  role          = var.lambda_role_arn
  timeout       = local.timeout
  environment {
    variables = merge(local.lambda_config, {
      "_HANDLER" = "workflow_auth"
    })
  }
  tags          = merge({}, var.resource_tags)
  reserved_concurrent_executions = var.reserved_concurrent_executions
}
