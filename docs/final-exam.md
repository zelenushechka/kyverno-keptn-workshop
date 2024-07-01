# Final Exam

Great! You made it to the final exam, where you will apply all the knowledge you have gained throughout the course.

You should have running `v3` of our Application in Production, which has a little special feature. It is using a Feature Flag to enable a new functionality.

For Feature Flagging we are using [OpenFeature](https://openfeature.dev/) and it's Evaluation Engine [FlagD](https://flagd.dev/).

Unfortunately, the feature is not working as expected, and slows down the already deployed application.

Your task is to create three Kyverno Policies to automatically fix the issue.

1. The first policy should create a load test job as soon as a flag is changed in the flagdefinition custom resource.
  - You can find the Job in `/excercises/final-exam/load-test-job.yaml`
  - You need to watch the API Group Kind `FeatureFlag` for an `UPDATE` operation.
  - Use label selectors to explicitly watch the flagdefinition for the `demo-app`
    ```
    app: sample-app
    type: feature-flag
    ```

2. The second policy should check if the load test job is completed and trigger a KeptnAnalysis.
  - You can find a KeptnAnalysis in `/excercises/final-exam/keptn-analysis.yaml`.
  - You need to watch the API Group Kind `Job/status` for an `UPDATE` operation.
  - Use a precondition to check if the Job is completed `request.object.status.succeeded`


3. The third policy should check the result of the KeptnAnalysis and roll back the flagdefinition to it's previous state. 
  - You can find the Job in `/excercises/final-exam/argo-rollback.yaml`.
  - You need to watch the API Group Kind `Analysis/status` for an `UPDATE` operation.
  - Use a precondition to check if the Analysis is completed `request.object.status.state == "Completed"` and the roll back if  `request.object.status.pass != "pass"` be aware that this field does not exist if the Analysis is not completed or failed.

## Tips and Lessons Learned 

### Permissions for Kyverno

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

### Permissions for ArgoCD

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

### Argo Sync for Generate Policies

Generate Policies are immutable. To force ArgoCD to recreate the policies during the Sync process, apply the following annotations to each of your policy custom resource:

```yaml
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  annotations:
    argocd.argoproj.io/sync-options: Force=true,Replace=true
```

### Resource Naming

Make sure to append this Kyverno Variable to each of your generated resources to prevent naming conflicts:

```
{% raw %}resource-name-{{request.object.metadata.resourceVersion}}{% endraw %}
```

### Argo Sync and Rollback

To disable autosync in argo for a specific application, you can't do this via the CLI, instead you need to patch the application resource:

```bash
kubectl -n argocd patch --type='merge' application kyverno-keptn-workshop -p "{\"spec\":{\"syncPolicy\":null}}"
```

Please see `/src/agro-cli/rollout.sh` for an example.

### Kyverno Policies with Preconditions and ArgoCD

Use within your precondition check instead of Quotes or Double Quotes a Backtick to prevent issues with the YAML Parser.

