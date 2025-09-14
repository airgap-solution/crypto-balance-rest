.PHONY: openapi openapi-ts-build openapi-ts-publish

openapi:
	scripts/openapi.sh

openapi-ts-build:
	@if [ -d "openapi/clientgen/ts/src" ]; then \
		cd openapi/clientgen/ts && npm run build; \
	else \
		echo "No TypeScript client generated yet. Run 'make openapi' first."; \
	fi

openapi-ts-publish: openapi-ts-build
	cd openapi/clientgen/ts && npm publish --access public
