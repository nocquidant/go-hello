# Release process

<!-- TOC -->

- [Release process](#release-process)
    - [From master branch](#from-master-branch)
    - [From another dev branch](#from-another-dev-branch)

<!-- /TOC -->

## From master branch

1. Commit all staging files and push to remote
1. Add a git tag: `git tag -a v1.0.0 -m "Version 1.0.0"`
1. Push the tag to remote: `git push origin v1.0.0`

Travis CI will then push the Docker image to DockerHub.

## From another dev branch

All releases are made from master.
