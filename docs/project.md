# Go Modules

#### Introduction

- Dependencies
- Module Mirrors
- Checksum
- Vendoring
- MVS Algorithm

#### Dependencies

Go modules manages third party projects for your project.

go clean modcache // to clean module cache

#### Module Mirrors

To install packages from local mod cache thats has already been installed, reason 
to eliminate going to internet. Perform below:

go env -w GOPROXY='file:///home/bersen/go/pkg/mod/,https://proxy.golang.org,direct'

then

go get [pkg]@v???

#### Checksum Database

ensures that the code beening download for a specific version is correct, works with go.mod
and go.sum.

#### Vendoring

go mod vendor localizes third party packages to the module project for better management of code.

# Kubernetes

#### Tooling

- Install k8s, kubectl, docker and kustomize

#### Clusters, Nodes and Pods

Clusters are abstractions over physical or virtual hardware providing compute resources via
nodes which run pods to run containers to run applications.

#### Write basic service

Prototype a basic service
