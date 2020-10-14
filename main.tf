# This is required to get the AWS region via ${data.config.default}.
data "aws_region" "current" {}



provider "aws" {
  region = "us-east-1"
}


# Define a Lambda function.
#
# The handler is the name of the executable for go1.x runtime.
resource "aws_lambda_function" "goportfolio" {
  function_name    = "gocounting"
  filename         = "gocounting.zip"
  handler          = "gocounting"
  source_code_hash = "${base64sha256(("gocounting"))}"
  role             = "${aws_iam_role.gocounting.arn}"
  runtime          = "go1.x"
  memory_size      = 128
  timeout          = 1

}

# A Lambda function may access to other AWS resources such as S3 bucket. So an
# IAM role needs to be defined.
resource "aws_iam_role" "hello" {
  name               = "hello"
  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": {
    "Action": "sts:AssumeRole",
    "Principal": {
      "Service": "lambda.amazonaws.com"
    },
    "Effect": "Allow"
  }
}
POLICY
}

# Allow API gateway to invoke the hello Lambda function.
resource "aws_lambda_permission" "hello" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.gocounting.arn}"
  principal     = "apigateway.amazonaws.com"
}

resource "aws_dynamodb_table" "basic-dynamodb-table" {
  name           = "movies"
  billing_mode   = "PROVISIONED"
  read_capacity  = 20
  write_capacity = 20
  hash_key       = "UserId"
  range_key      = "MovieTitle"

  attribute {
    name = "Movie"
    type = "S"
  }

  attribute {
    name = "Year"
    type = "N"
  }

  ttl {
    attribute_name = "TimeToExist"
    enabled        = false
  }
}

# A Lambda function is not a usual public REST API. We need to use AWS API
# Gateway to map a Lambda function to an HTTP endpoint.
resource "aws_api_gateway_resource" "gocounting" {
  rest_api_id = "${aws_api_gateway_rest_api.gocounting.id}"
  parent_id   = "${aws_api_gateway_rest_api.gocounting.root_resource_id}"
  path_part   = "gocounting"
}

resource "aws_api_gateway_rest_api" "gocounting" {
  name = "gocounting"
}

#           GET
# Internet -----> API Gateway
resource "aws_api_gateway_method" "gocounting" {
  rest_api_id   = "${aws_api_gateway_rest_api.gocounting.id}"
  resource_id   = "${aws_api_gateway_resource.gocounting.id}"
  http_method   = "GET"
  authorization = "NONE"
}

#              POST
# API Gateway ------> Lambda
# For Lambda the method is always POST and the type is always AWS_PROXY.
#
# The date 2015-03-31 in the URI is just the version of AWS Lambda.
resource "aws_api_gateway_integration" "gocounting" {
  rest_api_id             = "${aws_api_gateway_rest_api.gocounting.id}"
  resource_id             = "${aws_api_gateway_resource.gocounting.id}"
  http_method             = "${aws_api_gateway_method.gocounting.http_method}"
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "arn:aws:apigateway:${data.aws_region.current.name}:lambda:path/2015-03-31/functions/${aws_lambda_function.gocounting.arn}/invocations"
}

# This resource defines the URL of the API Gateway.
resource "aws_api_gateway_deployment" "gocounting_v1" {
  depends_on = [
    aws_api_gateway_integration.gocounting
  ]
  rest_api_id = "${aws_api_gateway_rest_api.gocounting.id}"
  stage_name  = "v1"
}

# Set the generated URL as an output. Run `terraform output url` to get this.
output "url" {
  value = "${aws_api_gateway_deployment.gocounting_v1.invoke_url}${aws_api_gateway_resource.hello.path}"
}