# kettle

Kettle is a command line tool for creating and deploying machine learning pipelines and microservices, starting from from best-in-class templates.

This CLI has two primary commands:

* `kettle create <name>` creates a directory containing all the boiler plate code that you need to get going. 
* `kettle deploy <path>` deploys the code in that directory to the cloud. Deploy currently supports serverless functions on AWS and GCP.

### Templates

Kettle supports three types of templates:

1. Templates that are already on your computer, at a given path.
2. Templates that are git repositories
3. Templates that are in the `kettle-templates` [repository](https://github.com/operatorai/kettle-templates); browse that repo's [README](https://github.com/operatorai/kettle-templates/blob/main/README.md) to see the templates that it contains spanning AWS Lambda, GCP Functions, and GCP Run.

## Installing with brew

You can install `kettle` using `brew` and [the operatorai tap](https://github.com/operatorai/homebrew-tap).

```bash
❯ brew tap operatorai/tap
❯ brew install kettle-cli

# You can see that this works by running
❯ kettle version
```

## Usage

Here's an example that takes you from a template to a deployed AWS Lambda.

### Example from kettle-templates

In the example below, we use the [pyenv-aws-lambda](https://github.com/operatorai/kettle-templates/tree/main/pyenv-aws-lambda) template in the [kettle-templates](https://github.com/operatorai/kettle-templates) repository. Since we're using a `kettle-templates` template, we just need to use `kettle create <name>`, where `<name>` is the directory name in the templates repo:

```bash
❯ kettle create pyenv-aws-lambda
Project name: hello-world

✅  Created:  <path>/hello-world
```

This will prompt you for a project name, and will then create that directory and add all the boiler plate you need to get going. This particular template comes with a `Makefile`, that we can use to set up the local environment:

```bash
❯ cd hello-world

❯ make install
```

## Kettle deploy

Kettle `deploy` is the command to deploy your project as a serverless function. It currently supports:

### AWS Lambdas

You must have the [aws cli](https://aws.amazon.com/cli/) installed.

For Python, `kettle` supports Lambdas where Python is managed with `pyenv` or `conda`.

### Google Cloud Functions

You must have the [gcloud](https://cloud.google.com/sdk/gcloud) SDK installed. You also need to have enabled the Cloud Functions API in the GCP console.

### Google Cloud Run

You must have the [gcloud](https://cloud.google.com/sdk/gcloud) SDK installed, and optionally [Docker](https://docs.docker.com/get-docker/) to build and run Cloud Run containerized applications locally. You also need to have enabled the Cloud Run API in the GCP console.

## Bug Reports

Please report any bugs or issues to me (neal.lathia@gmail.com) or by raising an issue in this repo.
