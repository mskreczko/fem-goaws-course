resource "aws_dynamodb_table" "users-table" {
    name = "users"
    billing_mode = "PAY_PER_REQUEST"
    hash_key = "email"

    attribute {
        name = "Email"
        type = "S"
    }

    attribute {
        name = "Password"
        type = "S"
    }
}