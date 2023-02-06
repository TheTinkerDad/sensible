# What is Sensible?
A small tool that provides monitoring for your Linux server via Home Assistant sensors and MQTT discovery.

*Note*: Yes, it's probably a temporary name, but I wanted to have something that at least a bit makes sense... (Pun intended!)

# How it works?

It does stuff - obviously this needs to be updated :)

# Building Sensible

Being very early in development, Sensible is currently being built for Linux, using make:

This downloads the minimal dependencies required
```
make prepare  
```

This one builds the executable and packs it with UPX
```
make build    
```

Also builds the executable, but without apply UPX
```
make build    
```

It is also possible to build example code for Docker, etc - see the Example usage section for this.

# Configuration

TBD

# Example usage

 * [A standalone systemd service on Linux servers](examples/systemd/README.md)
 * [Plugged into Docker containers](examples/docker/README.md) as a background process
 * [Plugged into LXC/LXD containers](examples/lxc/README.md) as a service
 * [As a standalone service on Raspberry Pi's](examples/raspberry-pi/README.md)
 * [As a system service on Windows](examples/windows/README.md)
 
# Development and planned features

 * Sensors that take CLI output should be configured as plugins without the need to rebuild them
 * There should be a way to implement sensors in Go for fully customized sensor data
 * Authentication for the API and a way to disable it
 * Configuration via environment variables
 
