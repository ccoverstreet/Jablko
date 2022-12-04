local: 
	docker buildx build --platform="linux/amd64" . --tag=ccoverstreet/testmod:localdev --load


publish:
	docker buildx build --platform="linux/amd64,linux/arm64" . --tag=ccoverstreet/testmod:latest --push


