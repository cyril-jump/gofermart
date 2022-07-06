go-generate:
	@echo "  >  Generating dependency files..."
	oapi-codegen --config=internal/swagger/models.cfg.yaml internal/swagger/gofermart.yaml
	oapi-codegen --config=internal/swagger/server.cfg.yaml internal/swagger/gofermart.yaml