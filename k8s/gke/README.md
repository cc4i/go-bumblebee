# Provision Kubernetes Cluster on GCP

It's simple as googling, but have to ensure that you have enabled the Google Kubernetes Engine API and install [Cloud SDK](https://cloud.google.com/sdk/install).

```bash

gcloud container clusters create go-bumblebee-cluster\
    --zone asia-east1-a

```

where:

- 'go-bumblebee-cluster' - cluster-name is the name of your new cluster.
- 'asia-east1-a' - compute-zone is the compute zone for your cluster.
