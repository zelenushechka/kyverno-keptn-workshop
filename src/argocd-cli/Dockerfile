FROM ubuntu:latest

# Install dependencies
RUN apt-get update && apt-get install -y curl

# Install ArgoCD CLI
RUN curl -sSL -o /usr/local/bin/argocd https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-amd64 && \
    chmod +x /usr/local/bin/argocd

RUN curl -sSL -o /usr/local/bin/kubectl "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" \
    && chmod +x /usr/local/bin/kubectl    

COPY --chmod=755 rollback.sh /rollback.sh

# Set the entrypoint to Bash
ENTRYPOINT ["/bin/bash"]