# Application Source Code

Please use the `container-build.sh` script to build the neccessary Docker images and upload these to the repository of your choice.  `version` is used to tag the created images. `repo` is the path to the docker image repository you are using (i.e. dockerhub, ECR etc.).  If neccessary please authenticate to the Docker repository first.

```sh
./container-build.sh {version} {repo}
```

## Notes:
- Originally a customised image was created to run the AppDynamics Database monitoring agent.  This has been superceded with the official AppD DB monitoring container.  The original agent image has been moved the /old directory.
- None of the other scripts present have been updated. This includes the front-end and catalogue components modified to present a skinned "cool cars" shop.
- In order to integrate the AppD Go SDK library, the User, Payment and Catalogue components have been modified.  The Dockerfile has been modified to include the latest GO SDK library and also to upgrade certain Go components (gorilla mux) in the original Weaveworks build to support the recommended AppD integration points.  See the components `/scripts/build.sh` script for more details.
