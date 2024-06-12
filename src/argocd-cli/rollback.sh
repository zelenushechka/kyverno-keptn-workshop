#!/bin/bash

# Check if project name is provided as parameter
if [ -z "$1" ]; then
    echo "Please provide a project name as a parameter."
    exit 1
fi

# Read secret from /etc/argosecret
secret=$(cat /etc/argosecret/password)

# Login to ArgoCD server
argocd login argocd-server.argocd.svc.cluster.local --insecure --username admin --password "$secret"

# Rollback application in the specified project
argocd app rollback "$1"