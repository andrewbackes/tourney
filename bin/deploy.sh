#!/bin/bash -x
cd "$(dirname "$0")"/..
mkdir -p build

# Update the image tag
TAG=$(git rev-parse --short HEAD)
sed s/:latest/:${TAG}/g ./k8s/deployment.yml > build/deployment.yml
cp ./k8s/service.yml build/service.yml

# Use gcloud+kubectl to apply the new deployment
cat <<EOF >build/gke-deploy.sh
#!/bin/bash -x
PROJECT_ID=ringed-furnace-185604
ZONE=us-west1-a
CLUSTER_NAME=g1-small
GOOGLE_APPLICATION_CREDENTIALS=/home/k8s/travis-ci.credentials.json

gcloud auth activate-service-account --key-file "\${GOOGLE_APPLICATION_CREDENTIALS}"
gcloud config set project \$PROJECT_ID
gcloud config set compute/zone \$ZONE

gcloud container clusters get-credentials \$CLUSTER_NAME --zone=\$ZONE

kubectl apply -f /home/build/deployment.yml
kubectl apply -f /home/build/service.yml
EOF

chmod +x build/gke-deploy.sh

docker run \
    -v "$(pwd)/k8s/:/home/k8s/" \
    -v "$(pwd)/build:/home/build" \
        andrewbackes/gcloud \
        /home/build/gke-deploy.sh