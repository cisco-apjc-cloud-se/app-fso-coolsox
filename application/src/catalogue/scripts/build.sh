#!/usr/bin/env bash

ARGS="$@"
REPO="${1}"
VERSION="${2}"

if [ "${REPO}x" == "x" ]; then
    echo "REPO is not specified, please try again"
    exit 1
fi

if [ "${VERSION}x" == "x" ]; then
    echo "VERSION is not specified, please try again"
    exit 1
fi

set -ev

export BUILD_VERSION="0.0.2-SNAPSHOT"
export BUILD_DATE=`date +%Y-%m-%dT%T%z`

SCRIPT_DIR=$(dirname "$0")

if [[ "$(uname)" == "Darwin" ]]; then
    DOCKER_CMD=docker
else
    DOCKER_CMD="sudo docker"
fi
CODE_DIR=$(cd $SCRIPT_DIR/..; pwd)
echo $CODE_DIR

cp -r "${CODE_DIR}"/images/ "${CODE_DIR}"/docker/catalogue/images/
cp -r "${CODE_DIR}"/cmd/ "${CODE_DIR}"/docker/catalogue/cmd/
cp "${CODE_DIR}"/*.go "${CODE_DIR}"/docker/catalogue/
mkdir -p "${CODE_DIR}"/docker/catalogue/vendor && cp "${CODE_DIR}"/vendor/manifest "${CODE_DIR}"/docker/catalogue/vendor/
cp "${CODE_DIR}"/golang-sdk-x64-linux-4.5.2.0.tar "${CODE_DIR}"/docker/catalogue/
# cp -r "${CODE_DIR}"/appdynamics/ "${CODE_DIR}"/docker/catalogue/vendor/appdynamics/ ## RW 4/5/22

IMAGE_BASE="${REPO}:catalogue"

# $DOCKER_CMD build -t ${IMAGE_BASE}-dev "${CODE_DIR}"/docker/catalogue;
# $DOCKER_CMD create --name catalogue ${IMAGE_BASE}-dev;
# $DOCKER_CMD cp catalogue:/app/main "${CODE_DIR}"/docker/catalogue/app;
# $DOCKER_CMD rm catalogue;
# $DOCKER_CMD build \
#   --build-arg BUILD_VERSION=$BUILD_VERSION \
#   --build-arg BUILD_DATE=$BUILD_DATE \
#   -t ${IMAGE_BASE}-${VERSION} \
#   -f "${CODE_DIR}"/docker/catalogue/Dockerfile-release "${CODE_DIR}"/docker/catalogue;
# $DOCKER_CMD image rm ${IMAGE_BASE}-dev;

## Use the single build container instead of alpine version ##
$DOCKER_CMD build -t ${IMAGE_BASE}-${VERSION} "${CODE_DIR}"/docker/catalogue;

$DOCKER_CMD build \
  -t ${IMAGE_BASE}-db-${VERSION} \
  -f "${CODE_DIR}"/docker/catalogue-db/Dockerfile "${CODE_DIR}"/docker/catalogue-db;
