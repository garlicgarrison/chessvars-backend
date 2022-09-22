#!/bin/bash
cd /home/ec2-user/chessvars-monolith
docker-compose build --no-cache
docker-compose up -d