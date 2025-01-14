# Change Log

All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

This file itself is based on [Keep a CHANGELOG](https://keepachangelog.com/en/0.3.0/).

## [Unreleased]

### Changed
- Change license to Apache License Version 2.0.

## [2.1.1] - 2021-09-07

### Changed
- Fix lvmd is not previleged in deploying with Helm (#358)
- appVersion was changed to 0.9.2.

### Misc
- duplicate label causes YAML parsing errors (#351)

### Contributors
- @faruryo
- @khrisrichardson

## [2.1.0] - 2021-08-20

### Changed
- Make allowedHostPaths property of node's PSP configurable (#347)
- Update appVersion to 0.9.0 in Chart.yaml (#348)

### Contributors
- @debackerl

## [2.0.3] - 2021-08-19

### Removed
- Remove kustomize. (#336)

### Added
- Helm Chart: Support custom clusterIP for scheduler deployment (#346)

### Contributors
- @d-kuro
- @yuseinishiyama

## [2.0.1] - 2021-07-27

This is the first release.

[Unreleased]: https://github.com/topolvm/topolvm/compare/topolvm-chart-v2.1.1...HEAD
[2.1.1]: https://github.com/topolvm/topolvm/compare/topolvm-chart-v2.1.0...topolvm-chart-v2.1.1
[2.1.0]: https://github.com/topolvm/topolvm/compare/topolvm-chart-v2.0.3...topolvm-chart-v2.1.0
[2.0.3]: https://github.com/topolvm/topolvm/compare/topolvm-chart-v2.0.1...topolvm-chart-v2.0.3
[2.0.1]: https://github.com/topolvm/topolvm/releases/tag/topolvm-chart-v2.0.1
