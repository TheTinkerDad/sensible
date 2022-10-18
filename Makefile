prepare:
	mkdir -p dist
	go install github.com/eclipse/paho.mqtt.golang
	go install github.com/TheTinkerDad/go.rice

build: test compile upx add-resources
	echo "Successfully built Sensible."

build-noupx: test compile add-resources
	echo "Successfully built Sensible."

test:
	go test

compile:
	go build -ldflags="-s -w" -o dist/sensible

upx:
	upx -9 dist/sensible

add-resources:
	rice append-dir --exec dist/sensible -d data

run:
	cd dist
	./sensible
	cd ..