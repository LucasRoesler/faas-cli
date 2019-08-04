package commands

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/openfaas/faas-cli/stack"
	"github.com/spf13/cobra"
)

var (
	templateURL string
	customName  string
)

func init() {
	templatePullStackCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing templates?")
	templatePullStackCmd.Flags().BoolVar(&pullDebug, "debug", false, "Enable debug output")
	templatePullStackCmd.PersistentFlags().StringVarP(&customName, "name", "n", "", "The custom name of the template")

	templatePullCmd.AddCommand(templatePullStackCmd)
}

var templatePullStackCmd = &cobra.Command{
	Use:   `stack`,
	Short: `Downloads templates specified in the function definition yaml file`,
	Long: `Downloads templates specified in the function yaml file, in the current directory
	`,
	Example: `
  faas-cli template pull stack
  faas-cli template pull stack -f myfunction.yml
`,
	RunE: runTemplatePullStack,
}

func runTemplatePullStack(cmd *cobra.Command, args []string) error {
	templatesConfig, err := loadTemplateConfig()
	if err != nil {
		return err
	}
	err = pullTemplatesFromConfig(templatesConfig, customName)
	if err != nil {
		return err
	}
	return nil
}

func loadTemplateConfig() (map[string]string, error) {
	stackConfig, err := readStackConfig()
	if err != nil {
		return nil, err
	}
	return stackConfig.StackConfig.TemlpatesConfigs, nil
}

func readStackConfig() (stack.Configuration, error) {
	configField := stack.Configuration{}

	configFieldBytes, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return configField, fmt.Errorf("Error while reading files %s", err.Error())
	}
	unmarshallErr := yaml.Unmarshal(configFieldBytes, &configField)
	if unmarshallErr != nil {
		return configField, fmt.Errorf("Error while reading configuration: %s", err.Error())
	}
	if len(configField.StackConfig.TemlpatesConfigs) == 0 {
		return configField, fmt.Errorf("Error while reading configuration: no template repos currently configured")
	}
	return configField, nil
}

func pullTemplatesFromConfig(templateInfo map[string]string, customName string) error {
	if len(customName) > 0 {
		templateInfo = findTemplate(templateInfo, customName)
		if templateInfo == nil {
			return fmt.Errorf("Unable to find template with name: `%s`", customName)
		}
	}
	for key, val := range templateInfo {
		fmt.Printf("Pulling template: `%s` from configuration file: `%s`\n", key, yamlFile)
		pullErr := pullTemplate(val)
		if pullErr != nil {
			return pullErr
		}
	}
	return nil
}

func findTemplate(templateInfo map[string]string, customName string) (specificTemplate map[string]string) {
	for key, val := range templateInfo {
		if key == customName {
			return map[string]string{key: val}
		}
	}
	return nil
}
