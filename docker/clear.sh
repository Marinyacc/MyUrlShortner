#!/bin/bash
docker compose -f docker-compose.yaml --env-file vars.env down
read -p "Press enter to continue..."