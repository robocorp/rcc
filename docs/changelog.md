# rcc change log

## v9.5.4 (date: 2.3.2021)

- Updated rpaframework to version 7.6.0 in templates

## v9.5.3 (date: 2.3.2021)

- added `--interactive` flag to `rcc task run` command, so that developers
  can use debuggers and other interactive tools while debugging

## v9.5.2 (date: 25.2.2021)

- bug fix: now cloning sources are not removed during --liveonly action,
  even when that source seems to be invalid
- changed timeline to use percent (not permilles anymore)
- minor fix on env diff printout

## v9.5.1 (date: 25.2.2021)

- now also printing environment differences when live is dirty and base
  is not, just before restoring live from base

## v9.5.0 (date: 25.2.2021)

- added support for detecting environment corruption
- now dirhash command can be used to compare environment content

## v9.4.4 (date: 24.2.2021)

- fix: added panic protection to telemetry sending, this closes #13
- added initial support for execution timeline tracking

## v9.4.3 (date: 23.2.2021)

- added generic reading and parsing diagnostics for JSON and YAML files

## v9.4.2 (date: 23.2.2021)

- fix: marked --report flag required in issue reporting
- added account-email to issue report, as backup contact information

## v9.4.1 (date: 17.2.2021)

- added conda.yaml diagnostics (initial take)
- made `rcc env variables` to be not silent anymore
- log level changes in environment creation
- env creation workflow has now 6 steps, added identity visibility

## v9.4.0 (date: 17.2.2021)

- added initial robot diagnostics (just robot.yaml for now)
- integrated robot diagnostics into configuration diagnostics (optional)
- integrated robot diagnostics to issue reporting (optional)
- fix: windows paths were wrong; "bin" to "usr" change

## v9.3.12 (date: 17.2.2021)

- introduced 48 hour delay to recycling temp folders (since clients depend on
  having temp around after rcc process is gone); this closes #12

## v9.3.11 (date: 15.2.2021)

- micromamba upgrade to 0.7.14
- made process fail early and visibly, if micromamba download fails

## v9.3.10 (date: 11.2.2021)

- Windows automation made environments dirty by generating comtypes/gen
  folder. Fix is to ignore that folder.
- Added some more diagnostics information.

## v9.3.9 (date: 8.2.2021)

- micromamba cleanup bug fix (got error if micromamba is missing)
- micromamba download bug fix (killed on MacOS)

## v9.3.8 (date: 4.2.2021)

- making started and finished subprocess PIDs visible in --debug level.

## v9.3.7 (date: 4.2.2021)

- micromamba version printout changed, so rcc now parses new format
- micromamba is 0.x, so it does not follow semantic versioning yet, so
  rcc will now "lockstep" versions, with micromamba locked to 0.7.12 now

## v9.3.6 (date: 3.2.2021)

- removing "defaults" channel from robot templates

## v9.3.5 (date: 2.2.2021)

- micromamba upgrade to 0.7.12
- REGRESSION: `rcc task shell` got broken when micromamba was introduced,
  and this version fixes that

## v9.3.4 (date: 1.2.2021)

- fix: removing environments now uses rename first and then delete,
  to get around windows locked files issue
- warning: on windows, if environment is somehow locked by some process,
  this will fail earlier in the process (which is good thing), so be aware
- minor change on cache statistics representation and calculation

## v9.3.3 (date: 1.2.2021)

- adding `--dryrun` option to issue reporting

## v9.3.2 (date: 29.1.2021)

- added environment variables for installation identity, opt-out status as
  `RCC_INSTALLATION_ID` and `RCC_TRACKING_ALLOWED`

## v9.3.1 (date: 29.1.2021)

- fix: when environment is leased, temporary folder is will not be recycled
- cleanup command now cleans also temporary folders based on day limit

## v9.3.0 (date: 28.1.2021)

- support for applications to submit issue reports thru rcc
- print "robot.yaml" to logs, to make it visible for support cases
- diagnostics can now print into a file, and that is used as part
  of issue reporting
- added links to diagnostic checks, for user guidance

## v9.2.0 (date: 25.1.2021)

- experiment: carrier PoC

## v9.1.0 (date: 25.1.2021)

- new command `rcc configure diagnostics` to help identify environment
  related issues
- also requiring new version of micromamba, 0.7.10

## v9.0.2 (date: 21.1.2021)

- fix: prevent direct deletion of leased environment

## v9.0.1 (date: 20.1.2021)

- BREAKING CHANGES
- removal of legacy "package.yaml" support

## v9.0.0 (date: 18.1.2021)

- BREAKING CHANGES
- new cli option `--lease` to request longer lasting environment (1 hour from
  lease request, and next requests refresh the lease)
- new environment variable: `RCC_ENVIRONMENT_HASH` for clients to use
- new command `rcc env unlease` to stop leasing environments
- this breaks contract of pristine environments in cases where one application
  has already requested long living lease, and other wants to use environment
  with exactly same specification (if pristine, it is shared, otherwise it is
  an error)

## v8.0.12 (date: 18.1.2021)
- Templates conda -channel ordering reverted pending conda-forge chagnes.

## v8.0.10 (date: 18.1.2021)

- fix: when there is no pip dependencies, do not try to run pip command

## v8.0.9 (date: 15.1.2021)

- fix: removing one verbosity flag from micromamba invocation

## v8.0.8 (date: 15.1.2021)

- now micromamba 0.7.8 is required
- repodata TTL is reduced to 16 hours, and in case of environment creation
  failure, fall back to 0 seconds TTL (immediate update)
- using new --retry-with-clean-cache option in micromamba

## v8.0.7 (date: 11.1.2021)

- Now rcc manages TEMP and TMP locations for its subprocesses

## v8.0.6 (date: 8.1.2021)

- Updated to robot templates
- conda channels in order for `--strict-channel-priority`
- library versions updated and strict as well (rpaframework v7.1.1)
- Added basic guides for what to do in conda.yaml for end-users.

## v8.0.5 (date: 8.1.2021)

- added robot test to validate required changes, which are common/version.go
  and docs/changelog.md

## v8.0.4 (date: 8.1.2021)

- now requires micromamba 0.7.7 at least, with version check added
- micromamba now brings --repodata-ttl, which rcc currently sets for 7 days
- and touching conda caches is gone because of repodata ttl
- can now also cleanup micromamba binary and with --all
- environment validation checks simplified (no more separate space check)

## v8.0.3 (date: 7.1.2021)

- adding path validation warnings, since they became problem (with pip) now
  that we moved to use micromamba instead of miniconda
- also validation pattern update, with added "~" and "-" as valid characters
- validation is now done on toplevel, so all commands could generate
  those warnings (but currently they don't break anything yet)

## v8.0.2 (date: 5.1.2021)

- fixing failed robot tests for progress indicators (just tests)

## v8.0.1 (date: 5.1.2021)

- added separate pip install phase progress step (just visualization)
- now `rcc env cleanup` has option to remove miniconda3 installation

## v8.0.0 (date: 5.1.2021)

- BREAKING CHANGES
- removed miniconda3 download and installing
- removed all conda commands (check, download, and install)
- environment variables `CONDA_EXE` and `CONDA_PYTHON_EXE` are not available
  anymore (since we don't have conda installation anymore)
- adding micromamba download, installation, and usage functionality
- dropping 32-bit support from windows and linux, this is breaking change,
  so that is why version series goes up to v8

## v7.1.5 (date: 4.1.2021)

- now command `rcc man changelog` shows changelog.md from build moment

## v7.1.4 (date: 4.1.2021)

- bug fix for background metrics not send when application ends too fast
- now all telemetry sending happens in background and synchronized at the end
- added this new changelog.md file

## Older versions

Versions 7.1.3 and older do not have change log entries. This changelog.md
file was started at 4.1.2021.
