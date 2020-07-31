terraform {
  backend "s3" {
    bucket = "countries-ec1d8a43-4a61-4883-b299-51090a0f7a6e"
    key    = "config/terraform.tfstate"
    region = "eu-west-2"
  }
}

provider "aws" {
  profile = "default"
  region  = "eu-west-2"
}

variable "fcup_email" {
  type = string
}

variable "fcup_name" {
  type = string
}

variable "s3_bucket" {
  type = string
}

provider "archive" {}

data "archive_file" "zip" {
  type        = "zip"
  source_file = "bin/freedom-countries"
  output_path = "freedom-countries.zip"
}

resource "aws_s3_bucket_object" "email_template" {
  bucket = var.s3_bucket
  key    = "email_template.html"
  source = "templates/email.html"

  etag = filemd5("templates/email.html")
}

data "aws_iam_policy_document" "policy" {
  statement {
    sid    = ""
    effect = "Allow"

    principals {
      identifiers = ["lambda.amazonaws.com"]
      type        = "Service"
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "iam_for_lambda" {
  name               = "iam_for_lambda"
  assume_role_policy = data.aws_iam_policy_document.policy.json
}

# This is to optionally manage the CloudWatch Log Group for the Lambda Function.
# If skipping this resource configuration, also add "logs:CreateLogGroup" to the IAM policy below.
# resource "aws_cloudwatch_log_group" "example" {
# name              = "/aws/lambda/${var.lambda_function_name}"
# retention_in_days = 14
# }

# See also the following AWS managed policy: AWSLambdaBasicExecutionRole
resource "aws_iam_policy" "lambda_s3_full" {
  name        = "lambda_s3"
  path        = "/"
  description = "IAM policy for S3 from a lambda"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "s3:GetObject",
        "s3:PutObject"
      ],
      "Resource": "arn:aws:s3:::${var.s3_bucket}/*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda_s3" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = aws_iam_policy.lambda_s3_full.arn
}

resource "aws_iam_policy" "lambda_ses" {
  name        = "lambda_ses"
  path        = "/"
  description = "IAM policy for SES from a lambda"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "ses:SendEmail",
        "ses:SendRawEmail"
      ],
      "Resource": "arn:aws:ses:*:*:*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda_ses" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = aws_iam_policy.lambda_ses.arn
}

resource "aws_iam_policy" "lambda_logging" {
  name        = "lambda_logging"
  path        = "/"
  description = "IAM policy for logging from a lambda"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda_logs" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = aws_iam_policy.lambda_logging.arn
}

resource "aws_lambda_function" "update_freedom_countries" {
  function_name = "update-freedom-countries"

  filename         = data.archive_file.zip.output_path
  source_code_hash = data.archive_file.zip.output_base64sha256

  role    = aws_iam_role.iam_for_lambda.arn
  handler = "freedom-countries"
  runtime = "go1.x"

  environment {
    variables = {
      S3_BUCKET : var.s3_bucket,
      FCUP_EMAIL : var.fcup_email,
      FCUP_NAME : var.fcup_name
    }
  }
  depends_on = [aws_iam_role_policy_attachment.lambda_logs, aws_iam_role_policy_attachment.lambda_s3]
}

resource "aws_lambda_permission" "allow_cloudwatch" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.update_freedom_countries.function_name
  principal     = "events.amazonaws.com"
}

resource "aws_cloudwatch_event_rule" "watcher_freedom_countries" {
  name                = "watch_freedom_countries_website"
  schedule_expression = "rate(6 minutes)"
  # schedule_expression = "rate(5 days)"
}

resource "aws_cloudwatch_event_target" "watcher" {
  rule = aws_cloudwatch_event_rule.watcher_freedom_countries.name
  arn  = aws_lambda_function.update_freedom_countries.arn
}
