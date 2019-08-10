# Image URL to use all building/pushing image targets
REGISTRY ?= local
IMG_API ?= k8s_packet_trace/api
IMG_AGENT ?= k8s_packet_trace/agent

IMAGE_TAG ?= 0.1


# Build the docker image
docker-build:
	docker build -f./docker/agent/Dockerfile -t ${REGISTRY}/${IMG_AGENT}:${IMAGE_TAG} .
	docker build -f./docker/api/Dockerfile -t ${REGISTRY}/${IMG_API}:${IMAGE_TAG} .

# initlize kind k8s cluster
cluster-up:
	kind create cluster --config ./cluster/kind-cluster.yaml

# Copy images to the kinds cluster
copy-images:
	kind load docker-image ${REGISTRY}/${IMG_AGENT}:${IMAGE_TAG}
	kind load docker-image ${REGISTRY}/${IMG_API}:${IMAGE_TAG}
