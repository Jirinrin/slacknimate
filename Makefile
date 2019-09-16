.PHONY: build_deps build goinstall gouninstall package_deps package clobber
.DEFAULT_GOAL := build

build_deps:
	@type go >/dev/null 2>&1 || \
		{ echo >&2 "I require go but it is not installed.  Aborting."; exit 1; }

build: build_deps
	go build -o bin/slacknimate.exe


# Standard go install, for people with a valid go / $GOPATH setup
goinstall:
	go install .

gouninstall:
	rm $(GOPATH)/bin/slacknimate.exe


# For cross compiling and packaging releases
package_deps:
	go get github.com/laher/goxc

package: package_deps build
	goxc -pv=`./bin/slacknimate.exe -v | cut -d' ' -f3` \
	     --resources-include="README*,LICENSE*,examples" \
	     -d="builds" -bc="linux,!arm darwin windows" xc archive rmbin


clobber:
	rm -rf builds bin
