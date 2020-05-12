data "aws_iam_policy_document" "ip_oracle_verify_lambda_role_sts" {
  statement {
  	actions = ["sts:AssumeRole"]
    effect  = "Allow"
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "ip_oracle_verify_lambda_role_verify" {
  statement {
    actions   = ["kms:Verify"]
    resources = ["${var.ip_oracle_key_arn}"]
    effect 	  = "Allow"
  }
}

resource "aws_iam_role" "ip_oracle_verify_role" {
  name               = "${var.resource_name_prefix}-ip-oracle-verify-role"
  description        = "Keymaster IP oracle verify lambda role"
  assume_role_policy = "${data.aws_iam_policy_document.ip_oracle_verify_lambda_role_sts.json}"
}

resource "aws_iam_policy" "ip_oracle_verify_lambda_role_verify_policy" {
  name        = "${var.resource_name_prefix}-ip-oracle-verify-policy"
  description = "Lambda role policy for IP oracle verify"
  policy      = "${data.aws_iam_policy_document.ip_oracle_verify_lambda_role_verify.json}"
}

resource "aws_iam_role_policy_attachment" "ip_oracle_verify" {
  role       = "${aws_iam_role.ip_oracle_verify_role.name}"
  policy_arn = "${aws_iam_policy.ip_oracle_verify_lambda_role_verify_policy.arn}"
}

resource "aws_lambda_function" "ip_oracle_verify" {
  filename         = "${var.ip_oracle_verify_zip}"
  function_name    = "${var.resource_name_prefix}-ip_oracle_verify"
  role             = "${aws_iam_role.ip_oracle_verify_role.arn}"
  handler          = "${var.lambda_verify_handler}"
  source_code_hash = "${filebase64sha256("${var.ip_oracle_verify_zip}")}"
  runtime          = "go1.x"
  timeout          = var.timeout
}

resource "aws_api_gateway_rest_api" "ip_oracle_verify" {
  name        = "${var.resource_name_prefix}-ip-oracle-verify"
  description = "Keymaster API gateway"

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_resource" "ip_oracle_verify" {
  rest_api_id = "${aws_api_gateway_rest_api.ip_oracle_verify.id}"
  parent_id   = "${aws_api_gateway_rest_api.ip_oracle_verify.root_resource_id}"
  path_part   = "verify"
}

resource "aws_api_gateway_method" "ip_oracle_verify" {
  rest_api_id   = "${aws_api_gateway_rest_api.ip_oracle_verify.id}"
  resource_id   = "${aws_api_gateway_resource.ip_oracle_verify.id}"
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "ip_oracle_verify" {
  rest_api_id = "${aws_api_gateway_rest_api.ip_oracle_verify.id}"
  resource_id = "${aws_api_gateway_method.ip_oracle_verify.resource_id}"
  http_method = "${aws_api_gateway_method.ip_oracle_verify.http_method}"

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "${aws_lambda_function.ip_oracle_verify.invoke_arn}"
}

resource "aws_api_gateway_deployment" "ip_oracle_verify" {
  depends_on  = ["aws_api_gateway_integration.ip_oracle_verify"]
  rest_api_id = "${aws_api_gateway_rest_api.ip_oracle_verify.id}"
  stage_name  = "live"
}

resource "aws_lambda_permission" "apigw_verify" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.ip_oracle_verify.arn}"
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_deployment.ip_oracle_verify.execution_arn}/*"
}