#!/bin/bash
sudo amazon-linux-extras install docker -y
sudo service docker start
echo "$GOOGLE_CREDENTIALS" > google_key.json
