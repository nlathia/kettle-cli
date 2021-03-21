# kettle

⚠️ This repo is currently being renamed (it used to be called `operator`) and it has not been fully tested with its new name yet.

Kettle is a command line tool for creating and deploying machine learning pipelines and microservices, starting from from best-in-class templates.

This CLI has two primary commands:

* `kettle create <name>` creates a directory containing all the boiler plate code that you need to get going. 
* `kettle deploy <path>` deploys the code in that directory to the cloud.

### Templates

Kettle supports three types of templates:

1. Templates that are already on your computer, at a given path.
2. Templates that are git repositories
3. Templates that are in the `kettle-templates` monorepo.

## Installing with brew

⚠️ This is a work in progress.

You will be able to install `kettle` using `brew` and [the operatorai tap](https://github.com/operatorai/homebrew-tap).

```bash
❯ brew tap operatorai/tap
❯ brew install kettle
```

## Usage

Create a new deployment by pointing `kettle create` to a template. In the example below, we use the [pyenv-aws-lambda](https://github.com/operatorai/kettle-templates/tree/main/pyenv-aws-lambda) template in the [kettle-templates](https://github.com/operatorai/kettle-templates) repository.

```bash
❯ kettle create pyenv-aws-lambda
Directory name: hello-world

✅  Created:  ~/hello-world
```

This will prompt you for a directory name, and will then create that directory and add all the boiler plate you need to get going. 

## Kettle deploy

Kettle `deploy` currently supports the following

### AWS Lambdas

You must have the [aws cli](https://aws.amazon.com/cli/) installed.

### Google Cloud Functions & Google Cloud Run

You must have the [gcloud](https://cloud.google.com/sdk/gcloud) SDK installed, and optionally [Docker](https://docs.docker.com/get-docker/) to build and run Cloud Run containerized applications locally. You also need to have enabled the Cloud Functions (or Cloud Run) API in the GCP console.

## Bug Reports

Please report any bugs or issues to me (neal.lathia@gmail.com) or by raising an issue in this repo.
