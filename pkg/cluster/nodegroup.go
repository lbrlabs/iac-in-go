package cluster

import (
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/eks"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type NodeGroupArgs struct {
	ClusterName      pulumi.StringInput             `pulumi:"clusterName"`
	SubnetIds        pulumi.StringArrayInput        `pulumi:"subnetIds"`
	InstanceTypes    *pulumi.StringArrayInput       `pulumi:"instanceTypes"`
	NodeMaxCount     *pulumi.IntInput               `pulumi:"nodeMaxCount"`
	NodeMinCount     *pulumi.IntInput               `pulumi:"nodeMinCount"`
	NodeDesiredCount *pulumi.IntInput               `pulumi:"nodeDesiredCount"`
	Taints           eks.NodeGroupTaintArray        `pulumi:"taints"`
	Labels           pulumi.StringMapInput          `pulumi:"labels"`
	ScalingConfig    eks.NodeGroupScalingConfigArgs `pulumi:"scalingConfig"`
	Role             *iam.Role                      `pulumi:"nodeRole"`
}

type NodeGroup struct {
	pulumi.ResourceState

	NodeGroup *eks.NodeGroup `pulumi:"attachedNodeGroup"`
}

// NewNodeGroup creates a new EKS Node group component resource.
func NewNodeGroup(ctx *pulumi.Context,
	name string, args *NodeGroupArgs, opts ...pulumi.ResourceOption) (*NodeGroup, error) {
	if args == nil {
		args = &NodeGroupArgs{}
	}

	component := &NodeGroup{}
	err := ctx.RegisterComponentResource("iac-in-go:index:AttachedNodeGroup", name, component, opts...)
	if err != nil {
		return nil, err
	}

	var instanceTypes pulumi.StringArrayInput

	if args.InstanceTypes == nil {
		instanceTypes = pulumi.StringArray{
			pulumi.String("t3.medium"),
		}
	} else {
		instanceTypes = *args.InstanceTypes
	}

	

	nodeGroup, err := eks.NewNodeGroup(ctx, fmt.Sprintf("%s-nodes", name), &eks.NodeGroupArgs{
		ClusterName:   args.ClusterName,
		SubnetIds:     args.SubnetIds,
		NodeRoleArn:   args.Role.Arn,
		Taints:        args.Taints,
		InstanceTypes: instanceTypes,
		Labels:        args.Labels,
		ScalingConfig: args.ScalingConfig,
	}, pulumi.Parent(component), pulumi.IgnoreChanges([]string{"scalingConfig"}))
	if err != nil {
		return nil, fmt.Errorf("error creating system nodegroup provider: %w", err)
	}

	component.NodeGroup = nodeGroup

	if err := ctx.RegisterResourceOutputs(component, pulumi.Map{
		"nodeGroup": nodeGroup,
	}); err != nil {
		return nil, err
	}

	return component, nil

}
