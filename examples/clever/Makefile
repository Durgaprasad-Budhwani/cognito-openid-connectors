build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/clever main.go

create-keys:
	./create-key.sh

deploy:
	make build

	npx serverless deploy --stage $(STAGE) -v --region $(REGION)
