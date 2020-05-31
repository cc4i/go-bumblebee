
# Go-Bumblebee ![CircleCI](https://circleci.com/gh/cc4i/go-bumblebee.svg?style=svg) ![GoReport](https://goreportcard.com/badge/github.com/cc4i/go-bumblebee)

<img src="./docs/anything.jpg" alt="bumblebee" width="80">

It's sample application with Microservices style on Kubernetes, developers can leverage that to experience various features on Kubernets, such DevOps tooling, observability components, service mesh, etc. 




## Quick Start

### High level architecture 

![ss](./docs/arch0.jpg)

## Interesting Branches
- [Go-Bumblebee-Jazz](https://github.com/cc4i/go-bumblebee/tree/new-combo-jazz)

## Tutorial

### 1. Go-Bumblebee through CI/CD pipleline
Build up CI/CD capabilities around Kubernetes, we have so many ecllecnt choices due to in incredibly rich ecosystem. This tutorial will go with free lunch to taste modernized release pipeline.

- [GitHub + CircleCI + AgroCD](./docs/github-circleci-argocd.md)
    - [yaml](./docs/github-circleci-argocd.md#yaml)
    - [kustomized yaml](./docs/github-circleci-argocd.md#kustomized-yaml)
    - [helm](./docs/github-circleci-argocd.md#helm)
    - [blue/green deployment](./docs/github-circleci-argocd.md#blue-green-deployment)
    - [canary deployment](./docs/github-circleci-argocd.md#canary-deployment)

- GitHub + CircleCI + Spinnaker


### 2. Manage traffic with Istio
Leverage Istio is a great idea for manage traffic, especially Istio 1.5 was evolved a lot. 

- Routing 
- Fault injection
- Traffic shifting 
- Circuit breaking
- Mirroring
- Ingress gateway



### 3. Improve Observability with Istio

### 4. Enhance security with Istio

### 5. Go-Bumblebee on Knative

