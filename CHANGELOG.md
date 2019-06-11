# Changelog
## [Unreleased]

## [1.1.1_test] - 2019-06-03
### Added
- possibility to pass page and per_page parameters for all list commands
- list video encodings on video describe
- list status and createdAt for flip list encodings
- list createdAt for flip list factories
- list createdAt for flip list profiles
- list createdAt and status for flip list videos
- list createdAt and status for tts list projects

### Changed
- made command line interface more readable
- reformatted tables (list commands)
- set page to 1 if parameter per_page was passed but page not (list commands)

## [1.1.0] - 2019-05-27
### Added
- tts service client
- posibility to pass additional header key and it's value
### Changed
- moved credentials to home directory
- updated README.md


## [1.0.0] - 2019-05-20
### Added
- initial version