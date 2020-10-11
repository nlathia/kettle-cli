# operator

A command line tool for creating and deploying production machine learning functions.

These functions can be deployed on Google Cloud managed infrastructure as:

* [Cloud Functions](https://cloud.google.com/functions)
* [Cloud Run](https://cloud.google.com/run) containerized applications

Warning: this is a pre-release alpha version of this tool. Please send any bugs or feedback.

## Requirements

This project creates boiler plate code that uses:

* [pyenv](https://github.com/pyenv/pyenv)
* [pyenv-virtualenv](https://github.com/pyenv/pyenv-virtualenv)
* [gcloud](https://cloud.google.com/sdk/gcloud)

Specifically for Cloud Run:

* [Docker](https://docs.docker.com/get-docker/)

## Installing

Clone this repo and run:

```bash
❯ cd ~/src/github.com/operatorai/operator
❯ make install
```

## Usage

Set up the CLI tool using `operator init`.

```bash
❯ operator init
Use the arrow keys to navigate: ↓ ↑ → ← 
? Deployment type: 
  ▸ Google Cloud Function
    Google Cloud Run
```

Create a new deployment with `operator create`:

```bash
❯ operator create hello-world 
```

... and set it up:

```bash
❯ cd hello-world 
❯ make install # To create a pyenv-virtualenv
```

Launch it locally:

```bash
❯ make localhost
```

... and, when you're ready, deploy it!

```bash
❯ operator deploy .
```

## Notes

This tool has been built using the [Cobra Generator](https://github.com/spf13/cobra/blob/master/cobra/README.md#cobra-generator).

To add a new command:

```bash
❯ cobra add <command-name>
```
