# What to contribute

If you want to contribute to this, please first take a look at the [open issues](https://github.com/UnseenWizzard/pocket-cli/issues). 

Feel free to grab any or create your own one and open a PR.

# Building and Testing

This project contains a [Makefile](./Makefile). 

To get started just use the default target:  

```shell
make coverage-check
```

Which will build and run tests. 

## Test Coverage

This project aims for at least 80% test coverage. When adding code, add tests. 

To check test coverage you can run:
```shell
make ch
```

## Running/Installing CLI Locally 

To build a runnable binary a pocket App consumer-key is needed. 

1. Create an App with `Add, Modify, Retrieve` permissions at https://getpocket.com/developer 
2. Copy the `consumer key` of that app
3. Run 
   
   ```shell
   .ci/release.sh {your consumer key} dev-$(git rev-parse --short HEAD) pocket-cli
   ```
4. Run the binary: `./pocket-cli`

# Conventional Commits

This project uses [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/).

With the following types (based on [angular conventions](https://github.com/angular/angular/blob/main/CONTRIBUTING.md#-commit-message-format)): 

* ci: Changes to anything Build and Test related
* docs: Documentation only changes
* feat: A new feature
* fix: A bug fix
* refactor: A code change that neither fixes a bug nor adds a feature
* test: Adding missing tests or correcting existing tests
* style: A code change only affecting style/formatting


When creating a commit, strive for [atomic commits](https://www.freshconsulting.com/insights/blog/atomic-commits/) 
with [good commit messages](https://cbea.ms/git-commit/).