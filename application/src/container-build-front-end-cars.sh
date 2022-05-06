#!/bin/sh

ARGS="$@"
VERSION="${1}"

if [ "${VERSION}x" == "x" ]; then
    echo "VERSION is not specified, please try again"
    exit 1
fi

# front-end
echo ""
echo "Building front-end-cars container (version ${VERSION})"
docker build -t moskrive/fso-with-coolsox:front-end-$VERSION front-end-cars/.
echo ""
echo "Pushing front-end-cars container to dockerhub"
docker push moskrive/fso-with-coolsox:front-end-${VERSION}
echo ""

