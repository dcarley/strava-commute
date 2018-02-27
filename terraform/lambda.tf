resource "aws_iam_role" "iam_for_lambda" {
  name               = "iam_for_lambda"
  assume_role_policy = "${data.aws_iam_policy_document.assume_role.json}"
}

data "aws_iam_policy_document" "assume_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role_policy" "cloudwatch_logs" {
  name   = "cloudwatch_logs"
  role   = "${aws_iam_role.iam_for_lambda.name}"
  policy = "${data.aws_iam_policy_document.cloudwatch_logs.json}"
}

data "aws_iam_policy_document" "cloudwatch_logs" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]

    resources = [
      "*",
    ]
  }
}

resource "aws_lambda_function" "lambda" {
  filename         = "../dist/${var.name}.zip"
  function_name    = "${var.name}"
  role             = "${aws_iam_role.iam_for_lambda.arn}"
  source_code_hash = "${base64sha256(file("../dist/${var.name}.zip"))}"
  runtime          = "go1.x"
  handler          = "${var.name}"

  environment {
    variables = {
      STRAVA_API_TOKEN = "${var.strava_api_token}"
    }
  }
}
