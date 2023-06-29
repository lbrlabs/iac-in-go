# Network

The network project bootstraps the foundational network in AWS that's required for all other projects.

We create a VPC with public and private subnets that meet the basic needs for running workloads.

## Configuration

The project within this repo uses Pulumi [configuration](https://www.pulumi.com/docs/concepts/config/) to make changes to the project graph depending on inputs from the configuration.

This is a common practice when using Pulumi stacks across different cloud environments, such as for development/staging/production environments and across cloud provider regions.