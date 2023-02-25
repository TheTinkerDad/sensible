#!/bin/sh

df -h / | tail -n 1 | awk '/ / {print$4}'