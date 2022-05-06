#!/bin/sh

ARGS="$@"
VERSION="${1}"

if [ "${VERSION}x" == "x" ]; then
    echo "VERSION is not specified, please try again"
    exit 1
fi


# orders
echo ""
echo "Building orders container (version ${VERSION})"
cd orders
sh scripts/build.sh moskrive/fso-with-coolsox $VERSION
cd ..
echo ""
echo "Pushing orders container to dockerhub"
docker push moskrive/fso-with-coolsox:orders-${VERSION}
echo ""
