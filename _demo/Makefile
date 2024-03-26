SHELL=/bin/bash
ROOT ?= ../
APPNAME = menu
OUTDIR = $(ROOT)/out/$(APPNAME)

ifeq ($(MODE), debug)
    BUILDARGC = -gcflags "-N -l"
else
    BUILDARGC = -ldflags "-w -s"
endif

default:run

run:
	@go run *.go

build:
	go build -o $(OUTDIR)/$(APPNAME) -trimpath $(BUILDARGC) *.go

buildwin:
	GOOS=windows GOARCH=amd64 go build -o $(OUTDIR)/$(APPNAME).exe -trimpath $(BUILDARGC) *.go

clean:
	-rm -rf $(OUTDIR)
	@echo "all clean"

task:
	@grep -o -e "^\S*:" Makefile | grep -v default: | grep -v .PHONY | sort

.PHONY: all task
