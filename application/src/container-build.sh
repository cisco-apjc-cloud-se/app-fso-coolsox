#!/bin/sh

ARGS="$@"
VERSION="${1}"
REPO="${2}" # moskrive/fso-with-coolsox

if [ "${VERSION}x" == "x" ]; then
    echo "VERSION is not specified, please try again"
    exit 1
fi

# # carts
# echo "Building carts container (${VERSION})"
# cd carts
# sh scripts/build.sh ${REPO} ${VERSION}
# cd ..
# echo ""
# echo "pushing carts container to repository ${REPO}"
# docker push ${REPO}:carts-${VERSION}
#
# catalogue and catalogue-db
echo ""
echo "Building catalogue and catalogue-db containers (${VERSION})"
cd catalogue
sh scripts/build.sh ${REPO} ${VERSION}
cd ..
echo ""
echo "pushing catalogue container to repository ${REPO}"
docker push ${REPO}:catalogue-${VERSION}
echo ""
echo "pushing catalogue container to repository ${REPO}"
docker push ${REPO}:catalogue-db-${VERSION}
#
# front-end
# echo ""
# echo "Building front-end container (version ${VERSION})"
# docker build -t ${REPO}:front-end-${VERSION} front-end/.
# echo ""
# echo "Pushing front-end container to repository ${REPO}"
# docker push ${REPO}:front-end-${VERSION}
# echo ""
#
# # front-end with BRUM
# echo ""
# echo "Building front-end appd brum container (version ${VERSION})"
# cd front-end-appd-brum
# sh scripts/build.sh ${REPO} ${VERSION}
# cd ..
# echo ""
# echo "Pushing front-end appd brum container to repository ${REPO}"
# docker push ${REPO}:front-end-appd-brum-${VERSION}
# echo ""
#
# # orders
# echo ""
# echo "Building orders container (version ${VERSION})"
# cd orders
# sh scripts/build.sh ${REPO} ${VERSION}
# cd ..
# echo ""
# echo "Pushing orders container to repository ${REPO}"
# docker push ${REPO}:orders-${VERSION}
# echo ""
#
# payment
# echo ""
# echo "Building payment container (version ${VERSION})"
# cd payment
# sh scripts/build.sh ${REPO} ${VERSION}
# cd ..
# echo ""
# echo "Pushing payment container to repository ${REPO}"
# docker push ${REPO}:payment-${VERSION}
# echo ""
#
# # queue-master
# echo ""
# echo "Building queue-master container (version ${VERSION})"
# cd queue-master
# sh scripts/build.sh ${REPO} ${VERSION}
# cd ..
# echo ""
# echo "Pushing queue-master container to repository ${REPO}"
# docker push ${REPO}:queue-master-${VERSION}
# echo ""
#
# # shipping
# echo ""
# echo "Building shipping container (version ${VERSION})"
# cd shipping
# sh ./scripts/build.sh ${REPO} ${VERSION}
# cd ..
# echo ""
# echo "Pushing shipping container to repository ${REPO}"
# docker push ${REPO}:shipping-${VERSION}
# echo ""
#
# # user
# echo ""
# echo "Building user container (version ${VERSION})"
# cd user
# docker build -t ${REPO}:user-${VERSION} .
# cd ..
# echo ""
# echo "Pushing user container to repository ${REPO}"
# docker push ${REPO}:user-${VERSION}
# echo ""
#
# # user-db
# echo ""
# echo "Building user-db container (version ${VERSION})"
# cd user/docker/user-db
# docker build -t ${REPO}:user-db-${VERSION} .
# cd ../../..
# echo ""
# echo "Pushing user-db container to repository ${REPO}"
# docker push ${REPO}:user-db-${VERSION}
# echo ""
#
# # load-test
# echo ""
# echo "Building load-test container (version ${VERSION})"
# docker build -t ${REPO}:load-test-${VERSION} load-test/.
# echo ""
# echo "Pushing load-test container (version ${VERSION}) to repository ${REPO}"
# docker push ${REPO}:load-test-${VERSION}
