# Running Sensible as a standalone systemd service on Linux servers

## Currently tested OSes

* Ubuntu 18.x --> 20.x
* Debian 10.x --> 11.x

## Steps

### 1. Preparation
You need download and unpack/build Sensible binary program

### 2. Create user and basic structure of files
```shell
mv sensible /usr/local/bin/sensible
chmod a+x /usr/local/bin/sensible
useradd sensible
mkdir /var/log/sensible
chown sensible /var/log/sensible
mkdir -p /etc/sensible/scripts
```

### 3. Create file `/etc/sensible/settings.yaml`

You can use "sensible -r" to create the needed folders and the example settings.yaml file.

This will create /etc/sensible, /etc/sensible/scripts and /etc/sensible/settings.yaml

### 4. Give access to the created files to Sensible

```shell
chown -R sensible /etc/sensible
```

### 5. Add your scripts

```shell
cp <your-script-name-or-folder> /etc/sensible/scripts/
chown -R sensible /etc/sensible/scripts
```

### 6. Create and start service

create file `/lib/systemd/system/sensible.service`
with contents

```ini
[Unit]
Description=Sensible monitoring service
After=network.target

[Service]
ExecStart=/usr/local/bin/sensible
User=sensible
Restart=on-failure

[Install]
WantedBy=multi-user.target

```

Activate service

```shell
systemctl daemon-reload
systemctl enable sensible.service
systemctl start sensible.service
```

Check status of service

```shell
systemctl status sensible.service
```