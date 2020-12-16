.PHONY: all clean get

GOPATH=$(shell pwd)

all: get
	test -d dist || mkdir dist
	GOPATH=$(GOPATH) GOOS="linux" GOARCH="amd64" go build github.com/bernarpa/photo && mv photo dist/photo-linux && chmod +x dist/photo-linux
	GOPATH=$(GOPATH) GOOS="darwin" GOARCH="amd64" go build github.com/bernarpa/photo && mv photo dist/photo-mac && chmod +x dist/photo-mac
	GOPATH=$(GOPATH) GOOS="windows" GOARCH="amd64" go build github.com/bernarpa/photo && mv photo.exe dist/photo-win.exe

clean:
	rm -fr bin/ pkg/ dist/ src/github.com/tmc/ src/github.com/kballard/ src/github.com/rwcarlsen/ src/golang.org/

get: src/github.com/rwcarlsen/goexif/exif/exif.go src/golang.org/x/crypto/go.mod src/github.com/tmc/scp/scp.go

src/github.com/rwcarlsen/goexif/exif/exif.go:
	GOPATH=$(GOPATH) go get github.com/rwcarlsen/goexif/exif

src/golang.org/x/crypto/go.mod:
	GOPATH=$(GOPATH) go get golang.org/x/crypto/ssh

src/github.com/tmc/scp/scp.go:
	GOPATH=$(GOPATH) go get github.com/tmc/scp
