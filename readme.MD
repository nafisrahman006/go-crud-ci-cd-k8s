<!-- ## 🚀 Go CRUD App
A secure, lightweight, and containerized CRUD application built with Go and uses postgres as DB — equipped with CI/CD pipelines, Docker secrets, and GitHub Actions for smooth development and deployment.
## 📁 Project Structure

```
GO-CRUD/
├── .github/               # GitHub workflows for CI/CD
├── .env                   # Environment variables (local use)
|── Screenshots            # Screenshots
|── user-api               # HELM Charts
├── .gitignore             # Git ignored files
├── docker-compose.yml     # Multi-service orchestration
├── Dockerfile             # Image build configuration
├── go.mod                 # Go module file
├── go.sum                 # Go dependency checksums
├── main.go                # Main application entry point
|── readme.MD              # Project documentation
|── deployall.yaml         # K8s yaml file

```

## 🔧 Configuration
Use .env for local development:
Changed the hardcoded `DB_USER` & `DB_PASSWORD` so far we can pass the from `.env` or from `terminal`.




## 🛠️ CI/CD Pipelines (GitHub Actions)

✅ Lint Check

Ensures code adheres to Go best practices.

🔍 Docker Security & Vulnerability Scan Using Trivy

```
Report Summary
┌────────────────────────────┬──────────┬─────────────────┬─────────┐
│           Target           │   Type   │ Vulnerabilities │ Secrets │
├────────────────────────────┼──────────┼─────────────────┼─────────┤
│ test-image (alpine 3.21.3) │  alpine  │        0        │    -    │
│ app/api                    │ gobinary │        0        │    -    │
└────────────────────────────┴──────────┴─────────────────┴─────────┘
```



🚀 Release Pipeline
Automatically builds and tags Docker images from GitHub releases.

Uses tag-based versioning like:
```
IMAGE_TAG=${GITHUB_REF#refs/tags/}
```
## Deploy On kubernetes

This project is deployed on a Kubernetes cluster. The DB_URL was updated for Kubernetes deployment, and NodePort was used to expose the application for testing.


Apply the Kubernetes services with the following command:
```
kubectl apply -f deployall.yaml
```
![](screenshots/checkallsvc.png)
Access the PostgreSQL pod and create a test table:
```
kubectl exec -it deployment/postgres -- bash
```
After accessing the pod, create a table for testing:
```
psql -U myuser -d mydb
CREATE TABLE users (id SERIAL PRIMARY KEY, name VARCHAR(100), email VARCHAR(100));
```
![](screenshots/createdb-from-terminal.png)

After deploying, test the application’s CRUD operations using Postman. Use the following URL to test:
```
http://localhost:30080/users  # 30080 is the NodePort
```
![](screenshots/post.png)

![](screenshots/get.png)

You can also use TablePlus to visually inspect the data in the PostgreSQL table:

![](screenshots/table.png)

# Deployment with Helm Chart

This Helm chart defines the deployment of a PostgreSQL database using Kubernetes resources. Below is a summary of the key configurations:

## Deployment Overview
- **Deployment Name**: `postgres`
- **Replica Count**: 1 (single instance of PostgreSQL)
- **Container Image**: Configured dynamically using the Helm values file (`{{ .Values.postgres.image }}`).
- **Environment Variables**: Loaded from a Kubernetes Secret named `postgres-credentials`.
- **Port Configuration**: Exposes PostgreSQL on container port `5432`.
- **Persistent Storage**:
  - A Persistent Volume Claim (PVC) named `postgres-pvc` is used to store PostgreSQL data.
  - The data is mounted to the container at `/var/lib/postgresql/data`.
- **Helm installation**
```
helm install user-api ./user-api \
  --set postgres.credentials.user=  \
  --set postgres.credentials.password=  \
  --set postgres.credentials.db=
  ```

## Helm Chart Benefits
- **Dynamic Configuration**: The use of Helm templates allows for easy customization of the PostgreSQL image and other parameters.
- **Secure Secrets Management**: Sensitive credentials are managed securely using Kubernetes Secrets.
- **Persistent Data**: Ensures data durability with a Persistent Volume Claim.

This setup provides a scalable and configurable way to deploy PostgreSQL in a Kubernetes environment. -->


## 🚀 Go CRUD App
A secure, lightweight, and containerized CRUD application built with Go and uses postgres as DB — equipped with CI/CD pipelines, Docker secrets, and GitHub Actions for smooth development and deployment.
## 📁 Project Structure

```
GO-CRUD/
├── .github/               # GitHub workflows for CI/CD
├── .env                   # Environment variables (local use)
|── Screenshots            # Screenshots
|── user-api               # HELM Charts
├── argocd/                # ArgoCD Application definition
│   └── application.yaml
├── k8s/                   # GitOps-managed Kubernetes manifests (Vault-integrated)
│   └── base/
│       ├── manifest.yaml
│       └── kustomization.yaml
├── .gitignore             # Git ignored files
├── docker-compose.yml     # Multi-service orchestration
├── Dockerfile              # Image build configuration
├── go.mod                 # Go module file
├── go.sum                 # Go dependency checksums
├── main.go                # Main application entry point
|── readme.MD              # Project documentation
|── deployall.yaml         # K8s yaml file (legacy, pre-GitOps)

```

## 🔧 Configuration
Use .env for local development:
Changed the hardcoded `DB_USER` & `DB_PASSWORD` so far we can pass the from `.env` or from `terminal`.




## 🛠️ CI/CD Pipelines (GitHub Actions)

✅ Lint Check

Ensures code adheres to Go best practices.

🔍 Docker Security & Vulnerability Scan Using Trivy

```
Report Summary
┌────────────────────────────┬──────────┬─────────────────┬─────────┐
│           Target           │   Type   │ Vulnerabilities │ Secrets │
├────────────────────────────┼──────────┼─────────────────┼─────────┤
│ test-image (alpine 3.21.3) │  alpine  │        0        │    -    │
│ app/api                    │ gobinary │        0        │    -    │
└────────────────────────────┴──────────┴─────────────────┴─────────┘
```



🚀 Release Pipeline
Automatically builds and tags Docker images from GitHub releases.

Uses tag-based versioning like:
```
IMAGE_TAG=${GITHUB_REF#refs/tags/}
```
## Deploy On kubernetes

This project is deployed on a Kubernetes cluster. The DB_URL was updated for Kubernetes deployment, and NodePort was used to expose the application for testing.


Apply the Kubernetes services with the following command:
```
kubectl apply -f deployall.yaml
```
![](screenshots/checkallsvc.png)
Access the PostgreSQL pod and create a test table:
```
kubectl exec -it deployment/postgres -- bash
```
After accessing the pod, create a table for testing:
```
psql -U myuser -d mydb
CREATE TABLE users (id SERIAL PRIMARY KEY, name VARCHAR(100), email VARCHAR(100));
```
![](screenshots/createdb-from-terminal.png)

After deploying, test the application's CRUD operations using Postman. Use the following URL to test:
```
http://localhost:30080/users  # 30080 is the NodePort
```
![](screenshots/post.png)

![](screenshots/get.png)

You can also use TablePlus to visually inspect the data in the PostgreSQL table:

![](screenshots/table.png)

# Deployment with Helm Chart

This Helm chart defines the deployment of a PostgreSQL database using Kubernetes resources. Below is a summary of the key configurations:

## Deployment Overview
- **Deployment Name**: `postgres`
- **Replica Count**: 1 (single instance of PostgreSQL)
- **Container Image**: Configured dynamically using the Helm values file (`{{ .Values.postgres.image }}`).
- **Environment Variables**: Loaded from a Kubernetes Secret named `postgres-credentials`.
- **Port Configuration**: Exposes PostgreSQL on container port `5432`.
- **Persistent Storage**:
  - A Persistent Volume Claim (PVC) named `postgres-pvc` is used to store PostgreSQL data.
  - The data is mounted to the container at `/var/lib/postgresql/data`.
- **Helm installation**
```
helm install user-api ./user-api \
  --set postgres.credentials.user=  \
  --set postgres.credentials.password=  \
  --set postgres.credentials.db=
  ```

## Helm Chart Benefits
- **Dynamic Configuration**: The use of Helm templates allows for easy customization of the PostgreSQL image and other parameters.
- **Secure Secrets Management**: Sensitive credentials are managed securely using Kubernetes Secrets.
- **Persistent Data**: Ensures data durability with a Persistent Volume Claim.

This setup provides a scalable and configurable way to deploy PostgreSQL in a Kubernetes environment.

# 🔄 GitOps Deployment with ArgoCD

To move beyond manual `kubectl apply`, this project also implements a full GitOps workflow using **ArgoCD** + **HashiCorp Vault**, so the cluster state is continuously reconciled against this Git repository rather than applied by hand.

## Why GitOps here
- **Git as the single source of truth** — the `k8s/base/` manifests define the desired cluster state; ArgoCD continuously syncs the live cluster to match
- **Self-healing** — any manual drift (e.g. `kubectl scale`, `kubectl edit`) is automatically detected and reverted back to what's defined in Git
- **Zero secrets in Git** — credentials are no longer stored as base64 K8s Secrets in manifests. Instead, a **Vault Agent Injector** sidecar pulls credentials from HashiCorp Vault at pod startup and injects them via `/vault/secrets/config`

## Architecture

```
┌──────────────┐      watches      ┌─────────────┐      applies      ┌───────────────┐
│  This Repo    │ ───────────────▶ │   ArgoCD    │ ─────────────────▶│  K8s Cluster  │
│  k8s/base/    │                   │             │                    │               │
└──────────────┘                   └─────────────┘                    └───────┬───────┘
                                                                                │ Vault Agent
                                                                                ▼ injects secrets
                                                                        ┌───────────────┐
                                                                        │  HashiCorp    │
                                                                        │    Vault      │
                                                                        └───────────────┘
```

## Setup

**1. Vault** (Kubernetes auth method, KV secrets engine, scoped policy):
```bash
helm install vault hashicorp/vault \
  --set "server.dev.enabled=true" \
  --set "injector.enabled=true" \
  -n vault --create-namespace
```

**2. ArgoCD**:
```bash
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

**3. Register the app with ArgoCD** — see [`argocd/application.yaml`](./argocd/application.yaml):
```bash
kubectl apply -f argocd/application.yaml
```

ArgoCD then takes over — pulling manifests from [`k8s/base/`](./k8s/base), applying them, and continuously watching for drift.

![ArgoCD application synced and healthy](screenshots/argo.png)

## Drift correction in action

To demonstrate the self-healing behavior, the `go-crud-app` deployment was manually scaled outside of Git:
```bash
kubectl scale deployment go-crud-app --replicas=4
```

ArgoCD detected the deviation from the Git-defined state (`replicas: 1`) and automatically reverted it — no manual intervention required.

