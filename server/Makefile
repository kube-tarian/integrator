OPEN_API_CODEGEN := github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

${OPEN_API_CODEGEN}:
	$(eval TOOL=$(@:%=%))
	@echo Installing ${TOOL}...
	go install $(@:%=%)

tools: ${OPEN_API_CODEGEN}

OPEN_API_DIR = ./api

oapi-gen: tools oapi-gen-server

oapi-gen-server:
	$(eval APP_NAME=server)
	@echo Generating server for ${APP_NAME}
	@mkdir -p ${OPEN_API_DIR}
	${GOBIN}/oapi-codegen -config ./cfg.yaml ./openapi.yaml
