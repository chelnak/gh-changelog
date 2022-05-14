# Make your changelogs ✨

[![ci](https://github.com/chelnak/gh-changelog/actions/workflows/ci.yml/badge.svg)](https://github.com/chelnak/gh-changelog/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/release/chelnak/gh-changelog.svg)](https://github.com/chelnak/gh-changelog/releases/latest)

An opinionated [GitHub Cli](https://github.com/cli/cli) extension for creating changelogs that adhere to the [keep a changelog](https://keepachangelog.com/en/1.0.0/) specification.

## What is supported?

`gh-changelog` is the tool for you if:

- You want to closely follow the [keep a changelog](https://keepachangelog.com/en/1.0.0/) specification
- You are using tags to mark releases
- You are following a pull-request workflow

## Installation and Usage

Before you start make sure that:

- GitHub Cli is [installed](https://cli.github.com/manual/installation) and [authenticated](https://cli.github.com/manual/gh_auth_login)
- You are inside a git repository
- The repository contains commits and has been pushed to GitHub

### Install

```bash
gh extension install chelnak/gh-changelog
```

### Upgrade

```bash
gh extension upgrade chelnak/gh-changelog
```

### Create a new changelog

```bash
gh changelog new
```

### Create a new changelog for an untagged version

```bash
gh changelog new --next-version v1.2.0
```

### View your changelog in the terminal

```bash
gh changelog show
```

## Configuration

Configuration for `gh changelog` can be found at `~/.config/gh-changelog/config.yaml`.
However some sensible defaults are provided to help you get off to a flying start.

```yaml
# Labels added here will be ommitted from the changelog
excluded_labels:
  - maintenance
# This is the filename of the generated changelog
file_name: CHANGELOG.md
# This is where labels are mapped to the sections in a changelog entry
# The possible sections are restricted to: Added, Changed, Deprecated,
# Removed, Fixed, Security.
sections:
  changed:
    - backwards-incompatible
  added:
    - fixed
    - enhancement
  fixed:
    - bug
    - bugfix
    - documentation
# When set to true, unlabelled entries will not be included in the changelog.
# By default they will be grouped in a section named "Other".
skip_entries_without_label: false
# Adds an unreleased section to the changelog. This will contain any qualifying entries
# that have been added since the last tag.
# Note: The unreleased section is not created when the --next-version flag is used.
show_unreleased: true
# If set to false, the tool will not check remotely for updates
check_for_updates: true
```

You can also override any setting using environment variables. When configured from the environment
properties are prefixed with `CHANGELOG`. For example, overriding `check_for_updates` might look
something like this:

```bash
export CHANGELOG_CHECK_FOR_UPDATES=false
```
