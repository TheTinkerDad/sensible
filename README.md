# What is Sensible?
A small tool that provides monitoring for your Linux server via Home Assistant sensors and MQTT discovery.

By default Sensible comes with only a few example sensors, but it is basically a framework that enables you to quickly prototype and implement your own sensors.

The below video showcases the first public version of Sensible and its capabilities. (This is the first installment of a tutorial / development update series.)

[![IMAGE ALT TEXT HERE](https://img.youtube.com/vi/21pho997KuA/0.jpg)](https://www.youtube.com/watch?v=21pho997KuA)

# Why should you use Sensible?

 * It's tiny! Currently the binary is approximately 2.4Mb in size! You can put it into a Docker container and you won't even notice it's there!

 * Thanks to MQTT discovery, its integration with Home Assistant is as smooth as possible.

   ![Sensible as a device in Home Assistant](media/ha-device.png?raw=true "Sensible's MQTT based integration in Home Assistant")

 * Because it follows basic MQTT / Home Assistant standards, it's easy to use with things like Lovelace UI, Node Red, you name it!

   ![Sensible sensors on the Lovelace UI](media/ha-lovelace-big.png?raw=true "Sensible's example sensors on the Lovelace UI")

 * It's fully opensource with a permissive license! You can fork it on GitHub and make your own version!

 * It has a control REST API that enables disabling sensor data publishing, etc. (still WIP though)

 * The developer behind is a veteran with 20+ years of experience, so the project is here to stay, you can expect support and future updates!
 
*Note*: Yes, it's probably a temporary name, but I wanted to have something that at least a bit makes sense... (Pun intended!)

# How it works?

Sensible is currently a framework application that works with Home Assistant and MQTT discovery.
You can configure sensors as plugins for the framework and the sensors will appear in Home Assistant.
There are currently two ways to implement sensors, although this part is still under development.
First, you can code them in Golang and build them as part of Sensible.
Second, you can implement them as unix shell scripts. In this case, you don't need to build Sensible, but you can use a prebuilt binary.

# Quickstart guide

 - Currently only Linux is supported - if you're running other OSes, sorry, you'll have to wait!

 - Grab one of the releases from https://github.com/TheTinkerDad/sensible/releases or build Sensible on your own (see below)

 - The .tar.gz file only contains the binary, extract it somewhere convenient.
 
 - Run it the first time with "./sensible -r" and it'll generate the default config file: /etc/sensible/settings.yaml and the required folders

 - Edit the config file to customize your settings

 - Scripts should be located under /etc/sensible/scripts (or in the folder you've configured in the settings.yaml file)

 - You can find the example scripts [here](examples/scripts) or you can start by making your own, they are rather simple

 - Add a sensor entry in the config file for each of your scripts like:
   ```
   - name: Sensible Host IP Address
     kind: script
     sensorid: ip_address
     script: ip_address.sh
     icon: mdi:check-network
   ```

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

Also builds the executable, but without applying UPX
```
make build    
```

It is also possible to build example code for Docker, etc - see the Example usage section for this.

# Configuration

This is currently done via a the file /etc/sensible/settings.yaml

A sample file looks like this

```
general:
    logfile: /var/log/sensible/sensible.log
    scriptlocation: /etc/sensible/scripts/
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
      icon: mdi:check-network
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

 * Security! MQTT encryption and all the bells and whistles to make it production ready ASAP!
 * There should be a way to implement sensors in Go for fully customized sensor data (plugin architecture) without rebuilding Sensible itself
 * Authentication for the REST API and a way to disable it
 * Documentation for the REST API
 * A way to control Sensible via MQTT
 * Configuration via environment variables
 * A lot of small TODO items in the code
 
