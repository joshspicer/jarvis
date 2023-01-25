#!/bin/bash

# TIP: Set these three variables as Codespaces Secrets (Command Palette: 'Codespaces: Manage User Secrets')

az account set --subscription $AZ_SUBSCRIPTION
az aks get-credentials --resource-group $AZ_RESOURCE_GROUP --name $AZ_CLUSTER_NAME