package main

import (
	"fmt"

	xec2 "github.com/pulumi/pulumi-awsx/sdk/go/awsx/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		conf := config.New(ctx, "")
		cidrBlock := conf.Require("cidrBlock")
		numberOfAvailabilityZones := conf.GetInt("numberOfAvailabilityZones")

		org := ctx.Organization()
		project := ctx.Project()

		// can probably make this a function
		var azs int
		if numberOfAvailabilityZones == 0 {
			azs = 2
		} else {
			azs = numberOfAvailabilityZones
		}

		// if we have 3 azs, assume it's a production vpc
		// in which case, we want to have a nat gateway per az
		var NatGatewayStrategy xec2.NatGatewayStrategy
		if numberOfAvailabilityZones > 3 {
			NatGatewayStrategy = xec2.NatGatewayStrategyOnePerAz
		} else {
			NatGatewayStrategy = xec2.NatGatewayStrategySingle
		}

		tags := pulumi.StringMap{
			"Name":          pulumi.String("public"),
			"Owner":         pulumi.String("lbriggs"),
			"PulumiOrg":     pulumi.String(org),
			"PulumiProject": pulumi.String(project),
		}

		vpc, err := xec2.NewVpc(ctx, "vpc", &xec2.VpcArgs{
			CidrBlock:                 &cidrBlock,
			NumberOfAvailabilityZones: &azs,
			NatGateways: &xec2.NatGatewayConfigurationArgs{
				Strategy: NatGatewayStrategy,
			},
			EnableDnsSupport:   pulumi.Bool(true),
			EnableDnsHostnames: pulumi.Bool(true),
			SubnetSpecs: []xec2.SubnetSpecArgs{
				{
					Type: xec2.SubnetTypePublic,
					Tags: pulumi.StringMap{
						"Name":                   pulumi.String("public"),
						"Owner":                  pulumi.String("lbriggs"),
						"PulumiOrg":              pulumi.String(org),
						"PulumiProject":          pulumi.String(project),
						"kubernetes.io/role/elb": pulumi.String("1"),
					},
				},
				{
					Type: xec2.SubnetTypePrivate,
					Tags: pulumi.StringMap{
						"Name":                            pulumi.String("public"),
						"Owner":                           pulumi.String("lbriggs"),
						"PulumiOrg":                       pulumi.String(org),
						"PulumiProject":                   pulumi.String(project),
						"kubernetes.io/role/internal-elb": pulumi.String("1"),
					},
				},
			},
			Tags: tags,
		})
		if err != nil {
			return fmt.Errorf("error creating vpc: %v", err)
		}

		ctx.Export("vpcId", vpc.VpcId)
		ctx.Export("privateSubnetIds", vpc.PrivateSubnetIds)
		ctx.Export("publicSubnetIds", vpc.PublicSubnetIds)
		return nil
	})
}
