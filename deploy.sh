#!/bin/bash

git pull --rebase

docker build -t hah-open-cookbook:latest .

docker stop hah-open-cookbook 2>/dev/null || true
docker rm hah-open-cookbook 2>/dev/null || true

docker run -d -p 127.0.0.1:80:80 -v $(pwd)/data/cookbook.db:/app/cookbook.db -v $(pwd)/data/.env:/app/.env --restart unless-stopped --name hah-open-cookbook hah-open-cookbook:latest
    
sleep 3

docker ps | grep hah-open-cookbook

echo "Deployment complete!"