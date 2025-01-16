resource "allquiet_service" "payment_api" {
  display_name = "Payment API"
  public_title = "Payment API"
  public_description = "Payments"
  templates = [
    {
      display_name = "Refunds delayed"
      message = "Refunds are currently delayed. All refunds will be processed but can currently take longer than usual to complete."
    },
    {
      display_name = "Payment gateway down"
      message = "Our payment gateway is currently down. We are working to resolve the issue as soon as possible."
    }
  ]
}
