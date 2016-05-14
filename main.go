package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var re = regexp.MustCompile("[-.+~`!@#$%^&*(){}\\[\\]:;\"',?<>/]")

func cleanVar(name string) string {
	return re.ReplaceAllLiteralString(name, "_")
}

type envvar struct {
	Key   string
	Value string
}

func (e envvar) String() string {
	return fmt.Sprintf(`export %s=%s`, e.Key, e.Value)
}

func flatten(prefix string, vars map[string]interface{}) []envvar {
	envvars := []envvar{}
	for key, val := range vars {
		varname := cleanVar(prefix + key)
		switch val.(type) {
		case map[string]interface{}:
			envvars = append(envvars,
				flatten(
					fmt.Sprintf("%s_", varname),
					val.(map[string]interface{}))...,
			)
		default:
			envvars = append(envvars,
				envvar{
					Key:   strings.ToUpper(varname),
					Value: fmt.Sprintf("%#v", val),
				})

		}
	}
	return envvars
}

func Process(vcap string) []string {
	var svcs map[string][]Service
	json.Unmarshal([]byte(vcap), &svcs)
	vars := []string{}

	for _, instances := range svcs {
		for _, instance := range instances {
			for _, cred := range flatten(instance.Name+"_", instance.Credentials) {
				vars = append(vars, cred.String())
			}
		}
	}
	return vars
}

// Service is a VCAP_SERVICE instance
type Service struct {
	Name        string                 `json:"name"`
	Credentials map[string]interface{} `json:"credentials"`
}

func main() {
	fmt.Println(strings.Join(Process(os.Getenv("VCAP_SERVICES")), "\n"))
}
