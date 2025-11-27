# Langfuse Kubernetes Controller

A Kubernetes operator for managing Langfuse resources via Custom Resource Definitions (CRDs).

[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/langfuse-controller-helm)](https://artifacthub.io/packages/search?repo=langfuse-controller-helm)

## Features

- **Declarative Management**: Manage Langfuse resources using Kubernetes manifests
- **GitOps Ready**: Full support for GitOps workflows (ArgoCD, Flux, etc.)
- **Automatic Secret Management**: API keys automatically stored in Kubernetes Secrets
- **Helm Deployment**: Production-ready Helm chart included

## Supported Resources

- `LangfuseProject` - Langfuse projects
- `LangfuseAPIKey` - Project API keys (with Secret creation)
- `LangfuseModel` - Model definitions and pricing
- `LangfuseLlmConnection` - LLM provider connections
- `LangfusePrompt` - Prompt templates
- `LangfuseScoreConfig` - Score configurations

## Quick Start

### Prerequisites

- Kubernetes 1.28+
- Helm 3.x
- Langfuse Public and Secret API keys

### Installation

```bash
# Install via Helm
helm install langfuse-controller ./charts/langfuse-controller-helm \\
  --set langfuse.host="https://cloud.langfuse.com" \\
  --set langfuse.publicKey="pk-..." \\
  --set langfuse.secretKey="sk-..."
```

### Example Usage

```yaml
# Create a project
apiVersion: langfuse.io/v1alpha1
kind: LangfuseProject
metadata:
  name: my-project
spec:
  name: "My Project"
---
# Create an API key
apiVersion: langfuse.io/v1alpha1
kind: LangfuseAPIKey
metadata:
  name: my-api-key
spec:
  projectRef: my-project
  name: "production-key"
  secretName: langfuse-credentials
```

Apply the manifest:
```bash
kubectl apply -f project.yaml
```

The controller will:
1. Create the project in Langfuse
2. Generate an API key
3. Store the credentials in a Kubernetes Secret named `langfuse-credentials`

## Development

### Build

```bash
make build
```

### Run locally

```bash
export LANGFUSE_HOST="https://cloud.langfuse.com"
export LANGFUSE_PUBLIC_KEY="pk-..."
export LANGFUSE_SECRET_KEY="sk-..."
make run
```

### Generate manifests

```bash
make manifests
```

## Configuration

The controller is configured via environment variables:

- `LANGFUSE_HOST` - Langfuse API endpoint (default: `https://cloud.langfuse.com`)
- `LANGFUSE_PUBLIC_KEY` - Langfuse Public API key for authentication
- `LANGFUSE_SECRET_KEY` - Langfuse Secret API key for authentication

## Architecture

The controller uses:
- **kubebuilder** for scaffolding
- **controller-runtime** for Kubernetes integration
- **Helm** for deployment

## License

Apache License 2.0
