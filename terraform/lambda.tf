resource "aws_iam_role" "iam_for_lambda" {
  name = "iam_for_lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
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
