# Make your changelogs ✨

[![ci](https://github.com/chelnak/gh-changelog/actions/workflows/ci.yml/badge.svg)](https://github.com/chelnak/gh-changelog/actions/workflows/ci.yml) [![Release](https://img.shields.io/github/release/chelnak/gh-changelog.svg)](https://github.com/chelnak/gh-changelog/releases/latest)

An opinionated [GitHub Cli](https://github.com/cli/cli) extension for creating changelogs that adhere to the [keep a changelog](https://keepachangelog.com/en/1.0.0/) specification.

## Installation and Usage

```bash
# Install
gh extension install chelnak/gh-changelog
```

```bash
# Upgrade
gh extension upgrade chelnak/gh-changelog
```

Before you start make sure that:

- GitHub Cli is [installed](https://cli.github.com/manual/installation) and [authenticated](https://cli.github.com/manual/gh_auth_login)
- You are inside a git repository
- The repository contains commits and has been pushed to GitHub

Basic usage is simple. Just run the following command:

```bash
gh changelog
```
