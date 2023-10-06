SHELL := /bin/bash

GOLANG          := golang:1.21-alpine
ALPINE          := alpine:3.18
KIND            := kindest/node@sha256:3966ac761ae0136263ffdb6cfd4db23ef8a83cba8a463690e98317add2c9ba72
POSTGRES        := postgres:15.3
VAULT           := hashicorp/vault:1.13
GRAFANA         := grafana/grafana:9.5.2
PROMETHEUS      := prom/prometheus:v2.44.0
TEMPO           := grafana/tempo:2.1.1
TELEPRESENCE    := datawire/tel2:2.13.3
KIND_CLUSTER    := bersen-starter-cluster
NAMESPACE       := boost-sales-system
APP             := boost-sales
BASE_IMAGE_NAME := bersennaidoo
SERVICE_NAME    := boost-sales-api
VERSION         := 1.0.0
SERVICE_IMAGE   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(VERSION)
METRICS_IMAGE   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME)-metrics:$(VERSION)


.PHONY: run kind-up build kind-load kind-apply kind-down kind-status kind-logs kind-restart kind-update kind-update-apply kind-status-boost-sales

setup: install-go init-go

install-go:
	wget "https://golang.org/dl/go$(GO_VERSION).linux-amd64.tar.gz"
	sudo tar -C /usr/local -xzf go$(GO_VERSION).linux-amd64.tar.gz
	rm go$(GO_VERSION).linux-amd64.tar.gz

init-go:
	echo 'export PATH=$$PATH:/usr/local/go/bin' >> $${HOME}/.bashrc
	echo 'export PATH=$$PATH:$${HOME}/go/bin' >> $${HOME}/.bashrc


run:
	go run application/services/transport/rest/boost-sales-api/main.go

test:
	go test ./... -count=1
	staticcheck -checks=all ./...

kind-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config infrastructure/k8s/kind/kind-config.yaml

build:
	docker build \
		-f infrastructure/docker/Dockerfile.boost-sales-api \
		-t $(SERVICE_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

kind-load:
	cd infrastructure/k8s/kind/boost-sales-pod; kustomize edit set image bersennaidoo/boost-sales-api=bersennaidoo/boost-sales-api:$(VERSION)
	kind load docker-image bersennaidoo/boost-sales-api:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	#kustomize build zarf/k8s/kind/database-pod | kubectl apply -f -
	#kubectl wait --namespace=database-system --timeout=120s --for=condition=Available deployment/database-pod
	#kustomize build zarf/k8s/kind/zipkin-pod | kubectl apply -f -
	#kubectl wait --namespace=zipkin-system --timeout=120s --for=condition=Available deployment/zipkin-pod
	kustomize build infrastructure/k8s/kind/boost-sales-pod | kubectl apply -f -
	#cat infrastructure/k8s/base/boost-sales-pod/base-boost-sales.yaml | kubectl apply -f -

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces


kind-logs:
	kubectl logs -l app=boost-sales --all-containers=true -f --tail=100 -n boost-sales-system

kind-restart:
	kubectl rollout restart deployment boost-sales-pod -n boost-sales-system

kind-update: build kind-load kind-restart

kind-update-apply: build kind-load kind-apply

kind-status-boost-sales:
	kubectl get po -o wide -w -n boost-sales-system

