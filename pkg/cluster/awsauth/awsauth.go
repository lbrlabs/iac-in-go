package awsauth

import (
	"gopkg.in/yaml.v2"
)
type RoleMapping struct {
	RoleArn  string   `yaml:"roleArn"`
	Username string   `yaml:"username"`
	Groups   []string `yaml:"groups"`
}

func BuildAuthData(arns []string) ([]byte, error) {

	roles := []RoleMapping{}

	for _, arn := range arns {
		role := RoleMapping{
			RoleArn:  arn,
			Username: "system:node:{{EC2PrivateDNSName}}",
			Groups:   []string{"system:bootstrappers", "system:nodes"},
		}
		roles = append(roles, role)
	}

	return yaml.Marshal(roles)
}
