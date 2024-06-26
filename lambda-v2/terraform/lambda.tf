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
    timeout = 60

    filename = local.archive_path
    source_code_hash = data.archive_file.function_archive.output_base64sha256

    runtime = "provided.al2023"
}

resource "aws_cloudwatch_log_group" "log_group" {
    name = "/aws/lambda/${aws_lambda_function.function.function_name}"
    retention_in_days = 7
}