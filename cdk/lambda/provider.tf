terraform {
    required_providers {
        aws = {
            source = "hashicorp/aws"
            version = "~> 5.0"
        }
        archive = {
            source = "hashicorp/archive"
        }
        null = {
            source = "hashicorp/null"
        }
    }
}

provider "aws" {
  region = "us-west-2"
  shared_credentials_files = ["/home/mskreczko/.aws/credentials"]
  profile = "mskreczko"
}

resource "aws_vpc" "example_vpc" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "example_subnet" {
    vpc_id = aws_vpc.example_vpc.id
    cidr_block = "10.0.1.0/24"
}

resource "aws_security_group" "lambda_sg" {
    vpc_id = aws_vpc.example_vpc.id

    ingress {
        from_port = 0
        to_port = 0
        protocol = "-1"
        cidr_blocks = ["0.0.0.0/0"]
    }

    egress {
        from_port = 0
        to_port = 0
        protocol = "-1"
        cidr_blocks = ["0.0.0.0/0"]
    }
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

resource "null_resource" "function_binary" {
    provisioner "local-exec" {
        command = "GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ${local.binary_path} ${local.src_path}"
    }
}

data "archive_file" "function_archive" {
    depends_on = [ null_resource.function_binary ]

    type = "zip"
    source_file = local.binary_path
    output_path = local.archive_path
}

resource "aws_lambda_function" "function" {
    function_name = "hello-world"
    description = "My first lambda function"
    role = aws_iam_role.lambda.arn
    handler = local.binary_name
    memory_size = 128

    filename = local.archive_path
    source_code_hash = data.archive_file.function_archive.output_base64sha256

    runtime = "provided.al2023"
}

resource "aws_cloudwatch_log_group" "log_group" {
    name = "/aws/lambda/${aws_lambda_function.function.function_name}"
    retention_in_days = 7
}

locals {
    function_name = "bootstrap"
    src_path = "${path.module}"

    binary_name = local.function_name
    binary_path = "${path.module}/tf_generated/${local.binary_name}"
    archive_path = "${path.module}/tf_generated/${local.function_name}.zip"
}

output "binary_path" {
    value = local.binary_path
}