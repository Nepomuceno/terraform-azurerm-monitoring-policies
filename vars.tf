variable "name" {
  type        = string
  description = "Name of the initiative to be created"
  default     = "Auto Diagnostics Policy Initiative"
}

variable "management_group_name" {
  type        = string
  description = "Mangment group that the policy should be created"
  default     = null
}