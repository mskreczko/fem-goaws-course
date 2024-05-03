terraform {
    required_providers {
      aws = {
        source = "hashicorp/aws"
        version = "~> 5.0"
      }
      aws_dynamodb_table = {
        source = "terraform-aws-modules/dynamodb-table/aws"
        version = "4.0.1"
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