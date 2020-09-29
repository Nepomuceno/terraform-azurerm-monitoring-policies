package main

//go:generate go run main.go

import (
	"github.com/nepomuceno/terraform-azurerm-monitoring-policies/generator"
)

func main() {
	generator.Generate()
}
