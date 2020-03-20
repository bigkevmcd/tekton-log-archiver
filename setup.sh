#!/bin/sh
kubectl apply -f https://github.com/tektoncd/pipeline/releases/download/v0.10.1/release.yaml
kubectl apply -f https://github.com/tektoncd/triggers/releases/download/v0.3.1/release.yaml
kubectl create -f deploy/service_account.yaml
kubectl create -f deploy/role.yaml
kubectl create -f deploy/role_binding.yaml
kubectl create -f deploy/operator.yaml
