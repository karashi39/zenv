package main

import (
	"fmt"
	"os"
	"path/filepath"
	"gopkg.in/yaml.v3"
	"github.com/spf13/cobra"
)

var configFile = filepath.Join(os.Getenv("HOME"), ".zenv.yml")
var cfg map[string]map[string]interface{}

func loadConfig() {
	cfg = make(map[string]map[string]interface{})
	if data, err := os.ReadFile(configFile); err == nil {
		yaml.Unmarshal(data, &cfg)
	}
}

func saveConfig() {
	data, _ := yaml.Marshal(cfg)
	os.WriteFile(configFile, data, 0644)
}

// ----- add -----
func addCommand(cmd *cobra.Command, args []string) {
	loadConfig()
	varName := args[0]
	if _, exists := cfg[varName]; exists {
		fmt.Println(varName, "is already registered.")
		return
	}
	cfg[varName] = map[string]interface{}{
		"current": "",
		"options": []interface{}{},
	}
	saveConfig()
	fmt.Println(varName, "added.")
}

// ----- list -----
func listCommand(cmd *cobra.Command, args []string) {
	loadConfig()
	if len(cfg) == 0 {
		fmt.Println("No variables registered.")
		return
	}
	for k, v := range cfg {
		current := v["current"].(string)
		fmt.Printf("%s = %s\n", k, current)
	}
}

// ----- set -----
func setCommand(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: zenv set KEY [VALUE]")
		return
	}
	varName := args[0]
	loadConfig()
	varData, ok := cfg[varName]
	if !ok {
		fmt.Println("Variable not registered. Use `zenv add` first.")
		return
	}
	options := varData["options"].([]interface{})
	current := varData["current"].(string)

	if len(args) == 1 {
		// Show options
		for _, opt := range options {
			mark := " "
			if opt.(string) == current {
				mark = "*"
			}
			fmt.Println(mark, opt)
		}
		return
	}

	// Set value
	value := args[1]
	found := false
	for _, opt := range options {
		if opt.(string) == value {
			found = true
			break
		}
	}
	if !found {
		options = append(options, value)
	}
	varData["options"] = options
	varData["current"] = value
	cfg[varName] = varData
	saveConfig()
	fmt.Printf("%s switched to %s.\n", varName, value)
}

// ----- unset -----
func unsetCommand(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: zenv unset KEY")
		return
	}
	varName := args[0]
	loadConfig()
	if varData, ok := cfg[varName]; ok {
		varData["current"] = ""
		cfg[varName] = varData
		saveConfig()
		fmt.Println(varName, "value cleared.")
	} else {
		fmt.Println(varName, "not registered.")
	}
}

// ----- rm -----
func rmCommand(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: zenv rm KEY")
		return
	}
	varName := args[0]
	loadConfig()
	if _, ok := cfg[varName]; ok {
		delete(cfg, varName)
		saveConfig()
		fmt.Println(varName, "removed from management.")
	} else {
		fmt.Println(varName, "not registered.")
	}
}

func main() {
	var rootCmd = &cobra.Command{Use: "zenv"}

	var addCmd = &cobra.Command{
		Use:   "add [KEY]",
		Short: "Add a new environment variable",
		Args:  cobra.ExactArgs(1),
		Run:   addCommand,
	}
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List registered environment variables and current value",
		Run:   listCommand,
	}
	var setCmd = &cobra.Command{
		Use:   "set [KEY] [VALUE]",
		Short: "Show options for a variable or set its value",
		Args:  cobra.RangeArgs(1, 2),
		Run:   setCommand,
	}
	var unsetCmd = &cobra.Command{
		Use:   "unset [KEY]",
		Short: "Clear value of a variable",
		Args:  cobra.ExactArgs(1),
		Run:   unsetCommand,
	}
	var rmCmd = &cobra.Command{
		Use:   "rm [KEY]",
		Short: "Remove a variable from management",
		Args:  cobra.ExactArgs(1),
		Run:   rmCommand,
	}

	rootCmd.AddCommand(addCmd, listCmd, setCmd, unsetCmd, rmCmd)
	rootCmd.Execute()
}
