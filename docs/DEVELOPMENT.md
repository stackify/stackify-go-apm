# Go APM Development Guide

## Development Setup

1. Install Go 1.15

2. Clone Repository

    http: `git clone https://<user>@bitbucket.org/stackify/stackify-go-apm.git`
    ssh:  `git clone git@bitbucket.org:stackify/stackify-go-apm.git`

3. Install Dependencies
    `$ go mod tidy`


## Run Test
    `$ ./tests.sh`

## Build
    ```
    $ go mod tidy
    ```

## Publish
    - merge develop to master
    - create tag from master branch
    - push tag to github repo
