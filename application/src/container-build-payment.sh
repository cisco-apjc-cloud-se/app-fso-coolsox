#!/bin/sh

ARGS="$@"
VERSION="${1}"

if [ "${VERSION}x" == "x" ]; then
    echo "VERSION is not specified, please try again"
    exit 1
fi

# payment
echo ""
echo "Building payment container (version ${VERSION})"
cd payment
sh scripts/build.sh moskrive/fso-with-coolsox $VERSION
cd ..
echo ""
echo "Pushing payment container to dockerhub"
docker push moskrive/fso-with-coolsox:payment-${VERSION}
echo ""

