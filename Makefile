.PHONY: generate-kubeconfig
generate-kubeconfig:
	docker build -t openpitrix/generate-kubeconfig:latest .

.PHONY: dump-all
dump-all:
	docker build -t openpitrix/dump-all:latest -f ./Dockerfile.dump .