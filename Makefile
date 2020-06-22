.PHONY: generate-kubeconfig
generate-kubeconfig:
	docker build -t openpitrix/generate-kubeconfig:latest .