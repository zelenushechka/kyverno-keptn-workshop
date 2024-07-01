#!/usr/bin/env bash

kind create cluster --name workshop-cluster --config infrastructure/kind/config.yaml

helm install cert-manager --namespace cert-manager --create-namespace --wait --repo https://charts.jetstack.io cert-manager --values - <<EOF
crds:
  enabled: true
EOF

helm install keptn --namespace keptn-system --create-namespace --wait --repo https://charts.lifecycle.keptn.sh keptn \
--version 0.6.0 --values - <<EOF
lifecycleOperator:
  schedulingGatesEnabled: true
  promotionTasksEnabled: true
EOF

helm install kyverno --namespace kyverno --create-namespace --wait --repo https://kyverno.github.io/kyverno kyverno --version 3.1.4
