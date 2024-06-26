cluster_name := "auth-server"
namespace_name := "auth-server"
registry_name := "auth-server-registry.localhost"
registry_port := "5000"

clean:
    rm -rf system-tests/node_modules
    rm -rf system-tests/playwright-report
    rm -rf system-tests/test-results
    docker system prune --all --volumes

cert:
    #!/usr/bin/env zsh
    openssl req -x509 -out .scratch/localhost.crt -keyout .scratch/localhost.key \
        -newkey rsa:2048 -nodes -sha256 \
        -subj '/CN=localhost' -extensions EXT -config <( printf "[dn]\nCN=localhost\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS:localhost\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth")
    sudo security add-trusted-cert \
        -d \
        -r trustRoot \
        -k /Library/Keychains/System.keychain \
        .scratch/localhost.crt

remove-cert:
    sudo security delete-certificate -t -c localhost /Library/Keychains/System.keychain

start: create-cluster setup-context wait-for-traefik install-operator cert
stop: delete-cluster clear-context remove-cert

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
    kubectl create namespace operators --dry-run=client -o yaml | kubectl --kubeconfig .scratch/kubeconfig apply -f -
    helm upgrade --install db-operator oci://docker.io/benwright/db-operator-chart \
        --kubeconfig .scratch/kubeconfig \
        --namespace operators \
        --version=v2.0.7 \
        --wait

delete-cluster:
    if k3d cluster list | grep -qw {{ cluster_name }}; then \
        k3d cluster delete {{ cluster_name }}; \
    fi

clear-context:
    if [[ -f .scratch/kubeconfig ]]; then \
        rm .scratch/kubeconfig; \
    fi

test:
    go test --short -v ./...

int-test:
    go test --run Integration -v ./...

build PATH_TO_CODE IMAGE_TAG:
    mkdir -p "{{PATH_TO_CODE}}/dist";
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o "{{PATH_TO_CODE}}/dist/app" "{{PATH_TO_CODE}}/main.go";
    if [ -d "{{PATH_TO_CODE}}/static" ]; then \
        rm -rf "{{PATH_TO_CODE}}/dist/www/static"; \
        mkdir -p "{{PATH_TO_CODE}}/dist/www/static"; \
        cp -r {{PATH_TO_CODE}}/static/ {{PATH_TO_CODE}}/dist/www/static/; \
        docker build -t "{{IMAGE_TAG}}" -f deploy/golang-web.Dockerfile "{{PATH_TO_CODE}}/dist"; \
    else \
        docker build -t "{{IMAGE_TAG}}" -f deploy/golang.Dockerfile "{{PATH_TO_CODE}}/dist"; \
    fi

build-mig PATH_TO_CODE IMAGE_TAG:
    docker build -t "{{IMAGE_TAG}}" -f deploy/migrations.Dockerfile "{{PATH_TO_CODE}}"

e2e:
    cd system-tests && npm run test