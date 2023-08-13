#!/bin/bash
# https://learn.microsoft.com/en-us/azure/aks/ingress-tls?tabs=azure-cli

# Add the Jetstack Helm repository
helm repo add jetstack https://charts.jetstack.io

# Update your local Helm chart repository cache
helm repo update

kubectl label namespace ingress-basic cert-manager.io/disable-validation=true
helm install cert-manager jetstack/cert-manager \
  --namespace ingress-basic \
  --set installCRDs=true \
  --set nodeSelector."kubernetes\.io/os"=linux