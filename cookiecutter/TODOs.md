# Post Generation TODOs


There is still some steps to do before having your app deployed in GCP and accessible to the whole
internet. Here in the sub-chapters you can find explanation about each step, you can as well search
for TODO-CONFIG project-wide to see what you have to change in the source code.

## Introduction / Explanation

This should help you understand the why of the next steps.  
The goal is at the end to have a setup that:

- makes your code accessible for anyone in the internet
- allows easy deploying to a new version
- allows easy rollback to a previous version
- the same version of the code can be deployed on different environments (DEV, STAGE, PROD ...)

Now, we have some restrictions:

- We host the code (our git repositories) in BitBucket
- We use the GCP platform

So how do the different hold together to allow all these points?

- The code should be available from GCP: that's why we will mirror the repo via Cloud Source
  Repositories
- The code should not be execute "as is", but be wrapped in a Docker image, so we need a Dockerfile
  describing how to build the image. We will put this file in the git repo under the
  directory `deployment`. (Docker image will allow storing versions of the logic, be able to deploy
  on any OS, will be essential for scalability later)
- Now we don't want the image to be built manually by developers, but rather allow developers to "
  declare" a git commit as build-ready and GCP should do the rest. For this we use GCP Cloud Build
  Triggers, that can "listen" to a repo changes. In our case we want to react on git tags push (on
  any branch). We will restrict the labels triggering the build to those following SemVer
  structure, to allow using other tags without interference, but also to have a nice git history
  showing all versions and providing the version information to the build trigger.
- Well it is clear "when" the trigger should build, but not "what" this trigger should do. GCP
  provides a specific file format (cloudbuild.yaml) to help declare the steps in a build. We will
  keep those files in the `deployment` directory as well, to allow version controlling of it. This
  cloudbuild will describe those steps: build an image basing on the Dockerfile, (optionally run a
  container from this image and run its tests), push the image to Artifact Registry then deploy a
  new "revision" basing on this image. You just heard two new things here:
- First, the Artifact Registry. That's the place where we will
  store the built images (long-term). Those images are tagged with the semantic version from the
  tag, so any time if you want to rollback to a version, you can deploy from the one you wanted.
  Additionally, for a specific code version, we only want to build the image once, but deploy on
  multiple environments, so storing the images will allow re-using them for Stage and Prod for
  instance. Lastly, in case of debugging a bug, you can download the image and run a container
  locally.
- Second: a new "revision". It is relevant to understand, that you don't deploy a "container", but
  a "revision" in the Cloud Run. A revision is based on an image, provides an URL to make the
  image's logic accessible, but internally will have multiple containers (called "instances")
  running and will handle
  redirecting the traffic to one of these containers. (Actually it will 
  start or stop more containers depending on the traffic.) 

## Adapt the deployment files

- You probably want to add some environment variables in the cloudbuild.yaml files, add some
  substitutions, maybe change the Dockerfile. You can find these files in the
  [deployment folder](deployment).

## Connect the local repo to Bitbucket

- Either use a gitGUI tool (like gitkracken) to handle the connection to a new origin and creation
  of a repo there
- Or do it yourself through git cli (`git remote add`) ,
- Or the naive method (not recommended as you lose a bit of git history): create a new repo on
  Bitbucket, clone it, then copy everything you generated from this cookiecutter into it.

Create a new repository in GCP's
with the name of your app.

## Artifact Registry

- Make sure there is an image repository you want to use in
  the  [Artifact Registry](https://console.cloud.google.com/artifacts).
  If none is adequate, create one.
    - name: YOUR_APP_NAME.
      (Underscores are not allowed, so replace them with hyphens `-`)
    - format: Docker
    - mode: Standard
    - Region: europe-west3 (Frankfurt)
- Copy its name, add it in the cloudbuild.yaml files as value to the
  substitution `__IMAGE_REPOSITORY`

## Secrets

- If your app needs some new secrets, create them in
  the [Secret Manager](https://console.cloud.google.com/security/secret-manager?)
- Specify them in the cloudbuild.yaml.
    - For build time secrets (like dependency license keys etc.), you want them to be available as
      secrets in the trigger and in the build step (see top level field `availableSecrets` and the
      step field `secretEnv`)
    - For run time secrets, simply add the secret names in the `--set-secrets` param of the deploy
      step

## Prepare your requests to the admin

Some tasks have to be done by a devops admin, so prepare the information you need before asking
him:

- What deploy environments do you want for your service?
- What APIs will your service be using? what permissions does your run service account need?
    - MongoDB? what database? read or write?
    - CloudSQL?
    - PubSubs? what topic? publish or subscribe?
    - Buckets? which ones? read or write?
    - Secrets?
    - others?

## Ask an admin to connect the Bitbucket repo to GCP

Kindly ask a devops admin (happtiq or ELOS) to do these steps for you:

- Go to Bitbucket, Log in with a _devops_ (don't use personal account) Bitbucket account
- Go to GCP: [Google Cloud Source Repositories](https://source.cloud.google.com/repo/connect)
- click "Add repository"
- choose "Connect external repository"
- choose the correct GCP project and bitbucket as provider
- Choose the new repo in the list
- That will take a few seconds. That's it, you now can use it when creating the trigger.

## Ask an admin to Create a run Service Account

- Tell to the devops admin which permissions the service account will need
- He will create a service account, give him permissions and tell you its email address.
- You will then add this email as value of `gcloud run deploy --service-account` in the
  cloudbuild.yaml (probably as substitution value in the trigger)

## Create Triggers in GCP

### Build Trigger (on DEV)

Here you want to create the [GCP Trigger](https://console.cloud.google.com/cloud-build/triggers)
that will react to pushing a new tag in git, in order to
build a docker image, run the tests on it and save it into
the [GCP Artifacts Registry](https://console.cloud.google.com/artifacts). It will also deploy a new
revision.

- name -> `<your_service_name>-build`
- region: europe-west3
- (optional) Tags: `build` `test`
- Event: push new tag
- Source: 1st gen -> your repo
- regex:
  ```regexp
  ^v(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(-(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))?$
  ```
  > this is a variant of
  the [official semantic versioning regular expression](https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string)
  accepting MAJOR.MINOR.PATCH-prerelease (without + metadata as `+` is not compatible with gcp
  buildtags), with prerelease being optional.  
  > some examples: `v0.0.1`, `v2.1.23`,  `v2.1.23-1` , `v2.1.23-alpha1`
- Configuration:
    - Type: "Cloud Build configuration file (yaml or json)"
    - Location: "Repository", `deployment/cloudbuild.yaml`
- Substitution Variables:
    - Any substitutions the cloudbuild.yaml uses but does not define (or if you want to overwrite
      some)
- service account: leave it empty, so it uses the default compute service account. (we will change
  this once we have the new gcp project structure)

### Deploy Trigger (on STAGE+PROD)

Here you want to create this trigger for each remote environment (assuming you have a different GCP
project for dev, stage and prod). It will pull the image from the dev environment (where you built
the image with the previous trigger) and deploy a new revision based on this image.
This will be a manual trigger, that takes a TAG_NAME as parameter (to know which image to pull)

- name -> `<your_service_name>-deploy`
- region: europe-west3
- (optional) Tags: `deploy`
- Event: Manual invocation
- Source: 1st gen -> Your repo
- Configuration
    - Type: "Cloud Build configuration file (yaml or json)"
    - Location: "Repository",  `deployment/clouddeploy.yaml`
- Substitution Variables:
    - `_DEPLOY_ENVIRONMENT`: `DEV` or `STAGE` or `PROD` (depending on the GCP project)
    - `_MIN_INSTANCES`: 0 - 1 (depending on the GCP project, you might not want dev to have a
      processor idle, but on stage yes?)
    - `_MAX_INSTANCES` (depending on the GCP project as well)
    - Any substitutions the cloudbuild.yaml uses but does not define (or if you want to overwrite
      some)
- service account: leave it empty, so it uses the default compute service account. (we will change
  this once we have the new gcp project structure)



## Last Steps

- Run the trigger by git-label the first version `v0.0.1` and pushing the label to origin
    - You should see the trigger executing in
      the [Cloud Build History](https://console.cloud.google.com/cloud-build/builds)
- While the first version is building, use the time to write some information about your new app
  here in this readme and delete this todo-config list
- Once your revision is successfully deployed, go to
  the [Cloud Run page](https://console.cloud.google.com/run), find your revision (name
  matches `_SERVICE_NAME` from the cloudbuild.yaml), grab the url and add it
  to [this project's metadata in the pyproject.toml](pyproject.toml). (The base url redirects to
  the docs page)
    - If the url is not available, and the service only allows authenticated ingress, contact a
      devops admin. -> [relevant link](https://cloud.google.com/run/docs/authenticating/public)
