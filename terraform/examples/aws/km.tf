
resource "aws_iam_role" "keymaster" {
  name = "keymaster-role"
  assume_role_policy = data.aws_iam_policy_document.keymaster-assumerole.json
}

data "aws_iam_policy_document" "keymaster-assumerole" {
  statement {
    principals {
      type = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
    actions = ["sts:AssumeRole"]
  }
}

data "aws_iam_policy_document" "keymaster-policy" {
  statement {
    actions = [
      "s3:GetObject",
    ]
    resources = [
//      "${aws_s3_bucket.keymaster-secrets.arn}",
//      "${aws_s3_bucket.keymaster-secrets.arn}/*",
//      "${data.aws_s3_bucket.state-bucket.arn}",
//      "${data.aws_s3_bucket.state-bucket.arn}/*",
    ]
  }
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]
    resources = [
      "arn:aws:logs:*:*:*"
    ]
  }
}

locals {
  keymaster_config = {

  }
}

module "keymaster_api" {
  source = "../../aws/api_lambda"
  artifact_file = "../../../build/keymaster-api.zip"
  lambda_role_arn = aws_iam_role.keymaster.arn
  configuration = local.keymaster_config
}
