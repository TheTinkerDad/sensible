#!/bin/sh

ip a s eth0 | awk '/inet / {print$2}'