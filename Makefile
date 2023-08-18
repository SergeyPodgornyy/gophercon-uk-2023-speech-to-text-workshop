build:
	docker buildx build --platform linux/amd64 -t ytt-amd64 --load -f Dockerfile .

download-lang-model:
	./download_model.sh tiny.en

run:
	docker run --platform linux/amd64 -v "$(pwd)":/data -it --rm ytt-amd64 get "JzPfMbG1vrE"
