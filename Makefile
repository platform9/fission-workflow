DOCKER_REPO=fission
DOCKER_TAG=latest
PROTO_TARGETS=$(shell find pkg -type f -name "*.pb.go")
PROTO_TARGETS+=$(shell find pkg -type f -name "*.pb.gw.go")
SRC_TARGETS=$(shell find pkg -type f -name "*.go" | grep -v version.gen.go )
CHART_FILES=$(shell find charts/fission-workflows -type f)
VERSION=head

.PHONY: build generate prepush verify test changelog build-linux build-osx build-windows

build build-linux fission-workflows fission-workflows-bundle fission-workflows-proxy:
	build/build.sh

build-osx:
	build/build-osx.sh

build-windows:
	build/build-windows.sh

docker-build:
	# DOCKER_REPO=${DOCKER_REPO}
	# DOCKER_TAG=${DOCKER_TAG}
	build/docker.sh ${DOCKER_REPO} ${DOCKER_TAG} 

generate: ${PROTO_TARGETS} examples/workflows-env.yaml pkg/api/events/events.gen.go

prepush: generate verify test

test:
	test/runtests.sh

verify:
	helm lint charts/fission-workflows/ > /dev/null
	hack/verify-workflows.sh
	hack/verify-gofmt.sh
	hack/verify-misc.sh
	hack/verify-govet.sh

clean:
	rm fission-workflows*

version pkg/version/version.gen.go: pkg/version/version.go ${SRC_TARGETS}
	hack/codegen-version.sh -o pkg/version/version.gen.go -v ${VERSION}

changelog:
	test -n "${GITHUB_TOKEN}" # $$GITHUB_TOKEN
	github_changelog_generator -t ${GITHUB_TOKEN} --future-release ${VERSION}

examples/workflows-env.yaml: ${CHART_FILES}
	hack/codegen-helm.sh

%.swagger.json: %.pb.go
	hack/codegen-swagger.sh

%.pb.gw.go %.pb.go: %.proto
	hack/codegen-grpc.sh

pkg/api/events/events.gen.go: pkg/api/events/events.proto
	python3 hack/codegen-events.py
