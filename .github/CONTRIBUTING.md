# Contributing to Mify

## Getting Started

### Ways to Contribute

- Reporting bugs.
- Suggest product enhancements.
- Fixing bugs.
- Implementing a requested enhancement in issues.
- Improving documentation and guides.

## Building Mify

Mify is a Go project, so you'll need [Go](https://golang.org/) installed.
Then, clone this repository and run `make build`. After this command you'll
have a working `mify` executable in the root directory of the repository.

## Making changes

The first step is to fork mify repository and clone your fork locally.
Afterwards, navigate to your cloned fork directory and prepare it:

1. Add upstream origin: `git remote add upstream <github url of upstream mify repo>`
2. Fetch upstream: `git fetch --all --tags`
3. Create your feature branch: `git checkout -b new-feature upstream/main`
4. Make changes and commit them
5. Push changes to the fork when ready: `git push -u origin new-feature`

Then you'll be able to create a PR to upstream repository.

## Testing

Right now there is script to test common use case, to run it call `./scripts/init.sh`
from the repository root.

## Documentation

Documentation is available under [./docs](https://github.com/mify-io/mify/tree/main/docs).
It is built on Docusaurus and is hosted at https://mify.io/docs
