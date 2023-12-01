# Advent Of Code 2023

Scripts and code used to complete the advent of code of 2023, in python

## Development

You can find additional documentation in the [docs directory](docs).

### ADRs

If there is need for bigger technical, architectural, stack decisions, remember to fill
an [Architectural Decision Record (ADR)](https://adr.github.io/) in the [docs/adr](docs/adr), with
file name formatted `[yyyy-mm-dd]_[topic].md`



## Remote Deployment

### Docker

This repo contains one single and simple Dockerfile, building a python image, installing the
requirements
through pip, then running a simple fastAPI app with uvicorn.

### GCP + CI CD Pipeline w/ GCP

There are two cloudbuild yaml files, one for building the image (on dev), one for deploying a new
revision (on dev/staging/production). The pipeline should look like this:

- Developer codes and tests a new feature in a feature branch, forked from the `main` branch.
- Developer opens a pull request to `main` (with changes described in
  the `## Unreleased` section of the [changelog](CHANGELOG.md). Other developers review, request
  changes, etc.
- Once the PR is accepted, developer can merge it into dev.
- He then moves the `## Unreleased` section to a new section with the new version and date, updates
  the version in the [pyproject.toml](pyproject.toml), then commits and tags this commit with the
  new version `vX.Y.Z` (optionally with a comment after a dash `-`), and pushes the commit and the
  tag to origin. (Fast forward a tag works as well)
- The new tag on origin triggers a GCP Cloud Build Trigger (which is connected to the origin
  repository in BitBucket).
  This trigger clones the repository, then follow the steps described in
  the [cloudbuild.yaml](deployment/cloudbuild.yaml) (build, test and save image).
  There are substitution
  variables, some are automatically given by GCP (no `_` as prefix), the other ones have to be
  specified (prefix `_`). The
  latter one are either defined inside the yaml file at the end under key `substitutions:`, or
  inside the GCP Console,
  when editing the Cloud Build Trigger (the latter ones overwrite the former ones). The
  [cloudbuild.yaml](deployment/cloudbuild.yaml) specifies what
  secrets and environment variables are to be used. Those can be overwritten or completed in the
  GCP interface though.
- To deploy an actual new revision, you will have to manually trigger the `deploy` trigger, no
  matter the deploy environment you are in right now (dev, stage, prod). Those
  triggers follow the steps in the  [clouddeploy.yaml](deployment/clouddeploy.yaml), simply pulling
  the already built image from the (dev) artifacts registry and updating the running
  service with a new revision.

## Local setup

### Uvicorn directly in terminal

- make sure you have an active virtual environment using python >= 3.11
  
- set the environment variable `DEPLOY_ENVIRONMENT` to `LOCAL` or write `DEPLOY_ENVIRONMENT=LOCAL`
  in the `.env` file
- install the projects requirements with `pip install -r requirements.txt`
- in same terminal, run `uvicorn src.main:app --reload` (reload makes the server update each time
  you change and save
  the code)
- You should see something
  like `INFO:     Uvicorn running on http://127.0.0.1:8000 (Press CTRL+C to quit)` in the
  terminal, so try to open the [docs page](http://127.0.0.1:8000/docs) to see if it
  works: http://127.0.0.1:8000/docs

### Testing

If all requirements are installed, you can simply run `pytest` at the root of the repository, it
will scan the directory
and find the directory `tests`, prepare the fixtures and configurations specified in
its `conftest.py`, then run all
functions with prefix `test_`.   
You can specify a subset of test by specifying the path to a directory or a file (
e.g. `pytest src/tests/test_routes.py`)  
You can as well use the IDE tools to run the tests (on pycharm, the "run" green triangle icons next
to each test or
class of tests, or after right-clicking on a file/directory and run). You can as well add some "
Configurations" of type
Python Testing > pytest for tests you might run often. This way you get a better overview of which
test failed than just
the cli with terminal command pytest.

## Contact

Responsible: Nils Cichy, email: nils.cichy@TODO
Maintainers : Nils Cichy,
