terraform {
  required_version = ">= 0.13"
}

provider "azurerm" {
  features {}
}


resource "azurerm_policy_definition" "base" {
  for_each     = fileset("${path.module}/templates", "**/rule.json")
  name         = substr("log-${replace(each.value, "/rule.json", "")}", 0, 54)
  policy_type  = "Custom"
  mode         = "All"
  display_name = "Diagnostic policy ${replace(each.value, "/rule.json", "")}"
  policy_rule = file(
    "${path.module}/templates/${each.value}",
  )
  parameters = file(
    "${path.module}/templates/${replace(each.value, "/rule.json", "/parameters.json")}",
  )
}


resource "azurerm_policy_set_definition" "basic_set" {
  name         = ""
  policy_type  = "Custom"
  display_name = "Auto Diagnostics Policy Initiative"

  parameters = <<PARAMETERS
    {
    "requiredRetentionDays": {
      "type": "String",
      "metadata": {
        "displayName": "Required retention (days)",
        "description": "The required diagnostic logs retention in days"
      },
      "defaultValue": "365"
    },
    "eventHubName": {
        "type": "String",
        "metadata":{
            "displayName": "Event hub to send the data to",
            "description": ""
        },
        "defaultValue": ""
    },
    "eventHubAuthorizationRuleId": {
        "type": "String",
        "metadata":{
            "displayName": "Event hub rule to be used to send data",
            "description": ""
        },
        "defaultValue": ""
    },
    "workspaceId": {
        "type": "String",
        "metadata":{
            "displayName": "Log analytics workspace id to send the data to",
            "description": ""
        },
        "defaultValue": ""
    },
    "storageAccountName": {
        "type": "String",
        "metadata":{
            "displayName": "Storage account to send the data to",
            "description": ""
        },
        "defaultValue": ""
    },
    "resourceLocation": {
        "type": "array",
        "metadata": {
          "description": "locations that you want to enable diagnotics to",
          "displayName": "location where disgnostics will be enabled",
          "strongType": "location"
        },
        "defaultValue": [
            "eastus",
            "eastus2",
            "southcentralus",
            "westus2",
            "australiaeast",
            "southeastasia",
            "northeurope",
            "uksouth",
            "westeurope",
            "centralus",
            "northcentralus",
            "westus",
            "southafricanorth",
            "centralindia",
            "eastasia",
            "japaneast",
            "koreacentral",
            "canadacentral",
            "francecentral",
            "germanywestcentral",
            "norwayeast",
            "switzerlandnorth",
            "uaenorth",
            "brazilsouth",
            "centralusstage",
            "eastusstage",
            "eastus2stage",
            "northcentralusstage",
            "southcentralusstage",
            "westusstage",
            "westus2stage",
            "asia",
            "asiapacific",
            "australia",
            "brazil",
            "canada",
            "europe",
            "global",
            "india",
            "japan",
            "uk",
            "unitedstates",
            "eastasiastage",
            "southeastasiastage",
            "eastus2euap",
            "westcentralus",
            "southafricawest",
            "australiacentral",
            "australiacentral2",
            "australiasoutheast",
            "japanwest",
            "koreasouth",
            "southindia",
            "westindia",
            "canadaeast",
            "francesouth",
            "germanynorth",
            "norwaywest",
            "switzerlandwest",
            "ukwest",
            "uaecentral",
            "brazilsoutheast"
          ]
    }
}
PARAMETERS

  policy_definitions = jsonencode([
    for def in azurerm_policy_definition.base : {
      parameters = {
        requiredRetentionDays = {
          value = "[parameters('requiredRetentionDays')]"
        },
        eventHubName = {
          value = "[parameters('eventHubName')]"
        },
        eventHubAuthorizationRuleId = {
          value = "[parameters('eventHubAuthorizationRuleId')]"
        },
        workspaceId = {
          value = "[parameters('workspaceId')]"
        },
        storageAccountName = {
          value = "[parameters('storageAccountName')]"
        },
        resourceLocation = {
          value = "[parameters('resourceLocation')]"
        },
      },
      policyDefinitionId = def.id
    }
  ])

  lifecycle {
    ignore_changes = [metadata] // hacky hack hack, always says it has changed!
  }
}
