<!-- markdownlint-disable MD024 -->
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/) and this project adheres to [Semantic Versioning](http://semver.org).

## [v0.12.1](https://github.com/chelnak/gh-changelog/tree/v0.12.1) - 2023-04-12

[Full Changelog](https://github.com/chelnak/gh-changelog/compare/v0.12.0...v0.12.1)

### Fixed

- Fix config initialization [#129](https://github.com/chelnak/gh-changelog/pull/129) ([chelnak](https://github.com/chelnak))

## [v0.12.0](https://github.com/chelnak/gh-changelog/tree/v0.12.0) - 2023-04-11

[Full Changelog](https://github.com/chelnak/gh-changelog/compare/v0.11.0...v0.12.0)

### Added

- Allow local configs [#119](https://github.com/chelnak/gh-changelog/pull/119) ([chelnak](https://github.com/chelnak))
- Add a Markdown parser [#117](https://github.com/chelnak/gh-changelog/pull/117) ([chelnak](https://github.com/chelnak))

### Fixed

- Handle Pre-Releases [#126](https://github.com/chelnak/gh-changelog/pull/126) ([chelnak](https://github.com/chelnak))

## [v0.11.0](https://github.com/chelnak/gh-changelog/tree/v0.11.0) - 2022-12-01

[Full Changelog](https://github.com/chelnak/gh-changelog/compare/v0.10.1...v0.11.0)

### Added

- Convert changelog datastructure [#112](https://github.com/chelnak/gh-changelog/pull/112) ([chelnak](https://github.com/chelnak))

### Fixed

- Fix usage on repositories without tags [#114](https://github.com/chelnak/gh-changelog/pull/114) ([chelnak](https://github.com/chelnak))

### Other

- Fix markown formatting [#113](https://github.com/chelnak/gh-changelog/pull/113) ([chelnak](https://github.com/chelnak))

## [v0.10.1](https://github.com/chelnak/gh-changelog/tree/v0.10.1) - 2022-10-20

[Full Changelog](https://github.com/chelnak/gh-changelog/compare/v0.10.0...v0.10.1)

### Fixed

- Fix tags append [#109](https://github.com/chelnak/gh-changelog/pull/109) ([chelnak](https://github.com/chelnak))

## [v0.10.0](https://github.com/chelnak/gh-changelog/tree/v0.10.0) - 2022-10-14

[Full Changelog](https://github.com/chelnak/gh-changelog/compare/v0.9.0...v0.10.0)

### Added

- Rename text logger [#105](https://github.com/chelnak/gh-changelog/pull/105) ([chelnak](https://github.com/chelnak))
- Better logging control [#102](https://github.com/chelnak/gh-changelog/pull/102) ([chelnak](https://github.com/chelnak))
- Refactor builder & changelog in to pkg [#101](https://github.com/chelnak/gh-changelog/pull/101) ([chelnak](https://github.com/chelnak))
- Scoped changelogs [#100](https://github.com/chelnak/gh-changelog/pull/100) ([chelnak](https://github.com/chelnak))

## [v0.9.0](https://github.com/chelnak/gh-changelog/tree/v0.9.0) - 2022-10-07

[Full Changelog](https://github.com/chelnak/gh-changelog/compare/v0.8.1...v0.9.0)

### Added

- Help and error UX improvements [#92](https://github.com/chelnak/gh-changelog/pull/92) ([chelnak](https://github.com/chelnak))
- Refactoring and code hygiene [#88](https://github.com/chelnak/gh-changelog/pull/88) ([chelnak](https://github.com/chelnak))
- Remove internal/pkg [#86](https://github.com/chelnak/gh-changelog/pull/86) ([chelnak](https://github.com/chelnak))
- Add a new excluded label default [#82](https://github.com/chelnak/gh-changelog/pull/82) ([chelnak](https://github.com/chelnak))

### Fixed

- Properly handle a repo with no tags [#95](https://github.com/chelnak/gh-changelog/pull/95) ([chelnak](https://github.com/chelnak))

## [v0.8.1](https://github.com/chelnak/gh-changelog/tree/v0.8.1) - 2022-06-07

[Full Changelog](https://github.com/chelnak/gh-changelog/compare/v0.8.0...v0.8.1)

### Fixed

- Fix lexer in PrintYAML method [#80](https://github.com/chelnak/gh-changelog/pull/80) ([chelnak](https://github.com/chelnak))

## [v0.8.0](https://github.com/chelnak/gh-changelog/tree/v0.8.0) - 2022-05-20

[Full Changelog](https://github.com/chelnak/gh-changelog/compare/v0.7.0...v0.8.0)

### Added

- Colorize config output [#77](https://github.com/chelnak/gh-changelog/pull/77) ([chelnak](https://github.com/chelnak))
- Simplify config printing [#70](https://github.com/chelnak/gh-changelog/pull/70) ([chelnak](https://github.com/chelnak))
- Add a command to view current config [#63](https://github.com/chelnak/gh-changelog/pull/63) ([chelnak](https://github.com/chelnak))
- Enable configuration from environment [#61](https://github.com/chelnak/gh-changelog/pull/61) ([chelnak](https://github.com/chelnak))

### Fixed

- Validate next version [#76](https://github.com/chelnak/gh-changelog/pull/76) ([chelnak](https://github.com/chelnak))
- Handle orphaned commits [#74](https://github.com/chelnak/gh-changelog/pull/74) ([chelnak](https://github.com/chelnak))

## [v0.7.0](https://github.com/chelnak/gh-changelog/tree/v0.7.0) - 2022-05-14

[Full Changelog](https://github.com/chelnak/gh-changelog/compare/v0.6.1...v0.7.0)

### Added

- Lowercase section names [#52](https://github.com/chelnak/gh-changelog/pull/52) ([chelnak](https://github.com/chelnak))
- Remove additional newlines in markdown [#51](https://github.com/chelnak/gh-changelog/pull/51) ([chelnak](https://github.com/chelnak))
- Rework configuration [#47](https://github.com/chelnak/gh-changelog/pull/47) ([chelnak](https://github.com/chelnak))

### Fixed

- Ensure that Keep a Changelog format is followed [#53](https://github.com/chelnak/gh-changelog/pull/53) ([chelnak](https://github.com/chelnak))

## [v0.6.1](https://github.com/chelnak/gh-changelog/tree/v0.6.1) - 2022-05-08

[Full Changelog](https://github.com/chelnak/gh-changelog/compare/v0.6.0...v0.6.1)

### Fixed

- Remove Println [#44](https://github.com/chelnak/gh-changelog/pull/44) ([chelnak](https://github.com/chelnak))

## [v0.6.0](https://github.com/chelnak/gh-changelog/tree/v0.6.0) - 2022-05-08

[Full Changelog](https://github.com/chelnak/gh-changelog/compare/v0.5.0...v0.6.0)

### Added

- Add update check [#41](https://github.com/chelnak/gh-changelog/pull/41) ([chelnak](https://github.com/chelnak))
- Break up changelog package [#38](https://github.com/chelnak/gh-changelog/pull/38) ([chelnak](https://github.com/chelnak))

### Fixed

- Fix error messages [#37](https://github.com/chelnak/gh-changelog/pull/37) ([chelnak](https://github.com/chelnak))

## [v0.5.0](https://github.com/chelnak/gh-changelog/tree/v0.5.0) - 2022-05-07

[Full Changelog](https://github.com/chelnak/gh-changelog/compare/v0.4.0...v0.5.0)

### Added

- Implementation of next-version and unreleased [#35](https://github.com/chelnak/gh-changelog/pull/35) ([chelnak](https://github.com/chelnak))
- Refactor & (some) tests [#34](https://github.com/chelnak/gh-changelog/pull/34) ([chelnak](https://github.com/chelnak))

## [v0.4.0](https://github.com/chelnak/gh-changelog/tree/v0.4.0) - 2022-04-26

[Full Changelog](https://github.com/chelnak/gh-changelog/compare/v0.3.1...v0.4.0)

### Added

- Migrate to GitHubs v4 API [#29](https://github.com/chelnak/gh-changelog/pull/29) ([chelnak](https://github.com/chelnak))
- Better config [#28](https://github.com/chelnak/gh-changelog/pull/28) ([chelnak](https://github.com/chelnak))
- Better config [#27](https://github.com/chelnak/gh-changelog/pull/27) ([chelnak](https://github.com/chelnak))

### Fixed

- Clarify functionality in README [#25](https://github.com/chelnak/gh-changelog/pull/25) ([chelnak](https://github.com/chelnak))

## [v0.3.1](https://github.com/chelnak/gh-changelog/tree/v0.3.1) - 2022-04-20

[Full Changelog](https://github.com/chelnak/gh-changelog/compare/v0.3.0...v0.3.1)

### Fixed

- Fixes root commit reference [#24](https://github.com/chelnak/gh-changelog/pull/24) ([chelnak](https://github.com/chelnak))

## [v0.3.0](https://github.com/chelnak/gh-changelog/tree/v0.3.0) - 2022-04-20

[Full Changelog](https://github.com/chelnak/gh-changelog/compare/v0.2.1...v0.3.0)

### Added

- Turbo boost! [#20](https://github.com/chelnak/gh-changelog/pull/20) ([chelnak](https://github.com/chelnak))

### Fixed

- Set longer line length for md render [#18](https://github.com/chelnak/gh-changelog/pull/18) ([chelnak](https://github.com/chelnak))

## [v0.2.1](https://github.com/chelnak/gh-changelog/tree/v0.2.1) - 2022-04-18

[Full Changelog](https://github.com/chelnak/gh-changelog/compare/v0.2.0...v0.2.1)

### Fixed

- Fix full changelog link [#14](https://github.com/chelnak/gh-changelog/pull/14) ([chelnak](https://github.com/chelnak))

## [v0.2.0](https://github.com/chelnak/gh-changelog/tree/v0.2.0) - 2022-04-15

[Full Changelog](https://github.com/chelnak/gh-changelog/compare/v0.1.0...v0.2.0)

### Added

- Add show command [#12](https://github.com/chelnak/gh-changelog/pull/12) ([chelnak](https://github.com/chelnak))
- Implement better errors [#9](https://github.com/chelnak/gh-changelog/pull/9) ([chelnak](https://github.com/chelnak))

## [v0.1.0](https://github.com/chelnak/gh-changelog/tree/v0.1.0) - 2022-04-15

[Full Changelog](https://github.com/chelnak/gh-changelog/compare/42d4c93b23eaf307c5f9712f4c62014fe38332bd...v0.1.0)
