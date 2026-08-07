#!/bin/bash
docker compose -p myurlshortener -f docker-compose.yaml --env-file vars.env up -d --build
read -p "Press enter to continue..."
