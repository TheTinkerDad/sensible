# What is Sensible?
A small tool that provides monitoring for your Linux server via Home Assistant sensors and MQTT discovery.

By default Sensible comes with only a few example sensors, but it is basically a framework that enables you to quickly prototype and implement your own sensors.

Thanks to MQTT discovery, its integration with Home Assistant is as smooth as possible.

![Sensible as a device in Home Assistant](media/ha-device.jpg?raw=true "Sensible's MQTT based integration in Home Assistant")

*Note*: Yes, it's probably a temporary name, but I wanted to have something that at least a bit makes sense... (Pun intended!)

# How it works?

Sensible is currently a framework application that works with Home Assistant and MQTT discovery.
You can configure sensors as plugins for the framework and the sensors will appear in Home Assistant.
There are currently two ways to implement sensors, although this part is still under development.
First, you can code them in Golang and build them as part of Sensible.
Second, you can implement them as unix shell scripts. In this case, you don't need to build Sensible, but you can use a prebuilt binary.

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

This is currently done via a the file /etc/sensible/settings.yaml

A sample file looks like this

```
mqtt:
    hostname: 127.0.0.1
    port: "1883"
    username: ""
    password: ""
    clientid: sensible_mqtt_client
discovery:
    devicename: some-unique-identifier
    prefix: homeassistant
plugins:
    - name: Sensible Heartbeat
      kind: internal
      sensorid: heartbeat
      script: ""
      icon: mdi:wrench-check
    - name: Sensible Heartbeat NR
      kind: internal
      sensorid: heartbeat_NR
      script: ""
      icon: mdi:wrench-check
    - name: Sensible Boot Time
      kind: internal
      sensorid: boot_time
      script: ""
      icon: mdi:clock
    - name: Sensible System Time
      kind: internal
      sensorid: system_time
      script: ""
      icon: mdi:clock
    - name: Sensible Root Disk Free
      kind: script
      sensorid: root_free
      script: root_free.sh
      icon: mdi:harddisk
    - name: Sensible Host IP Address
      kind: script
      sensorid: ip_address
      script: ip_address.sh
      icon: mdi:harddisk
```

# Example scripts

There are a couple of scripts under the examples/scripts folders that are also configured to act as sensors in the above example configuration file.

The only requirement for these scripts is that they should be simple, with an execution time no longer than 1-2 seconds and they should only output the value that is meant to be a sensor value. E.g. the ip_address.sh script only outputs an IP / CIDR.

# Example usage

 * [A standalone systemd service on Linux servers](examples/systemd/README.md)
 * [Plugged into Docker containers](examples/docker/README.md) as a background process
 * [Plugged into LXC/LXD containers](examples/lxc/README.md) as a service
 * [As a standalone service on Raspberry Pi's](examples/raspberry-pi/README.md)
 * [As a system service on Windows](examples/windows/README.md)
 
# Development and planned features

 * There should be a way to implement sensors in Go for fully customized sensor data (plugin architecture)
 * Authentication for the API and a way to disable it
 * Configuration via environment variables
 * A lot of small TODO items in the code
 
