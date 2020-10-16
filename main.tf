# This is required to get the AWS region via ${data.config.default}.
data "aws_region" "current" {}



provider "aws" {
  region = "us-east-1"
}


# Define a Lambda function.
#
# The handler is the name of the executable for go1.x runtime.
resource "aws_lambda_function" "portfolioGo" {
  function_name    = "portfolioGo"
  filename         = "portfolioGo.zip"
  handler          = "main"
  source_code_hash = base64sha256(("portfolioGo"))
  role             = aws_iam_role.portfolioGo.arn
  runtime          = "go1.x"
  memory_size      = 128
  timeout          = 3
  tags = {
    "key" = "value"
  }

}

# A Lambda function may access to other AWS resources such as S3 bucket. So an
# IAM role needs to be defined.

data "aws_iam_policy_document" "portfolioGo-assume-role-policy" {
  statement {
    actions = ["sts:AssumeRole"]
    # resources = [aws_dynamodb_table.visitorcountGo.arn,
    

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "portfolioGo" {
  name               = "portfolio_role"
  path               = "/system/"
  assume_role_policy = data.aws_iam_policy_document.portfolioGo-assume-role-policy.json
}

//attaches policy to the role giving full dynamo access
// TODO should be least privilege 
resource "aws_iam_role_policy_attachment" "portfolioGo_policy" {
  role       = "${aws_iam_role.portfolioGo.name}"
  policy_arn = "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess"
}


# Allow API gateway to invoke the hello Lambda function.
resource "aws_lambda_permission" "portfolioGo" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.portfolioGo.arn
  principal     = "apigateway.amazonaws.com"
}

resource "aws_dynamodb_table" "VisitorcountGo" {
  name           = "VisitorCountgo"
  billing_mode   = "PROVISIONED"
  read_capacity  = 20
  write_capacity = 20
  hash_key       = "key"

  attribute {
    name = "key"
    type = "S"
  }
}
//API Gateway creating the API
resource "aws_api_gateway_resource" "portfolioGo" {
  rest_api_id = aws_api_gateway_rest_api.portfolioGo.id
  parent_id   = aws_api_gateway_rest_api.portfolioGo.root_resource_id
  path_part   = "portfoliogo"
}

resource "aws_api_gateway_rest_api" "portfolioGo" {
  name = "portfolioGo"
}

#           GET
resource "aws_api_gateway_method" "portfolioGo" {
  rest_api_id   = aws_api_gateway_rest_api.portfolioGo.id
  resource_id   = aws_api_gateway_resource.portfolioGo.id
  http_method   = "GET"
  authorization = "NONE"
}

#              POST
# API Gateway ------> Lambda
resource "aws_api_gateway_integration" "portfolioGo" {
  rest_api_id             = aws_api_gateway_rest_api.portfolioGo.id
  resource_id             = aws_api_gateway_resource.portfolioGo.id
  http_method             = aws_api_gateway_method.portfolioGo.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "arn:aws:apigateway:${data.aws_region.current.name}:lambda:path/2015-03-31/functions/${aws_lambda_function.portfolioGo.arn}/invocations"
}

# This resource defines the URL of the API Gateway.
resource "aws_api_gateway_deployment" "portfolioGo_v1" {
  depends_on = [
    aws_api_gateway_integration.portfolioGo
  ]
  rest_api_id = aws_api_gateway_rest_api.portfolioGo.id
  stage_name  = "v1"
}

# Set the generated URL as an output. Run `terraform output url` to get this.
output "url" {
  value = "${aws_api_gateway_deployment.portfolioGo_v1.invoke_url}${aws_api_gateway_resource.portfolioGo.path}"
}

//Frontend
resource "aws_s3_bucket" "portfolioGo" {
  bucket = "chadiamond.iogo"
  acl = "private"
  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT","POST"]
    allowed_origins = [""]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
  versioning {
    enabled = true
  }
}
data archive_file lambda_zip {
type        = "zip"
  source_file = "${path.module}/portfolio.go"
  output_path = "${path.module}/files/portfoliogo.zip"
}