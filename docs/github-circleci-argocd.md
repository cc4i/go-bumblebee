# GitHub + CircleCI + AgroCD



### Preparation

- Install Argo CD

```bash
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```


# checkout password, admin user -> 'admin'
kubectl get pods -n argocd -l app.kubernetes.io/name=argocd-server -o name | cut -d'/' -f 2

# Optional for install Argo CD CLI (Mac)
brew tap argoproj/tap
brew install argoproj/tap/argocd

```


- Install Argo Rollouts

```bash
kubectl create namespace argo-rollouts
kubectl apply -n argo-rollouts -f https://raw.githubusercontent.com/argoproj/argo-rollouts/stable/manifests/install.yaml
```


## Yaml

## Kustomized yaml

## Helm

## Blue/Green deployment

## Canary deployment