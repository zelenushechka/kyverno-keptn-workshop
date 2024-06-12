#!/bin/bash

# Check if project name is provided as parameter
if [ -z "$1" ]; then
    echo "Please provide a project name as a parameter."
    exit 1
fi

# Read secret from /etc/argosecret
secret=$(cat /etc/argosecret/password)

# Disable autosync
kubectl -n argocd patch --type='merge' application kyverno-keptn-workshop -p "{\"spec\":{\"syncPolicy\":null}}"
kubectl -n argocd patch --type='merge' application $1 -p "{\"spec\":{\"syncPolicy\":null}}"

# Login to ArgoCD server
argocd login argocd-server.argocd.svc.cluster.local --insecure --username admin --password "$secret"

# Rollback application in the specified project
argocd app rollback "$1"