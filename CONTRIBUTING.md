# How to Contribute

Crust projects are [Apache 2.0 licensed](LICENSE) and accept contributions via
GitHub pull requests.  This document outlines some of the conventions on
development workflow, commit message formatting, contact points and other
resources to make it easier to get your contribution accepted.

# Certificate of Origin

By contributing to this project you agree to the Developer Certificate of
Origin (DCO). This document was created by the Linux Kernel community and is a
simple statement that you, as a contributor, have the legal right to make the
contribution. See the [DCO](DCO) file for details.

# Getting Started

- Fork the repository on GitHub
- Read the [Wiki](https://github.com/crusttech/crust/wiki) for build and test instructions
- Play with the project, submit bugs, submit patches!

## Contribution Flow

This is a rough outline of what a contributor's workflow looks like:

- Create a topic branch from where you want to base your work (usually master).
- Make commits of logical units.
- Make sure your commit messages are in the proper format (see below).
- Push your changes to a topic branch in your fork of the repository.
- Make sure the tests pass, and add any new tests as appropriate.
- Submit a pull request to the original repository.

Thanks for your contributions!

### Format of the Commit Message

We follow a rough convention for commit messages that is designed to answer
two questions: what kind of change it was, and what changed. The subject line
should feature all of this information.

We try to structure the commit messages as follows:

```
verb(location): description of change
location: verb + rest of the description
human readable sentence of changes
```

We are not terribly strict against a particular commit message format, but
the general rule is to avoid non-descriptive single-word commit messages like "fix".

The commit message body may feature a bullet list of multiple changes. This
is usual when a larger scope change has been made, or when squashing commits.

The first line is the subject and should be no longer than 70 characters, the
second line is always blank, and other lines should be wrapped at 80 characters.
This allows the message to be easier to read on GitHub as well as in various
git tools.
