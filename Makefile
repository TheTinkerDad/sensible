prepare:
	mkdir -p dist
	go mod download github.com/eclipse/paho.mqtt.golang
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
	# go build -ldflags="-s -w" -o dist/sensible
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -installsuffix cgo  -o dist/sensible

upx:
	upx -9 dist/sensible

run:
	$(SHELL) -c "cd dist;./sensible"