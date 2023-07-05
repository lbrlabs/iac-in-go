package cluster

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/lbrlabs/iac-in-go/pkg/cluster/awsauth"
)

func appendToStringArrayInput(ctx *pulumi.Context, arrayInput pulumi.StringArrayInput, value pulumi.StringOutput) pulumi.StringArrayInput {
	// Initialize arrayInput if it's nil
	if arrayInput == nil {
		arrayInput = pulumi.StringArray([]pulumi.StringInput{})
	}
	return pulumi.All(arrayInput, value).ApplyT(func(args []interface{}) ([]string, error) {
		arr := args[0].([]string)
		v := args[1].(string)
		return append(arr, v), nil
	}).(pulumi.StringArrayOutput)
}

func ArnToRoleMapping(arns []string) []awsauth.RoleMapping {
	mappings := []awsauth.RoleMapping{}
	for _, arn := range arns {
		mapping := awsauth.RoleMapping{
			RoleArn: arn,
			Username: "system:node:{{EC2PrivateDNSName}}",
			Groups: []string{"system:bootstrappers", "system:nodes"},
		}
		mappings = append(mappings, mapping)
	}
	return mappings
}
