#!/bin/sh

cat /sys/class/thermal/thermal_zone0/temp | awk '{printf "%.2f", $0 / 1000}'