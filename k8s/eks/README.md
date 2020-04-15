# Provision Kubernetes Cluster on AWS

It's pretty simple to provisin EKS cluster with simple command, but have to install eksctl, [here's how](https://eksctl.io/introduction/installation/).


```bash

eksctl create cluster --version=1.15

```

Or go with my recommended configration.

```bash

eksctl create cluster -f cluster-nodes.yaml

```