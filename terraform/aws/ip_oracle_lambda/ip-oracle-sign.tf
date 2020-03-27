data "aws_iam_policy_document" "ip_oracle_sign_lambda_role_sts" {
  statement {
  	actions = ["sts:AssumeRole"]
    effect  = "Allow"
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "ip_oracle_sign_lambda_role_sign" {
  statement {
    actions   = ["kms:Sign"]
    resources = ["${aws_kms_key.ip_oracle_key.arn}"]
    effect 	  = "Allow"
  }
}

resource "aws_iam_role" "ip_oracle_sign_role" {
  name               = "${var.resource_name_prefix}-ip-oracle-sign-role"
  description        = "Keymaster IP oracle sign lambda role"
  assume_role_policy = "${data.aws_iam_policy_document.ip_oracle_sign_lambda_role_sts.json}"
}

resource "aws_iam_policy" "ip_oracle_sign_lambda_role_sign_policy" {
  name               = "${var.resource_name_prefix}-ip-oracle-sign-policy"
  description        = "Lambda role policy for IP oracle sign"
  policy             = "${data.aws_iam_policy_document.ip_oracle_sign_lambda_role_sign.json}"
}

resource "aws_iam_role_policy_attachment" "ip_oracle_sign" {
  role       = "${aws_iam_role.ip_oracle_sign_role.name}"
  policy_arn = "${aws_iam_policy.ip_oracle_sign_lambda_role_sign_policy.arn}"
}

resource "aws_kms_key" "ip_oracle_key" {
	description					      = "KMS key for KeyMaster IP Oracle sign and verify"
	key_usage					        = "SIGN_VERIFY"
	customer_master_key_spec	= "RSA_2048"
}

resource "aws_lambda_function" "ip_oracle_sign" {
  filename         = "${var.ip_oracle_sign_zip}"
  function_name    = "${var.resource_name_prefix}-ip_oracle_sign"
  role             = "${aws_iam_role.ip_oracle_sign_role.arn}"
  handler          = "${var.lambda_sign_handler}"
  source_code_hash = "${filebase64sha256("${var.ip_oracle_sign_zip}")}"
  runtime          = "go1.x"
  timeout          = var.timeout
}

resource "aws_api_gateway_rest_api" "ip_oracle_sign" {
  name        = "${var.resource_name_prefix}-ip-oracle-sign"
  description = "Keymaster API gateway"

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

resource "aws_api_gateway_resource" "ip_oracle_sign" {
  rest_api_id = "${aws_api_gateway_rest_api.ip_oracle_sign.id}"
  parent_id   = "${aws_api_gateway_rest_api.ip_oracle_sign.root_resource_id}"
  path_part   = "sign"
}

resource "aws_api_gateway_method" "ip_oracle_sign" {
  rest_api_id   = "${aws_api_gateway_rest_api.ip_oracle_sign.id}"
  resource_id   = "${aws_api_gateway_resource.ip_oracle_sign.id}"
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "ip_oracle_sign" {
  rest_api_id = "${aws_api_gateway_rest_api.ip_oracle_sign.id}"
  resource_id = "${aws_api_gateway_method.ip_oracle_sign.resource_id}"
  http_method = "${aws_api_gateway_method.ip_oracle_sign.http_method}"

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "${aws_lambda_function.ip_oracle_sign.invoke_arn}"
}

resource "aws_api_gateway_deployment" "ip_oracle_sign" {
  depends_on  = ["aws_api_gateway_integration.ip_oracle_sign"]
  rest_api_id = "${aws_api_gateway_rest_api.ip_oracle_sign.id}"
  stage_name  = "live"
}

resource "aws_lambda_permission" "apigw_sign" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.ip_oracle_sign.arn}"
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_api_gateway_deployment.ip_oracle_sign.execution_arn}/*"
}