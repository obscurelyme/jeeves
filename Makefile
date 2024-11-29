build:
	go build

user-install: build
	go install

user-uninstall:
	rm $(GOBIN)/jeeves

brk:
	aws bedrock-runtime invoke-model \
--model-id amazon.titan-text-express-v1 \
--body '{"inputText": "Can you give me an example of an IAM role in yaml format that will give an arbitrary lambda function access to an S3 bucket?", "textGenerationConfig" : {"maxTokenCount": 512, "temperature": 0.5, "topP": 0.9}}' \
--cli-binary-format raw-in-base64-out \
invoke-model-output-text.json