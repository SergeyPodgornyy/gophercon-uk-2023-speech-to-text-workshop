build:
	docker buildx build --platform linux/amd64 -t ytt-amd64 --load -f Dockerfile .

download-lang-model:
	./download_model.sh tiny.en

network:
	docker network create surrealdb-network

database:
	docker run --name surrealdb --network surrealdb-network --restart always -v "$(pwd)/data:/data" -p 8000:8000 surrealdb/surrealdb:latest  start --user root --pass root

run:
	docker run --platform linux/amd64 --network surrealdb-network -v "$(pwd)":/data -e IP_ADDRESS=10.100.95.88 -it --rm ytt-amd64 get "JzPfMbG1vrE" 

