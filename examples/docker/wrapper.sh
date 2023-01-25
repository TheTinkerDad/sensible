#!/bin/bash

# Start Sensible as a background process
./sensible &
# Start the original process of your container (replace this as needed)
nginx -g 'daemon off;'