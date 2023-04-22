#!/bin/sh

df / | tail -n 1 | awk '/ / {printf "%.2f", $4 / 1048576}'
