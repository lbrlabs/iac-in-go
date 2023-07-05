package cluster

import (
	"encoding/json"
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type NodeRoleArgs struct {
}

type NodeRole struct {
	pulumi.ResourceState
	Role *iam.Role `pulumi:"nodeRole"`
}

// NewNodeRole creates a new EKS Node group component resource.
func NewNodeRole(ctx *pulumi.Context,
	name string, args *NodeRoleArgs, opts ...pulumi.ResourceOption) (*NodeRole, error) {

	if args == nil {
		args = &NodeRoleArgs{}
	}

	component := &NodeRole{}
	err := ctx.RegisterComponentResource("iac-in-go:index:NodeRole", name, component, opts...)
	if err != nil {
		return nil, err
	}

	nodePolicyJSON, err := json.Marshal(map[string]interface{}{
		"Statement": []map[string]interface{}{
			{
				"Action": "sts:AssumeRole",
				"Effect": "Allow",
				"Principal": map[string]interface{}{
					"Service": "ec2.amazonaws.com",
				},
			},
		},
		"Version": "2012-10-17",
	})
	if err != nil {
		return nil, fmt.Errorf("error marshalling node policy: %w", err)
	}

	nodeRole, err := iam.NewRole(ctx, fmt.Sprintf("%s-node-role", name), &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(nodePolicyJSON),
	}, pulumi.Parent(component))
	if err != nil {
		return nil, fmt.Errorf("error creating node role: %w", err)
	}

	_, err = iam.NewRolePolicyAttachment(ctx, fmt.Sprintf("%s-node-worker-policy", name), &iam.RolePolicyAttachmentArgs{
		PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"),
		Role:      nodeRole.Name,
	}, pulumi.Parent(nodeRole))
	if err != nil {
		return nil, fmt.Errorf("error attaching node worker policy: %w", err)
	}

	_, err = iam.NewRolePolicyAttachment(ctx, fmt.Sprintf("%s-node-ecr-policy", name), &iam.RolePolicyAttachmentArgs{
		PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"),
		Role:      nodeRole.Name,
	}, pulumi.Parent(nodeRole))
	if err != nil {
		return nil, fmt.Errorf("error attaching system node ecr policy: %w", err)
	}

	component.Role = nodeRole

	if err := ctx.RegisterResourceOutputs(component, pulumi.Map{
		"nodeRole": nodeRole,

	}); err != nil {
		return nil, err
	}

	return component, nil
}
