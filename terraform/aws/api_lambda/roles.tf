
data "aws_iam_policy_document" "lambda_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
    effect 	  = "Allow"
  }
}

data "aws_iam_policy_document" "km" {
  statement {
    actions = ["s3:GetObject"]
    resources = [
      "arn:aws:s3:::notreal-keymaster-secrets/*",
      "arn:aws:s3:::notreal-keymaster-secrets"
    ]
    effect = "Allow"
  }
  statement {
    actions = [
      "logs:PutLogEvents",
      "logs:CreateLogStream",
      "logs:CreateLogGroup"
    ]
    resources = [ "arn:aws:logs:*:*:*" ]
    effect = "Allow"
  }
}

resource "aws_iam_role" "km" {
  // TODO: name vars
  name = "km"
  description = "keymaster lambda role"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role.json
}

resource "aws_iam_policy" "km" {
  // TODO: name vars
  name = "km"
  description = "keymaster iam policy"
  policy = data.aws_iam_policy_document.km.json
}


resource "aws_iam_role_policy_attachment" "km" {
  role       = aws_iam_role.km.name
  policy_arn = aws_iam_policy.km.arn
}
