#!/bin/bash
# Sync RBAC files from config/rbac to Helm chart templates
# This script converts the generated and manually maintained RBAC files
# into Helm templates with proper templating

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
RBAC_DIR="${PROJECT_ROOT}/config/rbac"
HELM_RBAC="${PROJECT_ROOT}/charts/langfuse-controller-helm/templates/rbac.yaml"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Syncing RBAC to Helm chart...${NC}"

# Check if required files exist
if [[ ! -f "${RBAC_DIR}/role.yaml" ]]; then
    echo "Error: ${RBAC_DIR}/role.yaml not found. Run 'make manifests' first."
    exit 1
fi

if [[ ! -f "${RBAC_DIR}/leader_election_role.yaml" ]]; then
    echo "Error: ${RBAC_DIR}/leader_election_role.yaml not found."
    exit 1
fi

# Check if PyYAML is available
if ! python3 -c "import yaml" 2>/dev/null; then
    echo -e "${YELLOW}Warning: PyYAML not found. Installing...${NC}"
    pip3 install --user pyyaml || {
        echo "Error: Could not install PyYAML. Please install it manually: pip3 install pyyaml"
        exit 1
    }
fi

# Export environment variables for Python script
export RBAC_DIR="${RBAC_DIR}"
export HELM_RBAC="${HELM_RBAC}"

# Run the Python script to generate Helm RBAC template
python3 << PYTHON_SCRIPT
import yaml
import sys
import os

rbac_dir = os.environ['RBAC_DIR']
helm_rbac = os.environ['HELM_RBAC']

# Read role.yaml (ClusterRole)
with open(f'{rbac_dir}/role.yaml', 'r') as f:
    role_data = yaml.safe_load(f)

# Read leader_election_role.yaml (Role)
with open(f'{rbac_dir}/leader_election_role.yaml', 'r') as f:
    leader_election_data = yaml.safe_load(f)

# Generate Helm template
output = []

# ClusterRole
output.append("apiVersion: rbac.authorization.k8s.io/v1")
output.append("kind: ClusterRole")
output.append("metadata:")
output.append("  name: {{ include \"langfuse-controller-helm.fullname\" . }}")
output.append("  labels:")
output.append("    {{- include \"langfuse-controller-helm.labels\" . | nindent 4 }}")
output.append("rules:")

# Add rules from role.yaml with proper indentation
for rule in role_data.get('rules', []):
    rule_yaml = yaml.dump([rule], default_flow_style=False, sort_keys=False).rstrip()
    # Indent each line by 2 spaces (rules are already at root level)
    indented_lines = []
    for line in rule_yaml.split('\n'):
        if line.strip():  # Skip empty lines
            indented_lines.append('  ' + line)
        else:
            indented_lines.append('')
    output.append('\n'.join(indented_lines))

# Add secrets permissions (Helm-specific)
output.append("- apiGroups:")
output.append("  - \"\"")
output.append("  resources:")
output.append("  - secrets")
output.append("  verbs:")
output.append("  - create")
output.append("  - delete")
output.append("  - get")
output.append("  - list")
output.append("  - patch")
output.append("  - update")
output.append("  - watch")

# ClusterRoleBinding
output.append("---")
output.append("apiVersion: rbac.authorization.k8s.io/v1")
output.append("kind: ClusterRoleBinding")
output.append("metadata:")
output.append("  name: {{ include \"langfuse-controller-helm.fullname\" . }}")
output.append("  labels:")
output.append("    {{- include \"langfuse-controller-helm.labels\" . | nindent 4 }}")
output.append("roleRef:")
output.append("  apiGroup: rbac.authorization.k8s.io")
output.append("  kind: ClusterRole")
output.append("  name: {{ include \"langfuse-controller-helm.fullname\" . }}")
output.append("subjects:")
output.append("- kind: ServiceAccount")
output.append("  name: {{ include \"langfuse-controller-helm.serviceAccountName\" . }}")
output.append("  namespace: {{ .Release.Namespace }}")

# Leader Election Role
output.append("---")
output.append("# permissions to do leader election.")
output.append("apiVersion: rbac.authorization.k8s.io/v1")
output.append("kind: Role")
output.append("metadata:")
output.append("  name: {{ include \"langfuse-controller-helm.fullname\" . }}-leader-election")
output.append("  namespace: {{ .Release.Namespace }}")
output.append("  labels:")
output.append("    {{- include \"langfuse-controller-helm.labels\" . | nindent 4 }}")
output.append("rules:")

# Add leader election rules with proper indentation
for rule in leader_election_data.get('rules', []):
    rule_yaml = yaml.dump([rule], default_flow_style=False, sort_keys=False).rstrip()
    # Indent each line by 2 spaces
    indented_lines = []
    for line in rule_yaml.split('\n'):
        if line.strip():  # Skip empty lines
            indented_lines.append('  ' + line)
        else:
            indented_lines.append('')
    output.append('\n'.join(indented_lines))

# Leader Election RoleBinding
output.append("---")
output.append("apiVersion: rbac.authorization.k8s.io/v1")
output.append("kind: RoleBinding")
output.append("metadata:")
output.append("  name: {{ include \"langfuse-controller-helm.fullname\" . }}-leader-election")
output.append("  namespace: {{ .Release.Namespace }}")
output.append("  labels:")
output.append("    {{- include \"langfuse-controller-helm.labels\" . | nindent 4 }}")
output.append("roleRef:")
output.append("  apiGroup: rbac.authorization.k8s.io")
output.append("  kind: Role")
output.append("  name: {{ include \"langfuse-controller-helm.fullname\" . }}-leader-election")
output.append("subjects:")
output.append("- kind: ServiceAccount")
output.append("  name: {{ include \"langfuse-controller-helm.serviceAccountName\" . }}")
output.append("  namespace: {{ .Release.Namespace }}")

# Write output
with open(helm_rbac, 'w') as f:
    f.write('\n'.join(output) + '\n')

print("✅ RBAC synced successfully")
PYTHON_SCRIPT

echo -e "${GREEN}✅ RBAC synced successfully to ${HELM_RBAC}${NC}"
