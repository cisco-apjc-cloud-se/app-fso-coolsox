#!/bin/sh

ARGS="$@"
VERSION="${1}"

if [ "${VERSION}x" == "x" ]; then
    echo "VERSION is not specified, please try again"
    exit 1
fi


# user
echo ""
echo "Building user container (version ${VERSION})"
cd user
docker build -t moskrive/fso-with-coolsox:user-$VERSION .
cd ..
echo ""
echo "Pushing user container to dockerhub"
docker push moskrive/fso-with-coolsox:user-${VERSION}
echo ""

