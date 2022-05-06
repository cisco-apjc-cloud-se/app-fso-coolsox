#!/bin/sh

ARGS="$@"
VERSION="${1}"

if [ "${VERSION}x" == "x" ]; then
    echo "VERSION is not specified, please try again"
    exit 1
fi

# shipping
echo ""
echo "Building shipping container (version ${VERSION})"
cd shipping
sh ./scripts/build.sh moskrive/fso-with-coolsox $VERSION
cd ..
echo ""
echo "Pushing shipping container to dockerhub"
docker push moskrive/fso-with-coolsox:shipping-${VERSION}
echo ""

