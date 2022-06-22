# Application Source Code

Please use the `container-build.sh` script to build the neccessary Docker images and upload these to the repository of your choice.  `version` is used to tag the created images. `repo` is the path to the docker image repository you are using (i.e. dockerhub, ECR etc.).  If neccessary please authenticate to the Docker repository first.

```sh
./container-build.sh {version} {repo}
```

## Note:
None of the other scripts present have been updated. This includes the front-end and catalogue components modified to present a skinned "cool cars" shop.
