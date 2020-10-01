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

