.PHONY: dockerimage
dockerimage:
	@echo "building Jump Jump docker image..."
	docker build -t studiomj/jump-jump:latest -f build/package/Dockerfile .
