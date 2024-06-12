# Create a Lab Environment

## Fork the repository

Go to https://github.com/heckelmann/kyverno-keptn-workshop and fork the repository to your personal GitHub account.

Make sure the forked repository visability is set to `Public`.

![Fork Repository](assets/01-fork-repository.png)


## Start GitHub CodeSpace

In your fork, go to "Code" then switch to the "Codespaces" tab and click "Create codespace on main"

![Fork Repository](assets/01-create-codespace.png)

A new Window will open with the Codespace. This will take a few minutes to start.

## Change Application Path

To allow ArgoCD to use the application manifests stored in the `gitops` folder, select the gitops folder and search and replace the `https://github.com/heckelmann/kyverno-keptn-workshop.git` with your Repository URL.


## Create GitHub API Token and K8s Secret

Open the GitHub settings and navigate to `Developer settings` -> `Personal access tokens` -> `Fine-grained tokens` (https://github.com/settings/tokens?type=beta).

Select access only to your forked repository and set the permission on `Actions` to `read` and `write` access.

Note down the generated token.

![Fork Repository](assets/01-create-token.png)

Switch back to your Codespace and create a Kubernetes secret with the token:

```bash
GH_REPO_OWNER=<YOUR_GITHUB_USER>
GH_REPO=<YOUR_GITHUB_REPO>
GH_API_TOKEN=<YOUR_GITHUB_TOKEN>
kubectl create secret generic github-token -n demo-app-dev --from-literal=SECURE_DATA="{\"githubRepo\":\"${GH_REPO}\",\"githubRepoOwner\":\"${GH_REPO_OWNER}\",\"apiToken\":\"${GH_API_TOKEN}\"}"
```

## Start the Workshop Kubernetes Cluster

Open a new Terminal within your Codespace and run the following command:

```bash
make create
```