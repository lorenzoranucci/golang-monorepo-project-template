# CI/CD strategy

The [ci-checks.yml](ci-checks.yml), [ci-build.yml](ci-build.yml) and [cd.yml](cd.yml) files are the templates for
application-specific CI/CD pipelines.

Every application-specific CI/CD pipeline should be triggered just when the application folder is changed and should use
the templates.

The CI phase takes care of executing checks, and it tests and builds the artifact (docker image).
The CD phase deploys the artifact on the appropriate environment (test, stage or live).

There are four phases:

- Pull request: when a pull request is created or modified.
- (optional) Test deploy: when a pull request is deployed on the test environment.
- Merge: when a pull request is merged on master.
- Release: when a tag is created on master.

## Pull request

The pull request phase is triggered when a pull request is created or modified.

Each time the CI is executed including every check, and just the docker image push is skipped.

## Test deploy

The test deploy is triggered manually through the GitHub Actions UI.

It executes the CI build and in case of success it runs the CD.

The CI produces the docker image, it tags the image with `test-next`, and it pushes the image to the ECR repository.

The CD deploys the docker image tagged with `test-next` on the test environment.

## Merge

The merge phase is triggered when a pull request is merged on master.

It executes the CI and in case of success it runs the CD.

The CI produces the docker image, it tags the image with `$LATEST_TAG-next`, and it pushes the image to the ECR repository.

The CD deploys the docker image tagged with `$LATEST_TAG-next` on the stage environment.

## Release

The release phase is triggered when a tag is created.

It executes the CD only, since the CI has been already done on merge.

The CD phase deploys the docker image on the live environment referencing it with the `$LATEST_TAG-next` tag.
