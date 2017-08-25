#!/bin/bash

# Stop the ECS service from trying to join a cluster.
stop ecs || true && mv /etc/init/ecs.conf /etc/init/ecs.conf.disabled

# Pull and run the service.
docker pull dbhi/data-models-service:latest
docker run -d -p "80:8123" dbhi/data-models-service
