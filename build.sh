set -ex


docker build -t ohttps-webhook .
docker tag ohttps-webhook riba2534/ohttps-webhook
docker push riba2534/ohttps-webhook
