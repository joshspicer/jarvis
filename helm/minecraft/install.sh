#!/bin/bash

DOWNLOAD_WORLD_URL="*****"

helm install mcjava \
    --set minecraftServer.eula=true \
    --set minecraftServer.ops=1556450a-0dea-42d4-bb04-7ba2e9482411 \
    --set minecraftServer.serviceType=LoadBalancer \
    --set persistence.dataDir.enabled=true \
    --set minecraftServer.motd=SpiceCraft3 \
    --set minecraftServer.downloadWorldUrl="$DOWNLOAD_WORLD_URL" \
    itzg/minecraft