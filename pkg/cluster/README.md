# Cluster

This directory contains a [Pulumi Component](https://www.pulumi.com/docs/concepts/resources/components/) that defines an EKS cluster.

It's designed to be reusable across the Mono Repo where required.

In larger organizations, this would likely be a distinct, separate git repository that is versioned and imported distinctly. In this reference example, we keep the component within the repo to ease the resolution of packages for the user.

This component is lifted verbatim from a Pulumi Multi-Language package that lives [here](github.com/lbrlabs/pulumi-lbrlabs-eks). It has been copied to this repo to showcase the usage of components within a mono repo.

## Todo

Currently, these components are not unit or integration tested, which would be following best practices for components