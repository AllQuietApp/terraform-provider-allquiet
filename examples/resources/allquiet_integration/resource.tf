resource "allquiet_team" "root" {
  display_name = "Root"
}

resource "allquiet_integration" "datadog" {
  display_name = "My Datadog Integration"
  team_id      = allquiet_team.root.id
  type         = "Datadog"
}

resource "allquiet_integration" "amazon_cloudwatch" {
  display_name = "My Amazon CloudWatch Integration"
  team_id      = allquiet_team.root.id
  type         = "AmazonCloudWatch"
}

locals {
  computed_amazon_cloudwatch_webhook_url = allquiet_integration.amazon_cloudwatch.webhook_url
}