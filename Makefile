.EXPORT_ALL_VARIABLES:
VERSION = 0.2.0
BUILDDATE = $$(date)

prepare:
	mkdir -p dist
	go mod download github.com/eclipse/paho.mqtt.golang
	go mod download github.com/google/uuid
	go mod download github.com/TheTinkerDad/go.pipe

build: test compile upx
	echo "Successfully built Sensible."

build-noupx: test compile
	echo "Successfully built Sensible."

docker-example:
	cp examples/docker/* dist/
	cp -R examples/scripts dist/
	$(SHELL) -c "cd dist;docker build -t thetinkerdad/sensible-nginx-test ."

test:
	go test

compile:
	# go build -ldflags="-w" -o dist/sensible
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -X 'TheTinkerDad/sensible/releaseinfo.BuildTime=$(BUILDDATE)' -X 'TheTinkerDad/sensible/releaseinfo.Version=$(VERSION)'" -a -installsuffix cgo  -o dist/sensible

upx:
	upx -9 dist/sensible

run:
	$(SHELL) -c "cd dist;./sensible"

release-linux-amd64: build
	$(SHELL) -c "cd dist;tar cvzf sensible-linux-amd64-$(VERSION).tar.gz sensible"

release-rpi-armhf: build
	$(SHELL) -c "cd dist;tar cvzf sensible-rpi-armhf-$(VERSION).tar.gz sensible"

release-rpi-arm64: build
	$(SHELL) -c "cd dist;tar cvzf sensible-rpi-arm64-$(VERSION).tar.gz sensible"