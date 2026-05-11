.PHONY: dockerimage
dockerimage:
	@echo "building Jump Jump docker image..."
	docker build -t studiomj/jump-jump:latest -f build/package/Dockerfile .

.PHONY: docs
docs:
	swag init -g ./internal/app/routers/router.go

.PHONY: dev-ui
dev-ui:
	cd web-ui && npm run dev

.PHONY: build-ui
build-ui:
	cd web-ui && npm install && npm run build

.PHONY: build
build:
	go build ./...
