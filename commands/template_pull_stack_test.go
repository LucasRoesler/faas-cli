package commands

import (
	"reflect"
	"testing"
)

func Test_findTemplate(t *testing.T) {
	tests := []struct {
		title             string
		desiredTemplate   string
		existingTempaltes map[string]string
		expectedTemplate  map[string]string
	}{
		{
			title:             "Desired template is found",
			desiredTemplate:   "powershell",
			existingTempaltes: map[string]string{"powershell": "exampleURL", "rust": "exampleURL"},
			expectedTemplate:  map[string]string{"powershell": "exampleURL"},
		},
		{
			title:             "Desired template is not found",
			desiredTemplate:   "golang",
			existingTempaltes: map[string]string{"powershell": "exampleURL", "rust": "exampleURL"},
			expectedTemplate:  nil,
		},
	}
	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			result := findTemplate(test.existingTempaltes, test.desiredTemplate)
			if !reflect.DeepEqual(result, test.expectedTemplate) {
				t.Errorf("Wanted template: `%s` got `%s`", test.expectedTemplate, result)
			}
		})
	}
}

func Test_pullTemplatesFromConfig(t *testing.T) {
	tests := []struct {
		title             string
		desiredTemplate   string
		existingTempaltes map[string]string
		expectedError     bool
	}{
		{
			title:           "Pull specific tempalte",
			desiredTemplate: "my_powershell",
			existingTempaltes: map[string]string{
				"my_powershell": "https://github.com/openfaas-incubator/powershell-http-template",
				"my_rust":       "https://github.com/openfaas-incubator/openfaas-rust-template"},
			expectedError: false,
		},
		{
			title:           "Pull all templates",
			desiredTemplate: "",
			existingTempaltes: map[string]string{
				"my_powershell": "https://github.com/openfaas-incubator/powershell-http-template",
				"my_rust":       "https://github.com/openfaas-incubator/openfaas-rust-template"},
			expectedError: false,
		},
		{
			title:           "Pull non-existant template",
			desiredTemplate: "my_golang",
			existingTempaltes: map[string]string{
				"my_powershell": "exampleURL",
				"my_rust":       "https://github.com/openfaas-incubator/openfaas-rust-template"},
			expectedError: true,
		},
		{
			title:           "Pull template with invalid URL",
			desiredTemplate: "my_golang",
			existingTempaltes: map[string]string{
				"my_powershell": "invalidURL",
			},
			expectedError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			actualError := pullTemplatesFromConfig(test.existingTempaltes, test.desiredTemplate)
			if actualError != nil && test.expectedError == false {
				t.Errorf("Unexpected error: %s", actualError.Error())
			}
		})
	}
}
