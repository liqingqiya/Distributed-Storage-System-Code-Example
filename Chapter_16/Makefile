PROJECTDIR=$(shell pwd)
BINDIR=$(PROJECTDIR)/output/bin

.PHONY: clean all server client
all:server client

server:
	@echo "building server"
	@go build -gcflags=all="-N -l" -o $(BINDIR)/server example.com/hellofs/cmds/server

client:
	@echo "building client"
	@go build -gcflags=all="-N -l" -o $(BINDIR)/client example.com/hellofs/cmds/client
