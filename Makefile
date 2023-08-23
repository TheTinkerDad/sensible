.EXPORT_ALL_VARIABLES:
VERSION = 0.3.0
BUILDDATE = $$(date)
LASTCOMMIT = $$(git rev-parse --short HEAD)

build: test compile upx
	echo "Successfully built Sensible."

build-noupx: test compile
	echo "Successfully built Sensible."

docker-example:
	mkdir -p dist/docker-build/etc/sensible dist/docker-build/log
	cp dist/sensible dist/docker-build/
	cp -R examples/docker/* dist/docker-build
	cp -R examples/scripts dist/docker-build/etc/sensible
	$(SHELL) -c "cd dist/docker-build;docker build -t thetinkerdad/sensible-nginx-test .;cd -"

test:
	go test

compile:
	# go build -ldflags="-w" -o dist/sensible
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -X 'TheTinkerDad/sensible/releaseinfo.LastCommit=$(LASTCOMMIT)' -X 'TheTinkerDad/sensible/releaseinfo.BuildTime=$(BUILDDATE)' -X 'TheTinkerDad/sensible/releaseinfo.Version=$(VERSION)'" -a -installsuffix cgo  -o dist/sensible

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

release-amd64-deb: build
	$(SHELL) -c "mkdir -p dist/pkg-deb/amd64/etc/sensible/scripts"
	$(SHELL) -c "mkdir -p dist/pkg-deb/amd64/opt/sensible"
	$(SHELL) -c "mkdir -p dist/pkg-deb/amd64/var/log/sensible"
	$(SHELL) -c "cp dist/sensible dist/pkg-deb/amd64/opt/sensible"
	$(SHELL) -c "cd dist/pkg-deb/amd64/etc/sensible;../../opt/sensible/sensible --g"
	$(SHELL) -c "cd dist/pkg-deb;dpkg-deb --build --root-owner-group amd64"
	$(SHELL) -c "cd dist/pkg-deb;dpkg-name amd64.deb"