---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "allquiet_service Resource - allquiet"
subcategory: ""
description: |-
  The service resource represents a service in All Quiet.
---

# allquiet_service (Resource)

The `service` resource represents a service in All Quiet.

## Example Usage

```terraform
resource "allquiet_service" "payment_api" {
  display_name       = "Payment API"
  public_title       = "Payment API"
  public_description = "Payments"
  templates = [
    {
      display_name = "Refunds delayed"
      message      = "Refunds are currently delayed. All refunds will be processed but can currently take longer than usual to complete."
    },
    {
      display_name = "Payment gateway down"
      message      = "Our payment gateway is currently down. We are working to resolve the issue as soon as possible."
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) The display name of the service
- `public_title` (String) The public title of the service

### Optional

- `public_description` (String) The public description of the service
- `templates` (Attributes List) The templates of the service (see [below for nested schema](#nestedatt--templates))

### Read-Only

- `id` (String) Id

<a id="nestedatt--templates"></a>
### Nested Schema for `templates`

Required:

- `display_name` (String) The display name of the template
- `message` (String) The message of the template

Read-Only:

- `id` (String) Id