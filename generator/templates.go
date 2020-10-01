package generator

import "html/template"

const templateRule = `{
    "if": {
        "allOf": [
            {
              "field": "location",
              "in": "[parameters('resourceLocation')]"
            },
            {
              "field": "type",
              "equals": "{{.ResourceType}}"
            }
        ]
    },
    "then": {
        "effect": "DeployIfNotExists",
        "details": {
            "type": "Microsoft.Insights/diagnosticSettings",
            "existenceCondition": {
                "anyOf": [
                    {
                        "allOf": [
                            {
                                "field": "Microsoft.Insights/diagnosticSettings/logs[*].retentionPolicy.enabled",
                                "equals": "true"
                            },
                            {
                                "field": "Microsoft.Insights/diagnosticSettings/logs[*].retentionPolicy.days",
                                "equals": "[parameters('requiredRetentionDays')]"
                            },
                            {
                                "field": "Microsoft.Insights/diagnosticSettings/logs.enabled",
                                "equals": "true"
                            }
                        ]
                    },
                    {
                        "allOf": [
                            {
                                "not": {
                                    "field": "Microsoft.Insights/diagnosticSettings/logs[*].retentionPolicy.enabled",
                                    "equals": "true"
                                }
                            },
                            {
                                "field": "Microsoft.Insights/diagnosticSettings/logs.enabled",
                                "equals": "true"
                            }
                        ]
                    }
                ]
            },
            "roleDefinitionIds": [
                "/providers/Microsoft.Authorization/roleDefinitions/b24988ac-6180-42a0-ab88-20f7382dd24c"
            ],
            "deployment": {
                "properties": {
                    "mode": "incremental",
                    "template": {
                        "$schema": "http://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
                        "contentVersion": "1.0.0.0",
                        "parameters": {
                            "name": {
                                "type": "string"
                            },
                            "id": {
                                "type": "string"
                            },
                            "eventHubName": {
                                "type": "string"
                            },
                            "eventHubAuthorizationRuleId": {
                                "type": "string"
                            },
                            "workspaceId": {
                                "type": "string"
                            },
                            "storageAccountName": {
                                "type": "string"
                            },
                            "retentionDays": {
                                "type": "string"
                            }
                        },
                        "variables": {
                            "ehEnabled": "[greater(length(parameters('eventHubName')),0)]",
                            "laEnabled": "[greater(length(parameters('workspaceId')),0)]",
                            "saEnabled": "[greater(length(parameters('storageAccountName')),0)]"

                        },
                        "resources": [
                            {
                                "type": "{{.ResourceType}}/providers/diagnosticSettings",
                                "name": "[concat(parameters('name'), '/', 'Microsoft.Insights/setByPolicy')]",
                                "dependsOn": [],
                                "apiVersion": "2017-05-01-preview",
                                "properties": {
                                    "storageAccountId": "[if(variables('saEnabled'),resourceId('Microsoft.Storage/storageAccounts', parameters('storageAccountName')),json('null'))]",
                                    "eventHubAuthorizationRuleId": "[if(variables('ehEnabled'),parameters('eventHubAuthorizationRuleId'),json('null'))]",
                                    "eventHubName": "[if(variables('ehEnabled'),parameters('eventHubName'),json('null'))]",
                                    "workspaceId": "[if(variables('laEnabled'),parameters('workspaceId'),json('null'))]",
                                    "logs": [
										{{range  $index, $element := .Categories}}
                                        {{if $index}},{{end}}
                                        {
											"category": "{{$element}}",
											"enabled": true,
											"retentionPolicy": {
												"days": "[parameters('retentionDays')]",
												"enabled": true
											}
										}
										{{end}}
                                    ],
                                    "metrics": [
										{{if .HasMetrics}}
											{
												"category": "AllMetrics",
												"enabled": true,
												"retentionPolicy": {
													"enabled": true,
													"days": "[parameters('retentionDays')]"
												}
											}
										{{end}}
                                    ]
                                }
                            }
                        ]
                    },
                    "parameters": {
                        "name": {
                            "value": "[field('name')]"
                        },
                        "id": {
                            "value": "[field('fullName')]"
                        },
                        "eventHubName": {
                            "value": "[parameters('eventHubName')]"
                        },
                        "eventHubAuthorizationRuleId": {
                            "value": "[parameters('eventHubAuthorizationRuleId')]"
                        },
                        "workspaceId": {
                            "value": "[parameters('workspaceId')]"
                        },
                        "storageAccountName": {
                            "value": "[parameters('storageAccountName')]"
                        },
                        "retentionDays": {
                            "value": "[parameters('requiredRetentionDays')]"
                        }
                    }
                }
            }
        }
    }
}
`

const templateParam = `{
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
        "type": "Array",
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
`

const templateGenerated = `
[
    {{range .}}"{{.}}",
    {{end}}
]
`

const (
	paramTemplate     = "param"
	ruleTemplate      = "rule"
	generatedTemplate = "generated"
)

func getTemplates() (*template.Template, error) {
	temp, err := template.New(paramTemplate).Parse(templateParam)
	if err != nil {
		return temp, err
	}
	temp, err = temp.New(ruleTemplate).Parse(templateRule)
	if err != nil {
		return temp, err
	}
	temp, err = temp.New(generatedTemplate).Parse(templateGenerated)
	if err != nil {
		return temp, err
	}
	return temp, nil
}
