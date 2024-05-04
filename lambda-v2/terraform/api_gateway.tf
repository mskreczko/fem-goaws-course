resource "aws_api_gateway_rest_api" "my_api" {
    name = "my-api"
    description = "My API Gateway"

    endpoint_configuration {
      types = ["REGIONAL"]
    }
}

resource "aws_api_gateway_resource" "root" {
    for_each = local.endpoints
    rest_api_id = aws_api_gateway_rest_api.my_api.id
    parent_id = aws_api_gateway_rest_api.my_api.root_resource_id
    path_part = each.value.path_part
}

resource "aws_api_gateway_method" "api_method" {
    for_each = local.endpoints
    rest_api_id = aws_api_gateway_rest_api.my_api.id
    resource_id = aws_api_gateway_resource.root[each.key].id
    http_method = each.value.method
    authorization = "NONE"
}

resource "aws_api_gateway_integration" "lambda_integration" {
    for_each = local.endpoints
    rest_api_id = aws_api_gateway_rest_api.my_api.id
    resource_id = aws_api_gateway_resource.root[each.key].id
    http_method = aws_api_gateway_method.api_method[each.key].http_method
    integration_http_method = each.value.method
    type = "AWS_PROXY"
    uri = aws_lambda_function.function.invoke_arn
}

resource "aws_api_gateway_method_response" "api_method_response" {
    for_each = local.endpoints
    rest_api_id = aws_api_gateway_rest_api.my_api.id
    resource_id = aws_api_gateway_resource.root[each.key].id
    http_method = aws_api_gateway_method.api_method[each.key].http_method
    status_code = "200"

    response_parameters = {
      "method.response.header.Access-Control-Allow-Headers" = true,
      "method.response.header.Access-Control-Allow-Methods" = true,
      "method.response.header.Access-Control-Allow-Origin" = true
    }

}

resource "aws_api_gateway_integration_response" "api_integration_response" {
    for_each = local.endpoints
    rest_api_id = aws_api_gateway_rest_api.my_api.id
    resource_id = aws_api_gateway_resource.root[each.key].id
    http_method = aws_api_gateway_method.api_method[each.key].http_method
    status_code = aws_api_gateway_method_response.api_method_response[each.key].status_code

    response_parameters = {
      "method.response.header.Access-Control-Allow-Headers" = "'Content-Type,Authorization'",
      "method.response.header.Access-Control-Allow-Methods" = "'GET,OPTIONS,POST,PUT,DELETE'",
      "method.response.header.Access-Control-Allow-Origin" = "'*'"
    }

    depends_on = [
      aws_api_gateway_method.api_method,
      aws_api_gateway_integration.lambda_integration
    ]
}

resource "aws_api_gateway_method" "options" {
    for_each = local.endpoints
    rest_api_id = aws_api_gateway_rest_api.my_api.id
    resource_id = aws_api_gateway_resource.root[each.key].id
    http_method = "OPTIONS"
    authorization = "NONE"
}

resource "aws_api_gateway_integration" "options_integration" {
    for_each = local.endpoints
    rest_api_id             = aws_api_gateway_rest_api.my_api.id
    resource_id             = aws_api_gateway_resource.root[each.key].id
    http_method             = aws_api_gateway_method.options[each.key].http_method
    integration_http_method = "OPTIONS"
    type                    = "MOCK"
    request_templates = {
      "application/json" = "{\"statusCode\": 200}"
    }
}

resource "aws_api_gateway_method_response" "options_response" {
    for_each = local.endpoints
    rest_api_id = aws_api_gateway_rest_api.my_api.id
    resource_id = aws_api_gateway_resource.root[each.key].id
    http_method = aws_api_gateway_method.options[each.key].http_method
    status_code = "200"

    response_parameters = {
      "method.response.header.Access-Control-Allow-Headers" = true,
      "method.response.header.Access-Control-Allow-Methods" = true,
      "method.response.header.Access-Control-Allow-Origin"  = true
    }
}

resource "aws_api_gateway_integration_response" "options_integration_response" {
    for_each = local.endpoints
    rest_api_id = aws_api_gateway_rest_api.my_api.id
    resource_id = aws_api_gateway_resource.root[each.key].id
    http_method = aws_api_gateway_method.options[each.key].http_method
    status_code = aws_api_gateway_method_response.options_response[each.key].status_code

    response_parameters = {
      "method.response.header.Access-Control-Allow-Headers" = "'Content-Type,Authorization'",
      "method.response.header.Access-Control-Allow-Methods" = "'GET,OPTIONS,POST,PUT'",
      "method.response.header.Access-Control-Allow-Origin"  = "'*'"
    }

    depends_on = [
      aws_api_gateway_method.options,
      aws_api_gateway_integration.options_integration,
    ]
}

resource "aws_api_gateway_deployment" "deployment" {
    depends_on = [
      aws_api_gateway_integration.lambda_integration,
      aws_api_gateway_integration.options_integration,
    ]

    rest_api_id = aws_api_gateway_rest_api.my_api.id
    stage_name  = "dev"
}

locals {
    endpoints = {
        register = {
            method = "POST",
            path_part = "register",
            description = "Register endpoint"
        },
        login = {
            method = "POST",
            path_part = "login",
            description = "Login endpoint"
        }
    }
}

