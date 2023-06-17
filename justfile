cluster_name := "auth-server"
namespace_name := "auth-server"
registry_name := "auth-server-registry.localhost"
registry_port := "5000"

start: create-cluster setup-context wait-for-traefik install-operator
stop: delete-cluster clear-context

create-cluster:
    #!/usr/bin/env bash
    set -euxo pipefail

    if ! k3d cluster list | grep -qw {{ cluster_name }}; then
        k3d cluster create {{ cluster_name }} \
            --registry-create {{ registry_name }}:0.0.0.0:{{ registry_port }} \
            --kubeconfig-update-default=false \
            -p "80:80@loadbalancer" \
            --wait;
    else
        echo "cluster {{ cluster_name }} already exists!"
    fi

setup-context:
    @mkdir -p .scratch
    @k3d kubeconfig get {{ cluster_name }} > .scratch/kubeconfig
    chmod og-r .scratch/kubeconfig

wait-for-traefik:
    #!/usr/bin/env bash
    LAST_STATUS=""
    STATUS=""
    
    echo "Waiting for traefik to start..."

    while [[ "$STATUS" != "Running" ]]; do
        sleep 1
        STATUS=$(kubectl --kubeconfig .scratch/kubeconfig get pods -n kube-system -o json | jq '.items[] | select(.metadata.name | startswith("traefik")) | .status.phase' -r)
        if [[ "$STATUS" != "$LAST_STATUS" ]]; then
            echo "traefik pod is '$STATUS'"
        fi
        LAST_STATUS="$STATUS"
    done

    echo "done"

install-operator:
    kubectl create namespace db-operator --dry-run=client -o yaml | kubectl --kubeconfig .scratch/kubeconfig apply -f -
    helm upgrade --install db-operator oci://docker.io/benwright/db-operator-chart \
        --kubeconfig .scratch/kubeconfig \
        --namespace db-operator \
        --version=v1.0.5 \
        --set image=benwright/db-operator:v1.0.5 \
        --wait

delete-cluster:
    if k3d cluster list | grep -qw {{ cluster_name }}; then \
        k3d cluster delete {{ cluster_name }}; \
    fi

clear-context:
    if [[ -f .scratch/kubeconfig ]]; then \
        rm .scratch/kubeconfig; \
    fi

tilt:
    KUBECONFIG=.scratch/kubeconfig tilt up

test:
    go test --short -v ./...

int-test:
    go test --run Integration -v ./...

build PATH_TO_CODE IMAGE_TAG:
    mkdir -p "{{PATH_TO_CODE}}/dist";
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o "{{PATH_TO_CODE}}/dist/app" "{{PATH_TO_CODE}}/main.go";
    docker build -t "{{IMAGE_TAG}}" -f deploy/golang.Dockerfile "{{PATH_TO_CODE}}/dist"
