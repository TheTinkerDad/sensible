prepare:
	mkdir -p dist
	go install github.com/eclipse/paho.mqtt.golang

build: test compile upx
	echo "Successfully built Sensible."

build-noupx: test compile
	echo "Successfully built Sensible."

docker-example:
	cp examples/docker/* dist/
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