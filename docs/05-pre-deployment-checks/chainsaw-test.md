# Pre-Deployment Checks (chainsaw test)

With pre-deployment checks enabled and the maintenance window check in place, a KeptnTask is now supposed to execute before the deployment is considered successful.

However, if this automation is not working properly it can put your system at risk.
To exercise the automation and make sure it behaves as expected you can write and execute a chainsaw test.

Full documentation of chainsaw can be found [here](https://kyverno.github.io/chainsaw/latest/).

## Exercise

In this exercise, you will create a chainsaw test to verify that a KeptnTask is run when the deployment is updated.

### Create the chainsaw test

Create a new file `chainsaw-test.yaml` in the `tests/05-pre-deployment-checks` folder of your repository and add the following content:

```yaml
apiVersion: chainsaw.kyverno.io/v1alpha1
kind: Test
metadata:
  name: maintenance-window-check
spec:
  bindings:
  - name: repo
    value: (env('GITHUB_REPOSITORY'))
  steps:
  - try:
    - apply:
        resource:
          apiVersion: argoproj.io/v1alpha1
          kind: Application
          metadata:
            name: demo-app-test
            namespace: argocd
            annotations:
              argocd.argoproj.io/sync-wave: '10'
            finalizers:
              - resources-finalizer.argocd.argoproj.io
          spec:
            project: default
            source:
              path: charts/demo-app
              repoURL: (join('/', ['https://github.com', $repo]))
              targetRevision: main
              helm:
                valuesObject:
                  repo:
                    name: ($repo)
                  keptn:
                    appContext:
                      preDeploymentTasks:
                      - maintenance-window-check
                valueFiles:
                  - ../../gitops/dev/demo-app/values.yaml
                  - ../../gitops/dev/demo-app/values-specific.yaml
                parameters:
                  - name: commitID
                    value: $ARGOCD_APP_REVISION
                  - name: serviceVersion
                    value: v1
                  - name: service.nodePort
                    value: '31106'
            destination:
              server: https://kubernetes.default.svc
              namespace: ($namespace)
            syncPolicy:
              automated:
                prune: true
                selfHeal: true
    - assert:
        timeout: 10m
        resource:
          apiVersion: argoproj.io/v1alpha1
          kind: Application
          metadata:
            name: demo-app-test
            namespace: argocd
          status:
            health:
              status: Healthy
            sync:
              status: Synced
    - assert:
        timeout: 1m
        resource:
          apiVersion: lifecycle.keptn.sh/v1
          kind: KeptnApp
          metadata:
            name: demo-app
    - apply:
        resource:
          apiVersion: argoproj.io/v1alpha1
          kind: Application
          metadata:
            name: demo-app-test
            namespace: argocd
            annotations:
              argocd.argoproj.io/sync-wave: '10'
            finalizers:
              - resources-finalizer.argocd.argoproj.io
          spec:
            project: default
            source:
              path: charts/demo-app
              repoURL: (join('/', ['https://github.com', $repo]))
              targetRevision: main
              helm:
                valuesObject:
                  repo:
                    name: ($repo)
                  keptn:
                    appContext:
                      preDeploymentTasks:
                      - maintenance-window-check
                valueFiles:
                  - ../../gitops/dev/demo-app/values.yaml
                  - ../../gitops/dev/demo-app/values-specific.yaml
                parameters:
                  - name: commitID
                    value: $ARGOCD_APP_REVISION
                  - name: serviceVersion
                    value: v2
                  - name: service.nodePort
                    value: '31106'
            destination:
              server: https://kubernetes.default.svc
              namespace: ($namespace)
            syncPolicy:
              automated:
                prune: true
                selfHeal: true
    - assert:
        resource:
          apiVersion: lifecycle.keptn.sh/v1
          kind: KeptnTask
          spec:
            checkType: pre
            context:
              appName: demo-app
              appVersion: v2
              objectType: App
              taskType: pre
              workloadName: ""
            taskDefinition: maintenance-window-check
          status:
            status: Succeeded
```

This chainsaw test deletes the existing KeptnApp and the demo app Deployment.
ArgoCD will recreate the Deployment and Keptn will run the pre-deployment task.
Chainsaw will verify that a KeptnTask was executed and that the execution was successful.

### Run the chainsaw test

The test above can be run with:

```bash
chainsaw test tests/05-pre-deployment-checks
```

If the test succeeds you should see an output similar to:
```
Version: 0.2.5
Loading default configuration...
- Using test file: chainsaw-test
- TestDirs [./exercises/01-keptn-tasks/]
- SkipDelete false
- FailFast false
- ReportFormat ''
- ReportName ''
- Namespace ''
- FullName false
- IncludeTestRegex ''
- ExcludeTestRegex ''
- ApplyTimeout 5s
- AssertTimeout 30s
- CleanupTimeout 30s
- DeleteTimeout 15s
- ErrorTimeout 30s
- ExecTimeout 5s
- DeletionPropagationPolicy Background
- Template true
- NoCluster false
- PauseOnFailure false
Loading tests...
- maintenance-window-check-dev (./exercises/01-keptn-tasks/)
Loading values...
Running tests...
=== RUN   chainsaw
=== PAUSE chainsaw
=== CONT  chainsaw
=== RUN   chainsaw/maintenance-window-check-dev
=== PAUSE chainsaw/maintenance-window-check-dev
=== CONT  chainsaw/maintenance-window-check-dev
    | 00:14:44 | maintenance-window-check-dev | step-1   | TRY       | RUN   |
    | 00:14:44 | maintenance-window-check-dev | step-1   | DELETE    | RUN   | lifecycle.keptn.sh/v1/KeptnApp @ demo-app-dev/demo-app
    | 00:14:45 | maintenance-window-check-dev | step-1   | DELETE    | OK    | lifecycle.keptn.sh/v1/KeptnApp @ demo-app-dev/demo-app
    | 00:14:45 | maintenance-window-check-dev | step-1   | DELETE    | DONE  | lifecycle.keptn.sh/v1/KeptnApp @ demo-app-dev/demo-app
    | 00:14:45 | maintenance-window-check-dev | step-1   | DELETE    | RUN   | apps/v1/Deployment @ demo-app-dev/demo-app
    | 00:14:45 | maintenance-window-check-dev | step-1   | DELETE    | OK    | apps/v1/Deployment @ demo-app-dev/demo-app
    | 00:14:45 | maintenance-window-check-dev | step-1   | DELETE    | DONE  | apps/v1/Deployment @ demo-app-dev/demo-app
    | 00:14:45 | maintenance-window-check-dev | step-1   | ASSERT    | RUN   | lifecycle.keptn.sh/v1/KeptnApp @ demo-app-dev/demo-app
    | 00:15:15 | maintenance-window-check-dev | step-1   | ASSERT    | DONE  | lifecycle.keptn.sh/v1/KeptnApp @ demo-app-dev/demo-app
    | 00:15:15 | maintenance-window-check-dev | step-1   | ASSERT    | RUN   | lifecycle.keptn.sh/v1/KeptnTask @ demo-app-dev/*
    | 00:15:25 | maintenance-window-check-dev | step-1   | ASSERT    | DONE  | lifecycle.keptn.sh/v1/KeptnTask @ demo-app-dev/*
    | 00:15:25 | maintenance-window-check-dev | step-1   | TRY       | DONE  |
--- PASS: chainsaw (0.00s)
    --- PASS: chainsaw/maintenance-window-check-dev (40.17s)
PASS
Tests Summary...
- Passed  tests 1
- Failed  tests 0
- Skipped tests 0
Done.
```

### Summary

In this exercise, you have learned how to use Chainsaw to validate that pre-deployment tasks defined in your KeptnAppContext are correctly executed when the demo app is updated.
