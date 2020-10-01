terraform {
  required_version = ">= 0.13"
  
}

provider "azurerm" {
  features {}
}

module "policies" {
  source = "git@github.com:Nepomuceno/terraform-azurerm-monitoring-policies.git"
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "uksouth"
}


resource "azurerm_log_analytics_workspace" "example" {
  name                = "acctest-11-diagn"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_policy_assignment" "example" {
  name                 = "example-diag-policy-assignment"
  scope                = azurerm_resource_group.example.id
  policy_definition_id = module.policies.initiative.id
  description          = "Policy Assignment created via terraform"
  display_name         = "Diagnostic Logs application"
  location             = "uksouth"
  identity {
      type = "SystemAssigned"
  }

  metadata = <<METADATA
    {
    "category": "Logs"
    }
METADATA

  parameters = <<PARAMETERS
{
  "workspaceId": {
    "value": "${azurerm_log_analytics_workspace.example.id}"
  }
}
PARAMETERS

}