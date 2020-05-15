terraform {
  required_version = ">= 0.12"
}

provider "azurerm" {
  features {}
}


resource "azurerm_policy_definition" "base" {
  for_each     = fileset("${path.module}/out", "**/rule.json")
  name         = substr("log-${replace(each.value, "/rule.json", "")}", 0, 54)
  policy_type  = "Custom"
  mode         = "All"
  display_name = "Diagnostic policy ${replace(each.value, "/rule.json", "")}"
  policy_rule = file(
    "${path.module}/out/${each.value}",
  )
  parameters = file(
    "${path.module}/out/${replace(each.value, "/rule.json", "/parameters.json")}",
  )
}


resource "azurerm_policy_set_definition" "basic_set" {
  name         = "Auto Diagnostics policy initiative"
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
        "type": "String",
        "metadata": {
          "description": "location to put this diagnostic",
          "displayName": "location where this diagnostic is"
        }
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
