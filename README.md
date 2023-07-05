# iac-in-go

A reference architecture example for Pulumi in Go

## Usage

This repo contains multiple go modules, so a [go workspace](https://go.dev/doc/tutorial/workspaces) should be created to use it effectively:

```bash
go work init
go work use eks
go work use network
go work use pkg/cluster
go work use pkg/irsa
go work use pkg/cluster/kubeconfig
```