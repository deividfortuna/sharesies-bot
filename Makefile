.PHONY: build
.PHONY: publish

image = deividfortuna/auto-invest-sharesies
tag = 0.0.1

build: 
	docker build -t ${image}:${tag} .
publish: build
	docker push ${image}:${tag}