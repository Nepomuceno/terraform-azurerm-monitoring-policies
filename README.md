# Azure Monitoring Policies

This will add policies to your subscription that can guarantee that all resources that are mentioned in the [documentation](https://docs.microsoft.com/en-us/azure/azure-monitor/platform/resource-logs-categories) have monitoring enabled.

## Design Decision

This repository creates the policy definitions and a policy initiative, which is a collection of all of the policy definitions. You would need to apply the relevant definition (or the initiative) to the resource you want to monitor. This is a deliberate decision given that you might want to create the policy definition at the tenant level or the management group level and then apply the definition at the subscription or resource group level. If you assign a definition at the management group level all the logs from matching resources in all child subscriptions will be directed to the same log collection point (e.g. log analytics, event hub, storage blob).

If you do not specify one, the policy is created inside the tenant root management group since we do not require one to be created.

## How to use it

Using the module itself does not require any configuration since it just creates the definitions.  
According to [this link](https://www.terraform.io/docs/modules/sources.html#github):

```terraform
module "policies" {
  # use this form to clone using ssh creds
  # source = "git@github.com:Nepomuceno/terraform-azurerm-monitoring-policies.git"
  # and this one to use https
  source = "github.com/Nepomuceno/terraform-azurerm-monitoring-policies.git"
}
```

AzureRM version 2.29.0 or greater is required due to ['parameters' deprecation](https://github.com/terraform-providers/terraform-provider-azurerm/pull/8270).

You can also configure policy definition creation by overriding these terraform variables.

```terraform
module "policies" {
  source                = "git@github.com:Nepomuceno/terraform-azurerm-monitoring-policies.git"
  name                  = "My policy name"
  management_group_name = "My managment group name"
}
```

To use a policy you need to assign it to either a resource group, a subscription or a management group. In this example, the policy initiative is applied to a resource group. The logs are directed to a log analytics workspace.  

```terraform
terraform {
  required_version = ">= 0.13"
  
}

provider "azurerm" {
  version = ">=2.29.0"
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

You can find examples of the implementations at the `examples` folder.
