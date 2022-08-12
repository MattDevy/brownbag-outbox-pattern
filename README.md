# Outbox pattern demo

## Setup
```sh
docker-compose up -d
source ./dev.env
go run ./cmd/server
```

### Install postman
```sh
brew install --cask postman
```

### Use postman for demo
- Open postman
- Import file outbox_demo.postman_collection.json
- Run each request in before,after,transaction,outbox folders