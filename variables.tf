# Define config variables
variable "labelPrefix" {
  type        = string
  description = "Your college username. This will form the beginning of various resource names."
}

variable "region" {
  default = "canadacentral"
  type        = string
  description = "Azure region for all resources. Use a region allowed by your subscription policy."
}

variable "admin_username" {
  type        = string
  default     = "azureadmin"
  description = "The username for the local user account on the VM."
}
