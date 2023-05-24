set -ex


CGO_ENABLED=0 GOOS=linux go build -o main .
docker build -t ohttps-webhook .
docker tag ohttps-webhook riba2534/ohttps-webhook
docker push riba2534/ohttps-webhook
