# Final Exam

Great! You made it to the final exam, where you will apply all the knowledge you have gained throughout the course.

You should have running `v3` of our Application in Production, which has a little special feature. It is using a Feature Flag to enable a new functionality.

For Feature Flagging we are using OpenFeature and it's Evaluation Engine FlagD.

Unfortunately, the feature is not working as expected, and slows down the already deployed application.

Your task is to create three Kyverno Policies to automatically fix the issue.

1. The first policy should create a load test job as soon as a flag is changed in the flagdefinition custom resource.
2. The second policy should check if the load test job is completed and trigger a KeptnAnalysis.
3. The third policy should check the result of the KeptnAnalysis and roll back the flagdefinition to it's previous state. To do this, we have already prepared a Job which you can find in the root of the repository.

## Tips and Tricks

## Permissions for Kyverno

To enable the Kyverno Background Controller to watch and generate resources that are needed within this exam, you have to add addional ClusterRoles to the Kyverno Controller.

This is already applied in your cluster

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  labels:
    app.kubernetes.io/component: background-controller
    app.kubernetes.io/instance: kyverno
    app.kubernetes.io/part-of: kyverno
  name: kyverno:create-keptn-analysis
rules:
- apiGroups:
  - metrics.keptn.sh
  - batch
  resources:
  - 'analyses'
  - 'jobs'
  verbs:
  - create
  - list
  - get
  - watch
  - update
  - delete
```

## Permissions for ArgoCD

One of the Policies you will create needs the ability to patch ArgoCD Application resources. This is already applied in your cluster.

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: argocd-rollback
  namespace: argocd
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: argocd-rollback-role
rules:
  - apiGroups: ["argoproj.io"]
    resources: ["*"]
    verbs: ["patch", "get"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: argocd-rollback-crb
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: argocd-rollback-role
subjects:
  - kind: ServiceAccount
    name: argocd-rollback
    namespace: argocd
```

## Argo Sync for Generate Policies

Generate Policies are immutable. To force ArgoCD to recreate the policies during the Sync process, apply the following annotations to each of your policy custom resource:

```yaml
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  annotations:
    argocd.argoproj.io/sync-options: Force=true,Replace=true
```

## Resource Naming

Make sure to append this Kyverno Variable to each of your generated resources to prevent naming conflicts:

```
resource-name-{{request.object.metadata.resourceVersion}}
```