> ❗ While this is a working CLI this is largely a project for me to play around with some libraries and tools.
> 
> Notably this only really used/tested on Linux

# Pocket CLI

This is a simple CLI to access your [Pocket](https://getpocket.com) reading list.

## Installation

Download the [latest release](https://github.com/UnseenWizzard/pocket-cli/releases/latest) and put it somewhere on your path.
(e.g. move it to `/usr/local/bin/` on Linux)

Or see [CONTRIBUTING.md](./CONTRIBUTING.md#runninginstalling-cli-locally) for details on building from source.

## Authorize App

When first running you need to authorize `pocket-cli` to access your [Pocket](https://getpocket.com) Account.

Do so by running: 

```shell
pocket-cli login
```

## Access Your Reading List

After authorizing the App access your reading list with:

```shell
pocket-cli list
```