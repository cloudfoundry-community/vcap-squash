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
	return strings.ToUpper(re.ReplaceAllLiteralString(name, "_"))
}

type envvar struct {
	Key   string
	Value string
}

func (e envvar) String() string {
	return fmt.Sprintf(`export %s=%s`, e.Key, e.Value)
}

func flattenarr(prefix string, vararr []interface{}) []envvar {
	envvars := []envvar{}
	for key, val := range vararr {
		envvars = append(envvars,
			flatten(fmt.Sprintf("%s%v", prefix, key), val)...)
	}
	return envvars
}

func flattenmap(prefix string, varmap map[string]interface{}) []envvar {
	envvars := []envvar{}
	for key, val := range varmap {
		envvars = append(envvars,
			flatten(prefix+key, val)...)
	}
	return envvars
}

func flatten(prefix string, vars interface{}) []envvar {
	switch vars.(type) {
	case map[string]interface{}:
		return flattenmap(prefix+"_", vars.(map[string]interface{}))
	case []interface{}:
		return flattenarr(prefix+"_", vars.([]interface{}))
	default:
		return []envvar{envvar{
			Key:   cleanVar(prefix),
			Value: fmt.Sprintf("%#v", vars),
		}}
	}
}

func Process(vcap string) []string {
	var svcs map[string][]Service
	json.Unmarshal([]byte(vcap), &svcs)
	vars := []string{}

	for _, instances := range svcs {
		for _, instance := range instances {
			for _, cred := range flatten(instance.Name, instance.Credentials) {
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
