#!/bin/bash -x

cluster='Probe-A-Node'
echo "Deploying to cluster: ${cluster}"
namespace='Probe-A-Node'
echo "Deploying to namespace: ${namespace}"
imagetag=${DRONE_BUILD_NUMBER}
echo "Image tag: ${imagetag}"

bash /usr/local/bin/init-kubectl
echo "Settings up project, namespace, and kubeconfig"
rancher-projects --cluster-name ${cluster} --project-name TechOps --namespace ${namespace} --create-project true --create-namespace true --create-kubeconfig true --kubeconfig ~/.kube/config
export KUBECONFIG=~/.kube/config

if ! kubectl cluster-info
then
  echo "Problem connecting to the cluster"
  exit 1
fi

echo "Creating name and adding labels to namespace"
kubectl create ns ${namespace} --dry-run=client -o yaml | kubectl apply -f -
kubectl label ns ${namespace} team=techops --overwrite
kubectl label ns ${namespace} app=Probe-A-Node --overwrite

echo "Creating registry secret"
kubectl -n ${namespace} create secret docker-registry harbor-registry-secret \
--docker-server=harbor.coxedgecomputing.com \
--docker-username=${DOCKER_USERNAME} \
--docker-password=${DOCKER_PASSWORD} \
--dry-run=client -o yaml | kubectl apply -f -

echo "Deploying Helm Chart"
cd /drone/src/charts
helm upgrade --install Probe-A-Node ./Probe-A-Node \
--namespace ${namespace} \
-f ./Probe-A-Node/values.yaml \
--set image.tag=${DRONE_BUILD_NUMBER} \
--force
