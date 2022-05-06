#!/bin/sh

ARGS="$@"
VERSION="${1}"

if [ "${VERSION}x" == "x" ]; then
    echo "VERSION is not specified, please try again"
    exit 1
fi

docker image rm moskrive/fso-with-coolsox:carts-$VERSION
docker image rm moskrive/fso-with-coolsox:catalogue-$VERSION
docker image rm moskrive/fso-with-coolsox:catalogue-db-$VERSION
docker image rm moskrive/fso-with-coolsox:front-end-$VERSION
docker image rm moskrive/fso-with-coolsox:front-end-appd-brum-$VERSION
docker image rm moskrive/fso-with-coolsox:orders-$VERSION
docker image rm moskrive/fso-with-coolsox:payment-$VERSION
docker image rm moskrive/fso-with-coolsox:queue-master-$VERSION
docker image rm moskrive/fso-with-coolsox:shipping-$VERSION
docker image rm moskrive/fso-with-coolsox:user-$VERSION
docker image rm moskrive/fso-with-coolsox:user-db-$VERSION
docker image rm moskrive/fso-with-coolsox:load-test-$VERSION
