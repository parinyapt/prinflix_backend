install:
	go mod tidy

run:
	go run cmd/main.go -mode=development

runprod:
	go run cmd/main.go -mode=production

dockerbuild:
	docker build -t prinflix_backend:v1.0.0 --platform linux/amd64 .

dockertag:
	docker tag prinflix_backend:v1.0.0 ghcr.io/parinyapt/prinflix_backend:v1.0.0

dockerpush:
	docker push ghcr.io/parinyapt/prinflix_backend:v1.0.0