
data "aws_iam_policy_document" "lambda_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
    effect = "Allow"
  }
}

data "aws_iam_policy_document" "km" {
  statement {
    // TODO: read-only access
    actions = ["s3:*"]
    resources = [
      // TODO: we need to maintain a local list here...
      "${aws_s3_bucket.km_config[0].arn}",
      "${aws_s3_bucket.km_config[0].arn}/*",
    ]
    effect = "Allow"
  }
  statement {
    actions = [
      "logs:PutLogEvents",
      "logs:CreateLogStream",
      "logs:CreateLogGroup"
    ]
    resources = ["arn:aws:logs:*:*:*"]
    effect    = "Allow"
  }
  statement {
    actions = [
      "sts:AssumeRole",
    ]
    effect    = "Allow"
    resources = var.target_role_arns
  }
}

resource "aws_iam_role" "km" {
  count              = var.lambda_role_arn == "" ? 1 : 0
  name               = "km-${var.env_label}"
  description        = "keymaster issuing lambda role"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role.json
  tags               = merge({}, var.resource_tags)
}

resource "aws_iam_policy" "km" {
  name        = "km-${var.env_label}"
  description = "keymaster iam policy"
  policy      = data.aws_iam_policy_document.km.json
}

resource "aws_iam_role_policy_attachment" "km" {
  count      = var.lambda_role_arn == "" ? 1 : 0
  role       = aws_iam_role.km[0].name
  policy_arn = aws_iam_policy.km.arn
}
