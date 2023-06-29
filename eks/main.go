package main

import (
	"fmt"
	"github.com/lbrlabs/iac-in-go/pkg/cluster"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
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
		
		cluster, err := cluster.NewCluster(ctx, "eks", &cluster.ClusterArgs{
			ClusterSubnetIds:    pulumi.StringArrayOutput(publicSubnetIds),
			SystemNodeSubnetIds: pulumi.StringArrayOutput(privateSubnetIds),
			LetsEncryptEmail:       pulumi.String("mail@lbrlabs.com"),
		})
		if err != nil {
			return fmt.Errorf("error creating cluster: %v", err)
		}

		ctx.Export("clusterName", cluster.ControlPlane.Name)
		ctx.Export("kubeconfig", cluster.KubeConfig)

		return nil
	})
}
