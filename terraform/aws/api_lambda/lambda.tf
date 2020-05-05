
locals {
  timeout              = 30
  handler              = "api-linux-x64"
  runtime              = "go1.x"
  config_bucket_name   = var.config_bucket_name == "" ? "km-${var.env_label}" : var.config_bucket_name
  lambda_function_name = var.lambda_function_name == "" ? "km-${var.env_label}" : var.lambda_function_name
  lambda_role_arn      = var.lambda_role_arn == "" ? aws_iam_role.km[0].arn : var.lambda_role_arn
}

resource "aws_lambda_function" "km" {
  function_name = local.lambda_function_name
  handler       = local.handler
  runtime       = local.runtime
  filename      = var.artifact_file
  s3_bucket     = var.artifact_s3_bucket
  s3_key        = var.artifact_s3_key

  role    = local.lambda_role_arn
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

resource "null_resource" "lambda_changed" {
  triggers = {
    lambda_modified = aws_lambda_function.km.last_modified
  }
}

resource "aws_lambda_permission" "allow_invoke" {
  depends_on = [null_resource.lambda_changed]
  count         = length(var.client_account_arns)
  statement_id  = "AllowClientExecution-${count.index}"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.km.function_name
  principal     = var.client_account_arns[count.index]
}

resource "aws_s3_bucket" "km_config" {
  count         = var.config_bucket_enable ? 1 : 0
  bucket        = local.config_bucket_name
  region        = "ap-southeast-2" # TODO
  force_destroy = true
  tags          = merge({}, var.resource_tags)
}

resource "aws_s3_bucket_object" "km_config_file" {
  count  = var.config_file_upload_enable ? 1 : 0
  bucket = aws_s3_bucket.km_config[0].bucket
  key    = "km.yaml"
  source = var.config_file_name
}

