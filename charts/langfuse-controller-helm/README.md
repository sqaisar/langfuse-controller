# langfuse-controller-helm

[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/langfuse-controller-helm)](https://artifacthub.io/packages/search?repo=langfuse-controller-helm)

A Helm chart for deploying the Langfuse Kubernetes Controller, which manages Langfuse resources via Custom Resource Definitions (CRDs).

## Introduction

This chart deploys the Langfuse Kubernetes Controller, enabling you to manage Langfuse projects, API keys, models, LLM connections, prompts, and score configurations declaratively using Kubernetes manifests.

## Features

- **Declarative Management**: Manage Langfuse resources using Kubernetes manifests
- **GitOps Ready**: Full support for GitOps workflows (ArgoCD, Flux, etc.)
- **Automatic Secret Management**: API keys automatically stored in Kubernetes Secrets
- **Production Ready**: Includes RBAC, service accounts, and resource limits

## Prerequisites

- Kubernetes 1.28+
- Helm 3.x
- Langfuse Public and Secret API keys

## Installation

### Add the Helm repository

```bash
helm repo add langfuse-controller https://sqaisar.github.io/langfuse-controller
helm repo update
```

### Install the chart

```bash
helm install langfuse-controller langfuse-controller/langfuse-controller-helm \
  --set langfuse.host="https://cloud.langfuse.com" \
  --set langfuse.publicKey="pk-..." \
  --set langfuse.secretKey="sk-..." \
  --namespace langfuse-controller \
  --create-namespace
```

### Install using existing secret

If you prefer to store the API keys in a Kubernetes Secret:

```bash
# Create the secret
kubectl create secret generic langfuse-api-keys \
  --from-literal=LANGFUSE_PUBLIC_KEY="pk-..." \
  --from-literal=LANGFUSE_SECRET_KEY="sk-..." \
  --namespace langfuse-controller

# Install with existing secret
helm install langfuse-controller langfuse-controller/langfuse-controller-helm \
  --set langfuse.host="https://cloud.langfuse.com" \
  --set langfuse.existingSecret="langfuse-api-keys" \
  --namespace langfuse-controller \
  --create-namespace
```

## Configuration

The following table lists the configurable parameters and their default values:

| Parameter | Description | Default |
|-----------|-------------|---------|
| `image.repository` | Container image repository | `sqaisar/langfuse-controller` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `image.tag` | Container image tag | `latest` |
| `imagePullSecrets` | Image pull secrets | `[]` |
| `nameOverride` | Override the name of the chart | `""` |
| `fullnameOverride` | Override the full name of the chart | `langfuse-controller` |
| `watchNamespaces` | List of namespaces to watch (empty = all namespaces) | `[]` |
| `serviceAccount.create` | Create a service account | `true` |
| `serviceAccount.annotations` | Service account annotations | `{}` |
| `serviceAccount.name` | Service account name | `""` |
| `podAnnotations` | Pod annotations | `{}` |
| `podSecurityContext` | Pod security context | `{}` |
| `securityContext` | Container security context | `{}` |
| `resources.limits.cpu` | CPU limit | `500m` |
| `resources.limits.memory` | Memory limit | `128Mi` |
| `resources.requests.cpu` | CPU request | `10m` |
| `resources.requests.memory` | Memory request | `64Mi` |
| `nodeSelector` | Node selector | `{}` |
| `tolerations` | Tolerations | `[]` |
| `affinity` | Affinity rules | `{}` |
| `langfuse.host` | Langfuse API endpoint | `https://cloud.langfuse.com` |
| `langfuse.publicKey` | Langfuse Public API key | `""` |
| `langfuse.secretKey` | Langfuse Secret API key | `""` |
| `langfuse.existingSecret` | Name of existing secret with `LANGFUSE_PUBLIC_KEY` and `LANGFUSE_SECRET_KEY` | `""` |

## Usage Examples

### Basic Installation

```bash
helm install langfuse-controller langfuse-controller/langfuse-controller-helm \
  --set langfuse.host="https://cloud.langfuse.com" \
  --set langfuse.publicKey="pk-..." \
  --set langfuse.secretKey="sk-..."
```

### Custom Image Tag

```bash
helm install langfuse-controller langfuse-controller/langfuse-controller-helm \
  --set image.tag="v0.1.0" \
  --set langfuse.host="https://cloud.langfuse.com" \
  --set langfuse.publicKey="pk-..." \
  --set langfuse.secretKey="sk-..."
```

### Watch Specific Namespaces

```bash
helm install langfuse-controller langfuse-controller/langfuse-controller-helm \
  --set watchNamespaces[0]="default" \
  --set watchNamespaces[1]="production" \
  --set langfuse.host="https://cloud.langfuse.com" \
  --set langfuse.publicKey="pk-..." \
  --set langfuse.secretKey="sk-..."
```

### Custom Resource Limits

```bash
helm install langfuse-controller langfuse-controller/langfuse-controller-helm \
  --set resources.limits.cpu="1000m" \
  --set resources.limits.memory="256Mi" \
  --set resources.requests.cpu="100m" \
  --set resources.requests.memory="128Mi" \
  --set langfuse.host="https://cloud.langfuse.com" \
  --set langfuse.publicKey="pk-..." \
  --set langfuse.secretKey="sk-..."
```

## Managing Langfuse Resources

After installing the controller, you can create Langfuse resources using Kubernetes manifests:

### Create a Langfuse Project

```yaml
apiVersion: langfuse.io/v1alpha1
kind: LangfuseProject
metadata:
  name: my-project
spec:
  name: "My Project"
```

### Create an API Key

```yaml
apiVersion: langfuse.io/v1alpha1
kind: LangfuseAPIKey
metadata:
  name: my-api-key
spec:
  projectRef: my-project
  name: "production-key"
  secretName: langfuse-credentials
```

The controller will automatically:
1. Create the project in Langfuse
2. Generate an API key
3. Store the credentials in a Kubernetes Secret named `langfuse-credentials`

## Supported CRDs

- `LangfuseProject` - Langfuse projects
- `LangfuseAPIKey` - Project API keys (with Secret creation)
- `LangfuseModel` - Model definitions and pricing
- `LangfuseLlmConnection` - LLM provider connections
- `LangfusePrompt` - Prompt templates
- `LangfuseScoreConfig` - Score configurations

## Upgrading

```bash
helm upgrade langfuse-controller langfuse-controller/langfuse-controller-helm \
  --set langfuse.host="https://cloud.langfuse.com" \
  --set langfuse.publicKey="pk-..." \
  --set langfuse.secretKey="sk-..."
```

## Uninstallation

```bash
helm uninstall langfuse-controller --namespace langfuse-controller
```

## Troubleshooting

### Check controller logs

```bash
kubectl logs -n langfuse-controller deployment/langfuse-controller
```

### Verify CRDs are installed

```bash
kubectl get crds | grep langfuse.io
```

### Check controller status

```bash
kubectl get pods -n langfuse-controller
```

## License

Apache License 2.0

## Support

- [GitHub Repository](https://github.com/sqaisar/langfuse-controller)
- [Documentation](https://github.com/sqaisar/langfuse-controller/blob/main/README.md)

