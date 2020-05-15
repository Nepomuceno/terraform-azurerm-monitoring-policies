package generator

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// LogStructure Current log structure
type LogStructure struct {
	HasMetrics   bool
	Categories   []string
	ResourceType string
}

func formatName(name string) string {
	return strings.ToLower(strings.Replace(strings.Replace(name, "/", "_", -1), ".", "_", -1))
}

func getDefinitions() (map[string]LogStructure, error) {
	metrics, err := getMetrics()
	// Getting data from the azure
	resp, err := http.Get("https://raw.githubusercontent.com/MicrosoftDocs/azure-docs/master/articles/azure-monitor/platform/diagnostic-logs-schema.md")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	content := string(body)
	var foundSession, endSession bool = false, false
	response := make(map[string]LogStructure)
	for _, line := range strings.Split(content, "\n") {
		if line == "|Resource Type|Category|Category Display Name|" {
			foundSession = true
			continue
		}
		if foundSession && !endSession {
			if len(line) == 0 {
				endSession = true
				continue
			}
			if line == "|---|---|---|" || line == "|Resource Type|Category|Category Display Name|" {
				continue
			}
			logCategory := strings.Split(line, "|")
			logName := formatName(logCategory[1])
			cat, exist := response[logName]
			if exist {
				cat.Categories = append(cat.Categories, logCategory[2])
			} else {
				cat.ResourceType = logCategory[1]
				cat.Categories = []string{logCategory[2]}
				_, cat.HasMetrics = metrics[logName]
			}
			response[logName] = cat
		}
	}
	return response, nil
}

func getMetrics() (map[string]bool, error) {
	// Currently the only way to check whihc resources do support metrics.
	resp, err := http.Get("https://raw.githubusercontent.com/MicrosoftDocs/azure-docs/master/articles/azure-monitor/platform/metrics-supported.md")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	content := string(body)
	response := make(map[string]bool)
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "##") && strings.ContainsAny(line, "/") {
			resourceID := strings.Split(line, " ")[1]
			response[formatName(resourceID)] = true
		}
	}
	return response, nil
}

// Generate the role definitions
func Generate() error {
	logCategories, err := getDefinitions()
	if err != nil {
		return err
	}
	temp, err := getTemplates()
	if err != nil {
		return err
	}
	outputPath := os.Getenv("GENERATOR_OUTPUT_PATH")
	available := make([]string, 0)
	if len(outputPath) == 0 {
		outputPath = "./templates"
	}
	for k, content := range logCategories {
		available = append(available, content.ResourceType)
		os.MkdirAll(fmt.Sprintf("%s/%s/", outputPath, k), os.ModePerm)
		fr, err := os.Create(fmt.Sprintf("%s/%s/rule.json", outputPath, k))
		if err != nil {
			return err
		}
		_ = temp.ExecuteTemplate(fr, ruleTemplate, content)
		fp, err := os.Create(fmt.Sprintf("%s/%s/parameters.json", outputPath, k))
		if err != nil {
			return err
		}
		_ = temp.ExecuteTemplate(fp, paramTemplate, nil)
	}
	os.MkdirAll(outputPath, os.ModePerm)
	fa, err := os.Create(fmt.Sprintf("%s/available_resources.json", outputPath))
	if err != nil {
		return err
	}
	_ = temp.ExecuteTemplate(fa, generatedTemplate, available)
	return nil
}
