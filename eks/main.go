package main

import (
	"fmt"
	"github.com/lbrlabs/iac-in-go/pkg/cluster"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/eks"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		stack := ctx.Stack()
		org := ctx.Organization()

		// retrieve network details from network stack
		network, err := pulumi.NewStackReference(ctx, fmt.Sprintf("%v/go-network/%v", org, stack), nil)
		if err != nil {
			return err
		}

		publicSubnetIds := network.GetOutput(pulumi.String("publicSubnetIds"))
		privateSubnetIds := network.GetOutput(pulumi.String("privateSubnetIds"))

		nodeRole, err := cluster.NewNodeRole(ctx, "eks", &cluster.NodeRoleArgs{})
		if err != nil {
			return fmt.Errorf("error creating node role: %v", err)
		}

		theCluster, err := cluster.NewCluster(ctx, "eks", &cluster.ClusterArgs{
			ClusterSubnetIds:    pulumi.StringArrayOutput(publicSubnetIds),
			SystemNodeSubnetIds: pulumi.StringArrayOutput(privateSubnetIds),
			LetsEncryptEmail:    pulumi.String("mail@lbrlabs.com"),
			RoleMappings: []cluster.RoleMappingArgs{
				
			},
			NodeRoleArns: pulumi.StringArray{
				nodeRole.Role.Arn,
			},
		})
		if err != nil {
			return fmt.Errorf("error creating cluster: %v", err)
		}

		_, err = cluster.NewNodeGroup(ctx, "eks", &cluster.NodeGroupArgs{
			ClusterName: theCluster.ControlPlane.Name,
			Role:        nodeRole.Role,
			SubnetIds:   pulumi.StringArrayOutput(privateSubnetIds),
			ScalingConfig: eks.NodeGroupScalingConfigArgs{
				DesiredSize: pulumi.Int(1),
				MaxSize:     pulumi.Int(10),
				MinSize:     pulumi.Int(1),
			},
		})
		if err != nil {
			return fmt.Errorf("error creating node group: %v", err)
		}

		ctx.Export("clusterName", theCluster.ControlPlane.Name)
		ctx.Export("kubeconfig", theCluster.KubeConfig)

		return nil
	})
}
