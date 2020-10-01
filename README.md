# Azure Monitoring Policies

This will add policies to your subscription that can guarantee that all resources that are mentioned in the [documentation](https://docs.microsoft.com/en-us/azure/azure-monitor/platform/resource-logs-categories) do have monitoring enabled.

## Design Decision

This repository will only create the policy definitions and them you would need to apply this definition to the resource you want to. This is a deliberated decision given that you might want tto create the policy deficinion in the tenant level or the mangment group level and them apply this definition to your subscriptions in the subscription level. If you apply that in the managment group level all the logs from all the subscriptions will be directed to the same storage.

By default the policy is created inside the root managment group since we do nt requere one to be created.

## How to use it

Using the module itself does not require any configuration since it just create the  

```terraform
module "policies" {
  source = "git@github.com:Nepomuceno/terraform-azurerm-monitoring-policies.git"
}
```

You can also add some configurations.

```terraform
module "policies" {
  source                = "git@github.com:Nepomuceno/terraform-azurerm-monitoring-policies.git"
  name                  = "My policy name"
  management_group_name = "My managment group name"
}
```

To use the policy you need to assign that to you subcription:

```terraform
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
```

you can find examples of the implementations at the `examples` folder
