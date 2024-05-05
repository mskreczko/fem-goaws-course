data "aws_iam_policy_document" "allow_dynamodb_table_operations" {
    statement {
        effect = "Allow"
        actions = [
            "dynamodb:PutItem",
            "dynamodb:GetItem",
            "dynamodb:DeleteItem",
        ]

        resources = [
            aws_dynamodb_table.users.arn,
        ]
    }
}

resource "aws_iam_policy" "dynamodb_lambda_policy" {
    name = "DynamoDBPolicy"
    description = "Policy for lambda to operate on dynamodb table"
    policy = data.aws_iam_policy_document.allow_dynamodb_table_operations.json
}

resource "aws_iam_role_policy_attachment" "lambda_dynamodb_policy_attachment" {
    role = aws_iam_role.lambda.id
    policy_arn = aws_iam_policy.dynamodb_lambda_policy.arn
}

data "aws_iam_policy_document" "assume_lambda_role" {
    statement {
        actions = ["sts:AssumeRole"]

        principals {
            type = "Service"
            identifiers = ["lambda.amazonaws.com"]
        }
    }
}

resource "aws_iam_role" "lambda" {
    name = "AssumeLambdaRole"
    description = "Role for lambda to assume lambda"
    assume_role_policy = data.aws_iam_policy_document.assume_lambda_role.json
}

data "aws_iam_policy_document" "allow_lambda_logging" {
    statement {
        effect = "Allow"
        actions = [
            "logs:CreateLogStream",
            "logs:PutLogEvents",
        ]

        resources = [
            "arn:aws:logs:*:*:*",
        ]
    }
}

resource "aws_iam_policy" "function_logging_policy" {
    name = "AllowLambdaLoggingPolicy"
    description = "Policy for lambda cloudwatch logging"
    policy = data.aws_iam_policy_document.allow_lambda_logging.json
}

resource "aws_iam_role_policy_attachment" "lambda_logging_policy_attachment" {
    role = aws_iam_role.lambda.id
    policy_arn = aws_iam_policy.function_logging_policy.arn
}

data "aws_iam_policy_document" "allow_lambda_secrets" {
    statement {
        effect = "Allow"
        actions = [
            "secretsmanager:GetSecretValue",
        ]

        resources = ["*"]
    }
}

resource "aws_iam_policy" "allow_lambda_access_secrets" {
    name = "AllowLambdaToAccessSecretManager"
    description = "Policy for lambda to access secret manager"
    policy = data.aws_iam_policy_document.allow_lambda_secrets.json
}

resource "aws_iam_role_policy_attachment" "lambda_secrets_policy_attachment" {
    role = aws_iam_role.lambda.id
    policy_arn = aws_iam_policy.allow_lambda_access_secrets.arn
}

resource "aws_iam_role_policy_attachment" "lambda_execution" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role = aws_iam_role.lambda.name
}

resource "aws_lambda_permission" "apigw_lambda" {
    statement_id = "AllowExecutionFromAPIGateway"
    action = "lambda:InvokeFunction"
    function_name = aws_lambda_function.function.function_name
    principal = "apigateway.amazonaws.com"

    source_arn = "${aws_api_gateway_rest_api.my_api.execution_arn}/*/*/*"
}

data "aws_iam_policy_document" "assume_gateway_role" {
    statement {
        actions = ["sts:AssumeRole"]
        effect = "Allow"

        principals {
            type = "Service"
            identifiers = ["apigateway.amazonaws.com"]
        }
    }
}

resource "aws_iam_role" "gateway" {
    name = "AssumeGatewayRole"
    description = "Role for gateway to assume gateway"
    assume_role_policy = data.aws_iam_policy_document.assume_gateway_role.json
}

data "aws_iam_policy_document" "allow_gateway_logging" {
    statement {
        effect = "Allow"
        actions = [
            "logs:CreateLogStream",
            "logs:PutLogEvents",
        ]

        resources = [
            "arn:aws:logs:*:*:*",
        ]
    }
}

resource "aws_iam_policy" "gateway_logging_policy" {
    name = "AllowGatewayLoggingPolicy"
    description = "Policy for gateway cloudwatch logging"
    policy = data.aws_iam_policy_document.allow_gateway_logging.json
}

resource "aws_iam_role_policy_attachment" "gateway_logging_policy_attachment" {
    role = aws_iam_role.gateway.id
    policy_arn = aws_iam_policy.gateway_logging_policy.arn
}