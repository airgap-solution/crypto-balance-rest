.PHONY: openapi openapi-ts-setup openapi-ts-build openapi-ts-clean openapi-ts-publish

openapi:
	scripts/openapi.sh

# Setup TypeScript client (run once)
openapi-ts-setup:
	cd openapi/clientgen/ts && npm install

# Build TypeScript client (if generated files exist)
openapi-ts-build:
	@if [ -d "openapi/clientgen/ts/src" ]; then \
		cd openapi/clientgen/ts && npm run build; \
	else \
		echo "No TypeScript client generated yet. Run 'make openapi' first."; \
	fi

# Clean TypeScript build artifacts
openapi-ts-clean:
	@if [ -d "openapi/clientgen/ts" ]; then \
		cd openapi/clientgen/ts && npm run clean; \
	fi

# Publish TypeScript client
openapi-ts-publish: openapi-ts-build
	cd openapi/clientgen/ts && npm publish
