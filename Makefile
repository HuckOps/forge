.PHONY: depend clean
depend:
	@go mod tidy
	@cd web; yarn build

clean:
	@rm -rf bin
	@rm -rf web/dist

backend: clean depend
	@go build -o bin/server cmd/server/server.go
	@go build -o bin/agent cmd/agent/agent.go
	@go build -o bin/gateway cmd/gateway/gateway.go
	@go build -o bin/worker cmd/worker/worker.go

frontend: clean depend
	@cd web; yarn build
