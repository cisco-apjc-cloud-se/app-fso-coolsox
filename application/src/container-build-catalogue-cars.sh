#!/bin/sh

ARGS="$@"
VERSION="${1}"

if [ "${VERSION}x" == "x" ]; then
    echo "VERSION is not specified, please try again"
    exit 1
fi


# catalogue and catalogue-db
echo ""
echo "Building catalogue and catalogue-db containers (${VERSION})"
cd catalogue-cars
sh scripts/build.sh moskrive/fso-with-coolsox $VERSION
cd ..
echo ""
echo "pushing catalogue-cars container to dockerhub"
docker push moskrive/fso-with-coolsox:catalogue-${VERSION}
echo ""
echo "pushing catalogue-cars container to dockerhub"
docker push moskrive/fso-with-coolsox:catalogue-db-${VERSION}

