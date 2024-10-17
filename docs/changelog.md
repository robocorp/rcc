# rcc change log

## v18.1.7 (date: 17.10.2024)

- adding support for windows development of rcc
- using python/invoke instead of ruby/rake for building rcc
- code formatting python code with ruff
- also using windows runner for tests in github actions

## v18.1.6 (date: 21.8.2024)

- unit tests suite now works properly on MacOS and Windows

## v18.1.5 (date: 7.8.2024)

- developer directory tooling update

## v18.1.4 (date: 5.8.2024)

- bugfix: when there is PS1 in holotree variables, it is now filtered out
- developer directory with toolkit.yaml (hidden robot.yaml)

## v18.1.3 (date: 2.8.2024)

- bugfix: tlsCheck was giving nil TLS information without error
- CONTRIBUTING.md -- things to note when developing rcc

## v18.1.2 (date: 28.6.2024)

- updated default settings.yaml for Sema4.ai products.
- templates location also back in Sema4.ai settings.yaml.
- documentation updates

## v18.1.1 (date: 26.6.2024) WORK IN PROGRESS

- bug fix: too many commands were only visible with `--robocorp` product
  strategy, but they are needed also in `--sema4ai` strategy

## v18.1.0 (date: 26.6.2024) WORK IN PROGRESS

- new command `feedback batch` for applications to send many metrics at once
- disabling rcc internal metrics based on product strategy
- bug fix: journal appending as more atomic operation (just one write)

## v18.0.5 (date: 14.6.2024) WORK IN PROGRESS

- MAJOR breaking change: now command `rcc configuration settings` will require
  `--defaults` flag to show defaults template. Without it, default functionality
  now is to show effective/active settings in YAML format.

## v18.0.4 (date: 12.6.2024) WORK IN PROGRESS

- Additional `--robocorp` product flag added. To match `--sema4ai` flag.
- Now using `%ProgramData%` instead of hard coded `c:\ProgramData\` in code.
- Update on default `settings.yaml` for Sema4.ai products.

## v18.0.3 (date: 7.6.2024) WORK IN PROGRESS

- Windows bugfix: icacls now applied on shared holotree location from product
  strategy (was hardcoded before)

## v18.0.2 (date: 7.6.2024) WORK IN PROGRESS

- default `settings.yaml` is now behind product strategy (each product can
  have their own default settings)

## v18.0.1 (date: 5.6.2024) WORK IN PROGRESS

- added company name as strategy name (dynamic name handling for user messages)
- replaced static "Robocorp" references with strategy name
- renamed some Robocorp functions to more generic Product functions

## v18.0.0 (date: 3.6.2024) WORK IN PROGRESS

- MAJOR breaking change: rcc will now live in two product domains,
  Robocorp and Sema4.ai
- feature: initial support for `--sema4ai` strategy selection
- robot tests to test Sema4.ai support

## v17.29.1 (date: 29.5.2024)

- bugfix: when taking locks, some of those need to be in shared directory,
  while others should not; code was making too much directories shared

## v17.29.0 (date: 27.5.2024)

- bugfix: removing `VIRTUAL_ENV` when rcc is executing subprocesses
- adding warning about that environment variable also in diagnostics

## v17.28.4 (date: 26.4.2024)

- bugfix: when there is "rcc point of view" message, it was not showing
  who was controller, so now controller is visible

## v17.28.3 (date: 26.4.2024)

- bugfix: metrics sending was stating things as error, but they are not
  critical (so that is now mentioned in message)

## v17.28.2 (date: 24.4.2024)

- bugfix: more places are now using package/conda YAML loading

## v17.28.1 (date: 24.4.2024)

- bugfix: when exporting prebuild environments, include layer catalogs also
- bugfix: exporting was not adding all environments correctly to .zip file

## v17.28.0 (date: 22.4.2024)

- adding support for `package.yaml` as replacement for `conda.yaml`

## v17.27.0 (date: 17.4.2024)

- when pip dependencies has `--use-feature=truststore` those environments
  are identified as cacheable
- removed some robot.yaml file diagnostic checks since those are not valid
  anymore

## v17.26.0 (date: 17.4.2024)

- feature: `--no-retry-build` flag for tools to prevent rcc doing retry
  environment build in case of first build fails

## v17.25.0 (date: 17.4.2024)

- bug: when first build failed, original layers were expected to still be there
- fix: now second build always builds all layers (since failure needs that)

## v17.24.0 (date: 15.4.2024)

- micromamba upgrade to v1.5.8

## v17.23.2 (date: 15.4.2024)

- more github action upgrades

## v17.23.1 (date: 15.4.2024)

- github action upgrades

## v17.23.0 (date: 10.4.2024)

- cleanup improvement with option `--caches` to remove conda/pip/uv/hololib
  caches but leave holotree available
- also environment building now cleans up "building" space both before and
  after environment is build

## v17.22.0 (date: 27.3.2024) WORK IN PROGRESS

- compression flag is now globally accessible
- compression flag also switches using siphash as identity hasher
- dirtyness stats also now lists duplicates and linked files

## v17.21.0 (date: 25.3.2024) WORK IN PROGRESS

- experimental feature to disable compression of hololib content
- made cleanup much more strict on error detection
- bug fix: task shell was missing `--space` option; added now

## v17.20.0 (date: 12.3.2024) WORK IN PROGRESS

- uv experiment with limited scope and imperfect implementation

## v17.19.0 (date: 11.3.2024)

- micromamba upgrade to v1.5.7

## v17.18.1 (date: 4.3.2024)

- bugfix: template hash case-sensitivity fix

## v17.18.0 (date: 23.2.2024)

- new `--bundled` flag to support cases where rcc is bundled inside other apps
- first thing behind "bundled" flag is version check (so when flag is given,
  rcc will never check possible newer version existence)
- bugfix: flag handling defaults on peek initialized flags
- typofix: on certificate appending failure message

## v17.17.4 (date: 19.2.2024)

- added `venv.md` for start of documentation for `rcc venv` and depxtraction
  tooling and ideas how to use them
- added new man page command into rcc: `rcc man venv`

## v17.17.3 (date: 17.2.2024)

- bugfix: unwanted logging output can now be hidden using global option
  `--log-hide <text-fragment>` (and can be given multiple times)

## v17.17.2 (date: 14.2.2024)

- depxtraction output update and refactoring code

## v17.17.1 (date: 12.2.2024)

- fixed space .use file to be written only when path is actually known

## v17.17.0 (date: 9.2.2024)

- adding depxtraction as part of `rcc holotree venv` creation

## v17.16.0 (date: 7.2.2024)

- new `NO_PROXY` configuration addition to `settings.yaml` file.
- that `NO_PROXY` will override previous OS level configuration, so be careful
- this closes #57

## v17.15.2 (date: 5.2.2024)

- bugfix: venv activation script search performed after initialization

## v17.15.1 (date: 5.2.2024)

- bugfix: venv creation was missing `--system-site-packages` option, added
- bugfix: in venv creation, picking path from actual environment and then
  using python there to create venv

## v17.15.0 (date: 2.2.2024)

- pull request from https://github.com/SoloJacobs/rcc relating to Windows
  icacls usage. Thank you, Solomon Jacobs and Simon Meggle for bringing
  this up.
- this closes #54

## v17.14.0 (date: 2.2.2024)

- new experimental command `rcc holotree venv` to support python virtual
  environments; this is still "work in progress"

## v17.13.0 (date: 24.1.2024)

- removing internal ECC experiment code (since it never get proper support)
- this should also remove one security vulnerability (Terrapin) hopefully

## v17.12.1 (date: 27.11.2023)

- bugfix: removing duplicates and existing holotree from PATHs before adding
  new items in PATH

## v17.12.0 (date: 23.11.2023)

- reverted changes done in v17.8.0 (git hash 8771d622443efae2aa04c2d8c85b5b5c2e7aa3d6)

## v17.11.0 (date: 23.11.2023)

- adding functionality to mark holotree space as EXTERNALLY-MANAGED (PEP 668)

## v17.10.0 (date: 16.11.2023)

- functionality to tell rcc to not manage anything relating to temporaray
  directories (that is, something else is managing those)
- new environment variable `RCC_NO_TEMP_MANAGEMENT` and new command line flag
  `--no-temp-management` to control above thing
- functionality to tell rcc to not manage anything relating to python .pyc
  files (that is, something else is managing those)
- new environment variable `RCC_NO_PYC_MANAGEMENT` and new command line flag
  `--no-pyc-management` to control above thing
- added diagnostics warnings when above environment variables are set

## v17.9.0 (date: 15.11.2023)

- rcc is now checking if newer released versions are available, and adds
  notification into stderr if not using that version

## v17.8.1 (date: 14.11.2023)

- bug fix: made check of users sharing `ROBOCORP_HOME` case insenstive
- added note on `ROBOCORP_HOME` permissions into documentation
- also `journal.run` has event when multiple users share same home

## v17.8.0 (date: 14.11.2023) REVERTED

- expanded documentation on `rcc robot dependencies` command
- added warning when developer declared dependencies file is missing, and
  when environment configuration drift is shown
- added diagnostics to detect missing dependencies drift file

## v17.7.3 (date: 14.11.2023)

- changed subprocess monitoring from 200ms to 550ms, since on slower machines,
  that 200ms causes too much load (experiment; might need to change later again)

## v17.7.2 (date: 8.11.2023)

- documentation updates on maintenance, and vocabulary, etc.

## v17.7.1 (date: 8.11.2023)

- bugfix: changed used holotree space tracking so, that it is visible to
  everybody on file level

## v17.7.0 (date: 8.11.2023) INCOMPLETE

- added simple tracking of used holotree spaces
- added "Last used" and "Use count" to holotree space listings

## v17.6.1 (date: 6.11.2023) WORK IN PROGRESS

- removed experimental `SSL_CERT_DIR` as environment variable, since it might
  actually be confusing to have there (but diagnostics will remain)
- removed duplicate work on checking catalog integrity which was called
  during holotree restore

## v17.6.0 (date: 2.11.2023) WORK IN PROGRESS

- replaced trollhash with simple relocation detection algorithm and remove
  trollhash from codebase

## v17.5.0 (date: 30.10.2023)

- added `SSL_CERT_DIR` and `NODE_EXTRA_CA_CERTS` as environment variables
  when there is certificate bundle available
- also added diagnostics of those environment variables (plus others)
- minor documentation fixes
- tutorial: example of easy robot run

## v17.4.2 (date: 25.10.2023)

- minor fix: rcc point of view now has version number in it
- new `--anything` flag to allow adding to command line something unique or
  note worthy about that specific line (had no effect what so ever)
- technical: updated some go module dependencies

## v17.4.1 (date: 23.10.2023) WORK IN PROGRESS

- verifying that tlsexport bundle can imported into certificate pool
- using system certificate store as base (if available), and updating
  certificates there by default
- fix on conda.yaml merging on pip options case
- peeking `--debug` and `--trace` flags for preview of verbosity state

## v17.4.0 (date: 23.10.2023) WORK IN PROGRESS

- new subcommand, `rcc configuration tlsexport`, to export TLS certificates
  from given set of secure and unsecure URLs
- now tlsprobe reports fingerprint using sha256 from raw certificate, not
  just plain signature

## v17.3.1 (date: 18.10.2023)

- minor fix: now used micromamba version number is stored in separate asset
  file, to keep things in sync between build scripts and rcc binary

## v17.3.0 (date: 16.10.2023) WORK IN PROGRESS

- embedded micromamba inside rcc executable
- removed micromamba download support since it extract all the way
- removing support for arm64 architectures (linux, mac, windows) since
  embedded micromamba is not available on those architectures

## v17.2.0 (date: 12.10.2023)

- micromamba upgrade to v1.5.1

## v17.1.3 (date: 12.10.2023)

- fix: made used environment configuration visible on progress entry and
  also noted once on first unique contact

## v17.1.2 (date: 11.10.2023)

- bugfix: Windows micromamba activation failures
- bugfix: operating system information was leaking process STDERR
- added operating system information to speed test output

## v17.1.1 (date: 11.10.2023) UNSTABLE

- bugfix: operating system information executed differently in windows
- added hostname and user name to diagnostic information

## v17.1.0 (date: 10.10.2023) UNSTABLE

- operating system infomation on diagnostics and progress items

## v17.0.1 (date: 10.10.2023) UNSTABLE

- early detection of `--warranty-voided` flag to allow init usage
- more functionality skipped when "warranty voided", so that rcc is more
  read-only with that flag

## v17.0.0 (date: 4.10.2023) UNSTABLE

- MAJOR breaking change: removed interactive configuration command, since
  Setup Utility now better covers that functionality
- MAJOR breaking change: holotree is now layered by default and `--layered`
  option is gone
- few documentation updates

## v16.9.0 (date: 3.10.2023) UNSTABLE

- deterioration: added `--warranty-voided` mode to make system less robust but
  faster (do not use this mode, unless you really do know what you are doing)

## v16.8.0 (date: 2.10.2023)

- improvement: quick diagnostics now has settings.yaml age visible as seconds
- added `RCC_REMOTE_ORIGIN` variable to diagnostics output
- deprecated interactive configuration, since Setup Utility should be used

## v16.7.1 (date: 29.9.2023)

- bugfix: added process blacklist to prevent old processes shown as child
  processes in process tree (also recycled PIDs will become "grey listed")
  and this bug was detected in Windows
- improvement: changed command WaitDelay from 15 seconds to 3 seconds

## v16.7.0 (date: 27.9.2023)

- refactored profile commands into one file
- added support for removing configuration profiles
- updated robot tests to test profile removal
- fix: added 3 second timeout to TLS checks

## v16.6.0 (date: 22.9.2023)

- internal probe becomes `rcc configuration tlsprobe` command
- tlsprobe output improvements (address and DNS resolution)
- sending metrics of `rcc.cli.run.failure` when automation exit code is
  something else than zero

## v16.5.0 (date: 21.9.2023)

- new variables set into environments: `RC_DISABLE_SSL`, `WDM_SSL_VERIFY`,
  `NODE_TLS_REJECT_UNAUTHORIZED`, and `RC_TLS_LEGACY_RENEGOTIATION_ALLOWED`
- new settings option `legacy-renegotiation-allowed`
- removed `automation-studio` from `autoupdates:` in settings.yaml file
- settings.yaml version number updated to `2023.09`
- added 5 second timeout to probe connections

## v16.4.1 (date: 21.9.2023) INTERNAL

- improve: refining TLS probe (added cipher suite)

## v16.4.0 (date: 21.9.2023) INTERNAL

- feature: internal TLS probe implementation

## v16.3.1 (date: 20.9.2023)

- bug fix: extracting big template failed
- now some Progress steps have CPUs also visible, in addition to worker count

## v16.3.0 (date: 19.9.2023)

- extended using "rcc point of view" messaging to environment building,
  post-install and pre-run scripts
- holotree variables also now has "rcc point of view" visible
- changed robot tests to match "rcc point of view" changes
- highlighted Progress steps with cyan/green/red color (where available)

## v16.2.2 (date: 13.9.2023)

- bugfix: process tree 1 second delay to prevent "too fast" process snapshots
  on Windows
- refactoring some unused code out of codebase

## v16.2.1 (date: 12.9.2023)

- bugfix: detecting and truncating process tree with too deep child structure

## v16.2.0 (date: 11.9.2023)

- added relocations statistics on catalog listing (Relocate column)

## v16.1.3 (date: 7.9.2023)

- comment explaining why certain unsecure code forms is required when
  TLS diagnostics are done (to explain Github CodeQL security warnings)

## v16.1.2 (date: 6.9.2023) WORK IN PROGRESS

- bug fix: allowing detection of lower levels of TLS versions
- minor improvement: diagnostics TLS firewall/proxy detection
- minor improvement: full certificate chain is now behind `--debug` flag

## v16.1.1 (date: 5.9.2023) WORK IN PROGRESS

- bug fix: added missing proxies to micromamba phase

## v16.1.0 (date: 5.9.2023) WORK IN PROGRESS

- Now advanced network diagnostics also have separate `tls-verify` configuration
  to enable TLS verifications from custom addresses.

## v16.0.1 (date: 5.9.2023) WORK IN PROGRESS

- Added full signature chain "dump" in case where there is some kind of
  certificate failure in TLS verification. Network diagnostics still.

## v16.0.0 (date: 5.9.2023) WORK IN PROGRESS

- Breaking change: there is new TLS verification in place in diagnostic, and
  this can break some old setups because new warnings.

## v15.3.0 (date: 30.8.2023)

- added `journal.run` event log into artifacts directory
- tidying some golang dependencies and removing some unused files

## v15.2.0 (date: 23.8.2023)

- new strategy to manage micromamba, with its own directory based on version
  number: `ROBOCORP_HOME/micromamba/<version>/<executable>`
- updated cleanup to manage micromamba location change
- bugfix: speedtest now does timing also in debug/trace mode (and some other
  minor improvements)

## v15.1.0 (date: 22.8.2023)

- robot diagnostics now has indication of environment cacheability and also
  warnings (category 5010) when something prevents caching
- lack of public cacheability is also visible on environment creation
- documentation updates and improvements
- minor improvements on process tree debugging

## v15.0.0 (date: 21.8.2023) WORK IN PROGRESS

- breaking change: dropped default value `rcc robot initialize --template`
  option (now it must be given)
- breaking change: environment variable `RCC_VERBOSITY` with values "silent",
  "debug", and "trace" now override CLI options
- bugfix, process tree detecting and printing
- added debug/trace logging into process baby sitter
- work in progress: detecting cacheable environment configurations
- micromamba upgrade back to v1.4.9 (next trial)

## v14.15.4 (date: 17.8.2023)

- micromamba downgraded to v1.4.2 due to argument change

## v14.15.3 (date: 14.8.2023)

- added error message on canary failures
- added one diagnostics detail to show if global shared holotree is enabled

## v14.15.2 (date: 10.8.2023)

- bugfix, now giving little bit more stack to process tree

## v14.15.1 (date: 9.8.2023) BROKEN

- bugfix on process tree ending up eating too much stack (stack overflow)

## v14.15.0 (date: 3.8.2023) BROKEN

- micromamba upgrade to v1.4.9

## v14.14.0 (date: 27.6.2023)

- unless silenced, always show "rcc point of view" of success or failure of
  actual main robot run, on point of robot process exit

## v14.13.3 (date: 22.6.2023)

- faster heartbeat for snapshotting subprocesses during robot run (200ms)
- added guiding text on "non-empty artifacts directory case"

## v14.13.2 (date: 21.6.2023)

- micromamba downgrade to v1.4.2, because micromamba bug in Windows

## v14.13.1 (date: 20.6.2023) UNSTABLE

- bugfix: fixing exit code masking by subprocess handling
- predicting rcc exit code made visible
- making robot run exit code more visible
- robot tests now use special settings.yaml to prevent template updates and
  will only use internal templates for testing

## v14.13.0 (date: 15.6.2023) UNSTABLE

- improved listing of still running processes
- set process wait delay to 15 seconds after process has completed but has not
  released it IO pipes yet

## v14.12.0 (date: 13.6.2023) UNSTABLE

- adding listing of still running processes after robot run
- upgrading github actions to use go v1.20.x
- bugfix: panic when using lockpids with nil value

## v14.11.0 (date: 12.6.2023) UNSTABLE

- added `--switch` option to profile import to immediately switch to imported
  profile once it is successfully imported

## v14.10.0 (date: 12.6.2023) UNSTABLE

- saving separate info file for catalogs and holotrees (to speed up some
  commands in future)
- added interrupt signal ignoring around robot run, so that robot can actually
  react and respond to interrupt (and if send twice, then second interrupt
  will actually interrupt rcc)

## v14.9.2 (date: 8.6.2023) UNSTABLE

- more cleaning up of dead code and data structures
- made it visible if artifactsDir already have files before run starts

## v14.9.1 (date: 7.6.2023) UNSTABLE

- micromamba upgrade to v1.4.3

## v14.9.0 (date: 7.6.2023)

- added one user per `ROBOCORP_HOME` warnings
- added also diagnostics to warn about above issue
- full cleanup now also removes `rcccache.yaml` file
- removed "Robots" section from `rcccache.yaml` file

## v14.8.2 (date: 6.6.2023) UNSTABLE

- added missing golden yaml file saving on layers
- added worker count on second progress indicator
- reporting relative time ratios on setup/run balances
- fixed bug in buildstats, where it was using global variables (instead of "it")

## v14.8.1 (date: 5.6.2023) UNSTABLE

- added `RCC_HOLOTREE_SPACE_ROOT` to environment variables provided by rcc
- saving `rcc_plan.log` into intermediate layers as well (and it is now in
  memory presentation while building environment)
- restoring partial environment from layers and skipping already available
  layers (but still only behind `--layered` flag)
- layers add new Progress step to rcc, now total is 15 steps. Test changed
  to match that.

## v14.8.0 (date: 31.5.2023) UNSTABLE

- support for separating layers and calculating their fingerprints
- showing fingerprints on build output and in timeline (still only visualization)
- added controlling flag `--layered` to enable layer handling
- added recording of layers if above flag is given

## v14.7.0 (date: 15.5.2023)

- adding logical layers on holotree installation (visible on timeline)

## v14.6.0 (date: 4.5.2023)

- adding `--quick` flag to diagnostics to filter out slow diagnostics
- for now, "slow diagnostics" are mostly network related checks, some of
  subprocesses still get executed (like micromamba for example)

## v14.5.0 (date: 3.5.2023)

- subprocess exit codes are now visible, when subprocess fails (that is, when
  exit code in non-zero)
- minor update on "custom templates" documentation

## v14.4.1 (date: 19.4.2023)

- MAJOR BREAKING CHANGES:
  - under "spring cleaning" umbrella
  - virtual environment and `pyvenv.cfg` support removed after realization
    that holotree environments are not virtual environments, they are full
    environments and then some, they can also be called soft-containers
  - by trying to be virtual environment also caused bug in Windows, where
    `site-pacakges` and `Scripts` directories could be polluting all other
    environments as well, and that is why there is some cleanup in place now
  - removed old, unused functionality, specially commands `rcc robot fix`,
    `rcc robot libs`, and `rcc robot list` and their relating functionality

## v14.4.0 (date: 19.4.2023) UNSTABLE

- major breaking change: removed `rcc robot libs` command, since it is not
  used in tooling, and if needed, needs better design

## v14.3.0 (date: 19.4.2023) UNSTABLE

- major breaking change: removed `rcc robot fix` command (just command,
  internal functionality is still used by rcc)

## v14.2.0 (date: 18.4.2023) UNSTABLE

- cleanup functionality for "Scripts" and "site-packages" that are in wrong
  place (due virtual environment bug fixed in v14.0.0)

## v14.1.0 (date: 18.4.2023) UNSTABLE

- major breaking change: removed `rcc robot list` command and history handling
  support (this was old Lab requested functionality)

## v14.0.0 (date: 17.4.2023) UNSTABLE

- major breaking change: this will remove some old, now unwanted functionality
- this will be ongoing work for short while, making things unstable for now
- removal of "virtual environment" support (pyvenv.cfg), and `VIRTUAL_ENV`
  variable is no longer available

## v13.12.3 (date: 14.4.2023)

- improvement: more clear messaging on hololib corruption
- fix: full cleanup will first remove catalogs and then hololib

## v13.12.2 (date: 13.4.2023)

- updating documentation around `robot.yaml` and its functionality

## v13.12.1 (date: 12.4.2023)

- fix: added .poetry to list of ignored paths

## v13.12.0 (date: 12.4.2023)

- micromamba upgrade to v1.4.2
- test change: removed test that can fail because of probabilistic feature
  on some metric updates (which cause rcccache.yaml not to be written at all)

## v13.11.0 (date: 5.4.2023)

- tighter permissions restrictions of rcc.yaml and rcccache.yaml using
  os.Chmod, so probably works on Mac and Linux, but Windows is uncertain

## v13.10.1 (date: 20.3.2023)

- diagnostics: minor wording change (removing "toplevel" references)
- documentation: some refinements and additions

## v13.10.0 (date: 16.3.2023)

- documentation: added holotree maintenance documentation
- documentation: added vocabulary/glossary for rcc used terms
- both above also as `rcc docs` subcommands

## v13.9.2 (date: 9.3.2023) DOCUMENTATION

- bugfix: updated toc.py to generate improved table of contents

## v13.9.1 (date: 9.3.2023)

- bugfix: zip verification failed when Windows uses backslashes in paths
- adding diagnostics around `ROBOCORP_HOME` location and robots
- minor documentation updates in relation to `ROBOCORP_HOME` usage

## v13.9.0 (date: 8.3.2023)

- added initial support for verifying that holotree imported zip structure shape
  matches expected hololib catalog patterns (behind `--strict` flag, for now)

## v13.8.0 (date: 7.3.2023)

- new `--export` option to `rcc holotree prebuild` command, to enable direct
  export to given hololib.zip filename of new, successfully build catalogs
- bugfix: catalog was exported before its content, which would make it so, that
  catalog is present before its parts

## v13.7.1 (date: 27.2.2023)

- added missing `RCC_REMOTE_AUTHORIZATION` variable handling to rcc and passing
  that variable to rccremote on pull requests

## v13.7.0 (date: 27.2.2023)

- troubleshooting documentation added as `rcc man troubleshooting` command
- consolidated and streamlined documentation commands into fewer source files
- added robot tests for documentation commands

## v13.6.5 (date: 23.2.2023)

- dependabot raised update on golang.org/x/text module (upgraded)
- security related dependency upgrade

## v13.6.4 (date: 23.2.2023) DOCUMENTATION

- documentation updates on netdiagnostics and troubleshooting

## v13.6.3 (date: 20.2.2023)

- change: changed WorkGroup to not use buffers on incoming messages, since it
  will be more deterministic

## v13.6.2 (date: 16.2.2023) UNSTABLE

- bugfix: changed WaitGroup to WorkGroup (self implemented work synchronization)

## v13.6.1 (date: 15.2.2023)

- experiment: using probability to run some of maintenance functions and
  making rcc little bit faster depending on chance

## v13.6.0 (date: 10.2.2023)

- upgrade: upgrading github actions and also using newer golang and python there

## v13.5.8 (date: 10.2.2023)

- bugfix: holotree delete --space option was always set

## v13.5.7 (date: 10.2.2023)

- bugfix: holotree delete and plan were doing too many calls to find same
  environments (which mean they were really slow)
- some name refactorings to clarify intent of functions

## v13.5.6 (date: 8.2.2023)

- bugfix: create missing folders while creating and writing some files
- improvement: added optional top N biggest files sizes on catalog listing

## v13.5.5 (date: 7.2.2023)

- bugfix: contain output of checking long path support on Windows
- improvement: adding more structure to holotree pull timeline

## v13.5.4 (date: 6.2.2023) UNSTABLE

- bugfix: file syncing on pull commands

## v13.5.3 (date: 6.2.2023) UNSTABLE

- rccremote server zip file managementent improvements

## v13.5.2 (date: 2.2.2023)

- rccremote server timeout adjustments to much longer times (experimental)

## v13.5.1 (date: 2.2.2023)

- fixing progress counter on timeline output
- timeline output clarifications on hololib pull step

## v13.5.0 (date: 2.2.2023)

- support for pulling hololib catalogs as part of normal holotree environment
  creation process (new Progress step).

## v13.4.3 (date: 1.2.2023)

- bugfix: shortcutting to file resource on cloud.ReadFile if actual exiting
  file is given as resource link.

## v13.4.2 (date: 31.1.2023) UNSTABLE

- fixed broken holotree pull command, and made it allow pulling from plain
  http sources

## v13.4.1 (date: 30.1.2023) UNSTABLE

- prebuild now needs shared holotree to be enabled before building
- prebuilds can now be forced for full rebuilds

## v13.4.0 (date: 30.1.2023) UNSTABLE

- peercc is renamed to rccremote, and peercc package renamed to remotree

## v13.3.0 (date: 27.1.2023) UNSTABLE

- feature: command for prebuilding environments (from files or from URLs)
- improvement: rcc version visible in "Toplevel" command list
- added support for "cloud.ReadFile" functionality
- bugfix: wrapped os.TempDir functionality to ensure directory exists

## v13.2.0 (date: 24.1.2023)

- feature: peercc force pulling holotree catalog from other remote peercc
- self pulling should be prevented and so protect self loops
- new settings version, 2023.01 with autoupdates for lab removed and
  setup-utility added

## v13.1.2 (date: 23.1.2023)

- improvement: netdiagnostics with `--trace` flag will now list response
  header information

## v13.1.1 (date: 20.1.2023) UNSTABLE

- fix: netdiagnostics configuration flag change (now it is `--checks filename`)

## v13.1.0 (date: 19.1.2023) UNSTABLE

- feature: more network related configurable diagnostics

## v13.0.1 (date: 17.1.2023)

- bugfix: diagnostics of ignoreFiles was not using paths correctly

## v13.0.0 (date: 17.1.2023)

- major breaking change: various robot unzipping method now flatten directory
  tree so that paths used in robots are shorter and not so easily cause
  problems and confusion

## v12.3.1 (date: 16.1.2023) MAJOR BREAK

- bugfix: unwrap worked wrongly in case of "." dir prefix

## v12.3.0 (date: 13.1.2023) BUGGY MAJOR BREAK

- feature: unwrap command now removes extra middle parts of file paths when
  unzipping robot.zip files

## v12.2.0 (date: 11.1.2023)

- micromamba upgrade to v1.1.0

## v12.1.2 (date: 11.1.2023)

- bugfix: parallel long path checks failed because not unique path was used,
  added pid as part of that long path (just Windows), this closes #45

## v12.1.1 (date: 4.1.2023)

- bugfix: adding more info when zip extraction fails

## v12.1.0 (date: 3.1.2023)

- feature: on assistant runs, if CR does not give artifact URL for uploading
  artifacts, then it is now considered as disabled functionality (not error)
  and no artifacts are pushed into cloud

## v12.0.1 (date: 3.1.2023)

- added diagnostics on loading ignoreFiles entry, which does not contain
  any patterns in it
- updated documentation about `ignoreFiles:` in recipes, with hopefully
  better explanation of how it should be used

## v12.0.0 (date: 29.12.2022) UNSTABLE

- adding "grace period" in "token time" calculations, and this is breaking
  change, because token time calculation changes, and management of grace
  period is user/app responsibility (but there is default value) and tokens
  also will now have minimum period
- bugfix: when broken catalog was loaded, catalog listing failed

## v11.36.5 (date: 28.12.2022)

- fix: added more explanation to network diagnostics reporting, explaining
  what actual successful check option did

## v11.36.4 (date: 22.12.2022)

- bugfix: added missing lock protections around importing and pulling holotrees

## v11.36.3 (date: 21.12.2022)

- improvement: added more color and changed wording on lock wait messages

## v11.36.2 (date: 21.12.2022)

- improvement: when there is longer lock wait, possible lock holders are listed
  on console output and in timeline

## v11.36.1 (date: 20.12.2022)

- bugfix: diagnostics fail on new machine to touch lock files when directory
  does not exist, this closes #43
- bugfix: stale lock pid files are shown too often, this closes #42
- diagnostics will now show hopefully more human friendly message when active
  locks are detected
- added more runtime.Gosched calls to enable background go routines to have
  chance to finish before application closes

## v11.36.0 (date: 15.12.2022)

- added category field into diagnostics JSON output, to support applications
  to report better diagnostics messages

## v11.35.7 (date: 15.12.2022)

- this v11.35.x series adds new "peercc" executable and new holotree pull
  subcommand to rcc; these are work in progress, and not ready for production
  work yet; do not use, unless you know what you are doing
- added automatic import of delta environment update data
- tech: moved TryRemove, TryRemoveAll, and TryRename to pathlib
- tech: some zipper log verbosity was moved from Debug to Trace level

## v11.35.6 (date: 14.12.2022) UNSTABLE

- bug fix: ignoring dotfiles and directories in "pids" directory
- added new `rcc holotree pull` command to do delta environment update request
  to peercc (still incomplete, does not do automatic import of content)
- on delta export zip, catalog will now come as last part of that zip from wire
- added set membership map functionality (to make faster membership checks on
  bigger member sets)
- more failed parts of PoC removed (export specification and support functions)

## v11.35.5 (date: 9.12.2022) UNSTABLE

- fixed bug where last line of request was missing
- trying to fix CodeQL security warning (user input was already filtered based
  on known set of values, but analyzer did not understand that)

## v11.35.4 (date: 8.12.2022) UNSTABLE

- removing failed parts of PoC
- added handler for streaming of requested catalog and missing parts
- made robot tests to automatically disconnect from shared holotree

## v11.35.3 (date: 7.12.2022) UNSTABLE

- replaced deprecated "ioutil" with suitable functions elsewhere, thank you
  for Juneezee (Eng Zer Jun) for pointing these out in PR#40
- added ComSpec, LANG and SHELL from environment into diagnostics output

## v11.35.2 (date: 7.12.2022) UNSTABLE

- next try to fix ruby support in GHA

## v11.35.1 (date: 7.12.2022) UNSTABLE

- github actions updated to use ruby 2.7 (github stopped supporting used 2.5)

## v11.35.0 (date: 7.12.2022) UNSTABLE

- starting new PoC on topic of "peer rcc"
- export specification simplification: now supports exactly one "wants" value
  and it is not list anymore, but just plain and simple string
- added new "set" operations to support PoC functionality (generics)
- one part of PoC failed, but code is still there

## v11.34.0 (date: 29.11.2022)

- compiling rcc for arm64 architectures (linux, mac, windows)

## v11.33.2 (date: 24.11.2022)

- configuration diagnostics now measure and report time it takes to resolve
  set of hostnames found from settings files

## v11.33.1 (date: 23.11.2022) UNSTABLE

- some additional timeline markers on assistant runs

## v11.33.0 (date: 18.11.2022) UNSTABLE

- feature: holotree delta export (for missing things only)
- changes normal holotree export command to support ".hld" files

## v11.32.6 (date: 15.11.2022)

- bugfix: from now on, lock pid files will only give diagnostic "warning" when
  they are less than 12 hours old, after that they will be labeled as "stale"
  and will still be visible in diagnostics, but on "ok" level

## v11.32.5 (date: 15.11.2022)

- cleanup: removing dead code that was not used anymore

## v11.32.4 (date: 15.11.2022)

- cleanup: removing old run minutes and stat lines (holotree stats cover those)

## v11.32.3 (date: 14.11.2022)

- holotree statistics are now part of human readable diagnostics when there
  is 5 or more entries in statistics (but not available in JSON output)
- added cumulative statistics section into output
- bugfix: calculation mistakes in case of missing steps
- bugfix: detecting successful build

## v11.32.2 (date: 11.11.2022)

- added week limitation option for holotree statistics command
- added filter flags for assistants, robots, prepares, and variables for
  holotree statistics command

## v11.32.1 (date: 11.11.2022) UNSTABLE

- feature: command to show local holotree environment build statistics

## v11.32.0 (date: 10.11.2022) UNSTABLE

- feature: local recording of holotree environment build statistics events
- moved journals to `ROBOCORP_HOME/journals` directory (and build stats will
  be part of those journals)
- added pre run scripts to timeline

## v11.31.2 (date: 8.11.2022)

- bugfix: removing path separators from user name on lock pid files

## v11.31.1 (date: 8.11.2022)

- bugfix: changed lock pid filename not to contain extra dots
- added more info on pending lock files diagnostics check
- more debug information on Windows locking behaviour

## v11.31.0 (date: 7.11.2022)

- micromamba upgrade to v1.0.0

## v11.30.1 (date: 7.11.2022)

- bugfix: added more checks around shared holotree enabling and using
- bugfix: make all lockfiles readable and writable by all
- added "diagnostics" command to toplevel commands

## v11.30.0 (date: 2.11.2022)

- added warning when vulnerable openssl is installed in environment

## v11.29.1 (date: 26.10.2022) UNSTABLE

- robot tests for unmanaged holotree spaces (revealed bugs)
- bugfix: correct checking of unmanaged space conflicts (on creation)

## v11.29.0 (date: 25.10.2022) BROKEN

- started adding support for unmanaged holotree spaces, to enable IT managed
  holotree spaces (rcc will create them once, but integrity check are not
  done when unmanaged spaces are used)
- bugfix: removing also .lck files when removing space

## v11.28.3 (date: 19.10.2022)

- added configuration diagnostic reporting on locking pids information

## v11.28.2 (date: 19.10.2022)

- made lock wait messages little more descriptive and added more of them
- added "pids" folder to keep track who is holding locks (just information)

## v11.28.1 (date: 12.10.2022)

- bugfix: direct initializing robot did not update templates

## v11.28.0 (date: 5.10.2022)

- micromamba upgrade to v0.27.0
- refactored version micromamba version numbering into one place
- added used pip and micromamba versions in progress messages
- BUGFIX: now explicitely using environment python to run pip commands
  (using `python -m pip install ...` form instead old `pip install` form)

## v11.27.3 (date: 29.9.2022)

- fix: adding more "plan analyzer" identifiers to its output
- fix: adding detection to "failed to build" messages
- fix: added json output to new robot creation to cloud

## v11.27.2 (date: 27.9.2022)

- improving plan analyzer with more rules to show messages

## v11.27.1 (date: 26.9.2022)

- fixing CodeQL security warning about allocation overflow

## v11.27.0 (date: 23.9.2022)

- support for analyzing installation plans and their challenges and show it
  online, or afterwards
- analysis is visible in `rcc holotree plan` command and also in `pip`
  phase in environment creation

## v11.26.6 (date: 19.9.2022)

- try to upgraded cobra and viper dependencies, to get remove security warnings
  given by AWS container scanner tooling
- upgrade to use github.com/spf13/cobra v1.5.0
- upgrade to use github.com/spf13/viper v1.13.0
- upgrade to use gopkg.in/square/go-jose.v2 v2.6.0

## v11.26.5 (date: 16.9.2022)

- added architecture/platform metric with same interval as timezone metrics
- `docs/history.md` updated with v11 information so far
- `docs/troubleshooting.md` updated with additional points

## v11.26.4 (date: 14.9.2022)

- new `docs/troubleshooting.md` document added
- new `docs/history.md` document added
- updated `scripts/toc.py` with new documents and minor improvement

## v11.26.3 (date: 12.9.2022)

- bugfix: moved "ht.lck" inside holotree location, and renamed it to be
  `global.lck` file.
- added environment variable `SSL_CERT_FILE` to point into certificate bundle
  if one is provided by profile
- documentation updates

## v11.26.2 (date: 8.9.2022)

- converted assets to embedded resources (golang builtin embed module)
- go-bindata is not used anymore (replaced by "embed")

## v11.26.1 (date: 7.9.2022)

- minor documentation improvement, highlighting configuration settings help,
  that plain commands are showing vanilla rcc setting by default.

## v11.26.0 (date: 7.9.2022)

- experiment: pyvenv.cfg file written into created holotree before lifting
- update: cloud-linking in setting.yaml now points to new default location:
  https://cloud.robocorp.com/link/
- bugfix: settings.yaml version updated to 2022.09 (because options section)

## v11.25.1 (date: 6.9.2022)

- fix: symbolic link restoration, when target is actually non-symlink

## v11.25.0 (date: 5.9.2022)

- flag to show identity.yaml (conda.yaml) in holotree catalogs listing
  and functionality then just show it as part of output, both human readable
  and machine readable (JSON)

## v11.24.0 (date: 2.9.2022)

- refactoring some utility functions to more common locations
- adding rcc and micromamba binary locations to diagnostics
- added `RCC_EXE` environment variable available for robots
- added `RCC_NO_BUILD` environment variable support (in addition to
  previous settings options and CLI flag; see v11.19.0)
- some documentation updates
- added support for toplevel `--version` option

## v11.23.0 (date: 2.9.2022)

- added unused option to holotree catalog removal command
- added maintenance related robot test suite
- minor documentation updates

## v11.22.1 (date: 1.9.2022) BROKEN

- fix: using wrong file for age calculation on holotree catalogs
- fix: holotree check failed to recover on corrupted files; now failure
  leads to removal of broken file
- fix: empty hololib directories are now removed on holotree check

## v11.22.0 (date: 31.8.2022)

- new command `rcc holotree remove` added, and this will remove catalogs
  from holotree library (hololib)
- added repeat count to holotree check command (used also from remove command)

## v11.21.0 (date: 30.8.2022)

- added support to tracking when catalog blueprints are used
- if there is no tracking info on existing catalog, first reporting will
  reset it to zero (and report it as -1)
- added catalog age in days, and days since last used to catalog listing
- fixed bug on shared hololib location on catalog listing

## v11.20.0 (date: 26.8.2022)

- feature: allow holotree exporting using robot.yaml file.

## v11.19.1 (date: 25.8.2022)

- bug: empty entry on ignoreFiles caused unclear error
- fix: now empty entries are diagnosed and noted
- fix: also non-existing ignore files are diagnosed

## v11.19.0 (date: 24.8.2022)

- new global flag `--no-build` which prevents building environments, and
  only allows using previously cached, prebuild or imported holotrees
- there is also "no-build" option in "settings.yaml" options section
- added "no-build" information to diagnostics output

## v11.18.0 (date: 23.8.2022)

- new cleanup option `--downloads` to remove downloads caches (conda, pip,
  and templates)
- change: now conda pkgs is cleaned up also in quick cleanup (which now
  includes all "downloads" cleanups)
- robot cache is now part of full cleanup
- run commands now cleanup their temp folders immediately

## v11.17.2 (date: 19.8.2022)

- bugfix: adding missing symbolic link handling of files and directories
- hololib catalogs now have rcc version information included
- added timeout to account deletion, to speed up unit tests

## v11.17.1 (date: 18.8.2022) UNSTABLE

- fix continued: adding missing symbolic link handling of files and directories

## v11.17.0 (date: 17.8.2022) UNSTABLE

- fix started: adding missing symbolic link handling of files and directories
- this will be UNSTABLE, work in progress, for now

## v11.16.0 (date: 16.8.2022)

- micromamba upgrade to v0.25.1
- template upgrade of python to 3.9.13
- template upgrade of pip to 22.1.2
- template upgrade of rpaframework to 15.6.0
- upgraded tests to match above version changes and their effects

## v11.15.4 (date: 13.7.2022)

- go-bindata was accidentally removed, adding it back

## v11.15.3 (date: 13.7.2022) BROKEN

- refactoring module dependencies to help reusing parts of rcc in other apps

## v11.15.2 (date: 11.7.2022)

- fixed table of contents links to match Github generated ones
- also tried to make toc.py more OS neutral (was failing on Windows)

## v11.15.1 (date: 8.7.2022)

- added "old school CI" recipe into documentation

## v11.15.0 (date: 7.7.2022)

- micromamba upgrade to v0.24.0

## v11.14.5 (date: 22.6.2022)

- added `--once` flag to holotree shared enabling, in cases where costly
  sharing is required only once

## v11.14.4 (date: 15.6.2022)

- holotree share enabling now uses "icals" in Windows to set default properties
- added marker file "shared.yes" when shared has been executed

## v11.14.3 (date: 9.6.2022)

- upgraded rcc to be build using go v1.18.x

## v11.14.2 (date: 8.6.2022)

- retry on fixing codeql-analysis problem

## v11.14.1 (date: 8.6.2022)

- fixing codeql-analysis settings and problems
- no codeql analysis for ruby or python in this repo

## v11.14.0 (date: 7.6.2022)

- experimenting on setting `VIRTUAL_ENV` environment variable to point into
  environment rcc created environment
- made OS and architecture visible in rcc "Progress 2" marker

## v11.13.0 (date: 19.5.2022)

- new shared holotree should now be effective
- some instructions on recipes for enabling shared holotree
- micromamba upgrade to v0.23.2

## v11.12.9 (date: 17.5.2022) UNSTABLE

- bugfix: effective user id did not work on windows, removing it for all OSs
- diagnostics now has true/false flag to indicated shared/private holotrees

## v11.12.8 (date: 16.5.2022) UNSTABLE

- bugfix: making shared directories shared only when they really are
- new command `rcc holotree shared --enable` to enable shared holotrees
  in specific machine
- command `rcc holotree init` is now for normal users after shared command

## v11.12.7 (date: 12.5.2022) UNSTABLE

- micromamba upgrade to v0.23.1
- added checks for hololib shared locations mode requirements

## v11.12.6 (date: 10.5.2022) UNSTABLE

- bugfix: added additional directory for hololib, since it helps mounting
  on servers
- one recipe addition, for idea generation ...

## v11.12.5 (date: 9.5.2022) UNSTABLE

- micromamba upgrade to v0.23.0

## v11.12.4 (date: 9.5.2022)

- bugfix: rcc task script could not find any task (reason: internal quoting)
- this closes #32

## v11.12.3 (date: 5.5.2022) UNSTABLE

- Reverted the change in v11.12.2 based on further testing.

## v11.12.2 (date: 5.5.2022) UNSTABLE

- legacyfix: Adding `x-` prefix to custom header, due to some enterprise network proxies stripping headers.

## v11.12.1 (date: 29.4.2022)

- bugfix: duplicate devTask note when error needs to be shown

## v11.12.0 (date: 28.4.2022)

- `rcc holotree import` now supports URL imports

## v11.11.2 (date: 26.4.2022)

- more documentation update for devTasks

## v11.11.1 (date: 26.4.2022)

- bugfix: added v12 indicator in new form of holotree catalogs (to separate
  from old ones) to allow old and new versions to co-exist
- documentation update for devTasks and ToC

## v11.11.0 (date: 25.4.2022)

- in addition to normal tasks, now robot.yaml can also contain devTasks
- it is activated with flag `--dev` and only available in task run command

## v11.10.7 (date: 22.4.2022)

- bugfix/retry: lock files are now marked as shared files (actually this
  will not work on Windows on multi-user setup)
- changed robot test setup cleanup

## v11.10.6 (date: 21.4.2022)

- bugfix: lock files are now marked as shared files
- Rakefile: using `go install` for bindata from now on

## v11.10.5 (date: 20.4.2022)

- settings certificates section now has full path to CA bundle if available

## v11.10.4 (date: 20.4.2022)

- different preRunScripts for different operating systems
- acceptable differentiation patterns are: amd64/arm64/darwin/windows/linux

## v11.10.3 (date: 19.4.2022)

- fixed panic when settings.yaml is broken, now it will be blunt fatal failure
- fixed search path problem with preRunScripts (now robot.yaml PATH is used)
- removed direct mentions of Robocorp App (old name)

## v11.10.2 (date: 13.4.2022) UNSTABLE

- made presence of hololib.zip more visible on environment creation
- test changed so that new holotree relocation can be tested

## v11.10.1 (date: 12.4.2022) UNSTABLE

- adding support for user identity in relocation
- added holotree locations into diagnostics
- usage of hololib.zip now has debug entry in log files

## v11.10.0 (date: 12.4.2022) UNSTABLE

- work in progess: holotree relocation revisiting started
- now holotree library and catalog can be on shared location

## v11.9.16 (date: 7.4.2022)

- took pull request with documentation fixes

## v11.9.15 (date: 6.4.2022)

- improved table of contents with more clearer numbering

## v11.9.14 (date: 6.4.2022)

- improved table of contents with numbering

## v11.9.13 (date: 6.4.2022)

- added new script for generating table of contents for docs
- generated first table of contents as `docs/README.md` file

## v11.9.12 (date: 6.4.2022)

- added new `rcc man profiles` documentation command
- more documentation updates
- `Robocorp Cloud` to `Robocorp Control Room` related documentation changes

## v11.9.11 (date: 5.4.2022)

- documentation updates

## v11.9.10 (date: 4.4.2022)

- added current profile in JSON response from configuration switch command
- fixing bugs and typos in code and texts

## v11.9.9 (date: 31.3.2022)

- updated interactive create with `--task` option alternative
- updated run error message with `--task` option instructions
- this closes #28
- updated `recipes.md` with python conversion instructions

## v11.9.8 (date: 29.3.2022)

- updated profile documentation
- added integrity check on hololib to space extraction
- more robot tests added
- fixed ssl-no-revoke bug (found thru new robot tests)

## v11.9.7 (date: 28.3.2022)

- profiles should now be good enough to start testing them
- interactive configuration now has instructions for next steps (kind of
  scripted but not automated; copy-paste instructions)
- added placeholder `docs/profile_configuration.md` for future documentation
- settings.yaml now has Automation Studio autoupdate URL
- added `robot_tests/profiles.robot` to test new functionality

## v11.9.6 (date: 25.3.2022) UNSTABLE

- adding more setting options and environment variables
- added support for CA-bundles in pem format

## v11.9.5 (date: 23.3.2022) UNSTABLE

- refactoring variables exporting into one place
- adding `PIP_CONFIG_FILE`, `HTTP_PROXY`, and `HTTPS_PROXY` variables into
  conda environment if they are configured

## v11.9.4 (date: 22.3.2022) UNSTABLE

- profile exporting now works

## v11.9.3 (date: 22.3.2022) UNSTABLE

- started to add real support for profile switching/importing
- some documentation updates

## v11.9.2 (date: 18.3.2022) UNSTABLE

- settings are now layered, so that partial custom settings.yaml also works
- settings now have flat API interface, that is used instead of direct access
- settings.yaml version upgrade with new fields (still incomplete)
- endpoints in settings are now a map and not separate structure anymore
- partial "demo" work on interactive configuration (work in progress)

## v11.9.1 (date: 10.3.2022) UNSTABLE

- added condarc and piprc to be asked from user as configuration options
- refactoring some wizard code to better support new functionality

## v11.9.0 (date: 9.3.2022) UNSTABLE

- new work started around network configuration topics and this will be
  WIP (work in progress) for a while, and so it is marked as unstable
- added new command placeholders (no-op): `interactive configuration`,
  `configuration export`, `configuration import`, and `configuration switch`

## v11.8.0 (date: 8.3.2022)

- added initial alpha support for pre-run scripts from robot.yaml and executed
  right before actual task is run

## v11.7.1 (date: 8.3.2022)

- when timeline option is given, and operation fails, timeline was not shown
  and this change now makes timeline happen before exit is done
- speed test now allows using debug flag to actually see what is going on

## v11.7.0 (date: 8.3.2022)

- micromamba update to version 0.22.0

## v11.6.6 (date: 7.3.2022)

- JSON/YAML diagnostics is now ignoring anything that contains ".vscode"

## v11.6.5 (date: 2.3.2022)

- Still continuing GH#27 fixing issue where rcc finds executables outside of
  holotree environment.

## v11.6.4 (date: 23.2.2022)

- GH#27 fixing issue where rcc finds executables outside of holotree
  environment.
- this closes #27

## v11.6.3 (date: 10.1.2022)

- more patterns added ("pypoetry" and "virtualenv") to be removed from PATH,
  since they also can break isolation of our environments

## v11.6.2 (date: 7.1.2022)

- added "pyenv" and "venv" to patterns removed from PATH, since they can
  break isolation of our environments

## v11.6.1 (date: 7.1.2022)

- fixing micromamba version number parsing

## v11.6.0 (date: 7.12.2021) broken

- micromamba update to version 0.19.0
- now `artifactsDir` is explicitely created before robot execution

## v11.5.5 (date: 2.11.2021)

- bugfix: robot task format ignored artifacts directory, but now it uses it

## v11.5.4 (date: 29.10.2021)

- bugfix: path handling in robot wrap commands (now cross-platform)

## v11.5.3 (date: 28.10.2021) broken

- bugfix: path handling in robot wrap commands

## v11.5.2 (date: 27.10.2021)

- added `--json` option and output to catalogs listing
- bug fix: added missing file detection to holotree check

## v11.5.1 (date: 26.10.2021)

- adding holotree catalogs command to list available catalogs with more detail
- extending holotree list command to show all spaces reachable from hololib
  catalogs including imported holotree spaces
- holotree delete should now also remove space elsewhere (based on imported
  catalogs and their holotree locations)

## v11.5.0 (date: 20.10.2021)

- adding initial support for importing hololib.zips into local hololib catalog

## v11.4.3 (date: 20.10.2021)

- fixing bug where gzipped files in virtual holotree get accidentally
  expanded when doing `--liveonly` environments
- added global `--workers` option to allow control of background worker count

## v11.4.2 (date: 19.10.2021)

- one more improvement on abstract score reporting (time is also scored)

## v11.4.1 (date: 18.10.2021)

- minor textual improvements on abstract score reporting

## v11.4.0 (date: 18.10.2021)

- new command `rcc configuration speedtest` which gives abstract score to both
  network and filesystem speed
- some refactoring to enable above functionality

## v11.3.6 (date: 13.10.2021)

- bugfix: added retries to holotree file removal

## v11.3.5 (date: 12.10.2021)

- bugfix: added retries and better error message on holotree rename pattern

## v11.3.4 (date: 12.10.2021)

- new toplevel flag to turn on `--strict` environment handling, and for now
  this make rcc to run `pip check` after environment install completes
- added timeout to metrics sending

## v11.3.3 (date: 8.10.2021)

- micromamba update to version 0.16.0
- minor change on os.Stat usage in holotree functions
- changed minimum required worker count to 2 (was 4 previously)

## v11.3.2 (date: 7.10.2021)

- templates are removed when quick cleanup is requested
- bugfix: now debug and trace flags are also considered same as
  `VERBOSE_ENVIRONMENT_BUILDING` environment variable
- bugfix: added some jupyter paths as skipped ingored ones in diagnostics
- added canary checks into diagnostics for pypi and conda repos

## v11.3.1 (date: 5.10.2021)

- using templates from templates.zip in addition to internal templates
- command `rcc holotree bootstrap` update to use templates.zip
- command `rcc interactive create` now uses template descriptions
- command `rcc robot init` now has `--json` flag to produce template list
  as JSON
- settings.yaml updated to version 2021.10

## v11.3.0 (date: 4.10.2021)

- update robot templates from cloud (not used yet, coming up in next versions)

## v11.2.0 (date: 29.9.2021)

- updated content to [recipes](/docs/recipes.md) about Holotree controls
- two new documentation commands, features and usecases, and corresponding
  markdown documents in docs folder
- added env.json capability also into pure `conda.yaml` case in
  `rcc holotree variables` command (bugfix)

## v11.1.6 (date: 27.9.2021)

### What to consider when upgrading from series 10 to series 11 of rcc?

Major version break between rcc 10 and 11 was about removing the old base/live
way of managing environments (`rcc environment ...` commands). That had some
visible changes in rcc commands used. Here is a summary of those changes.

- Compared to base/live based management of environments, holotree needs
  a different mindset to work with. With the new holotree, users decide which
  processes share the same working space and which receive their own space.
  So, high level management of logical spaces has shifted from rcc to user
  (or tools), where in base/live users did not have the option to do so.
  Low level management is still rcc responsibility and based on "conda.yaml"
  content.
- All `rcc environment` commands were removed or renamed, since this was
  an old way of doing things.
- Old `rcc env hash` was renamed to `rcc holotree hash` and changed to show
  holotree blueprint hash.
- Old `rcc env plan` was renamed to `rcc holotree plan` and changed to show
  plan from given holotree space.
- Old `rcc env cleanup` was renamed to `rcc configuration cleanup` and
  changed to work in a way that only holotree things are valid from now on.
  This means that if you are using `rcc conf cleanup`, check help for changed
  flags also.
- In general, the old `--stage` flag is gone, since it was base/live specific.
- Holotree related commands, including various run commands, now have default
  values for the `--space` flag. So if no `--space` flag is given, that
  defaults to `user` value, and the same space will be updated based
  on requested environment specification.
- Output of some commands have changed, for example there are now more
  "Progress" steps in rcc output.

## v11.1.5 (date: 24.9.2021)

- bugfix: performance profiling revealed bottleneck in windows, where calling
  stat is expensive, so here is try to limit using it uneccessarily

## v11.1.4 (date: 23.9.2021)

- bugfix: adding concurrencty to catalog check
- performance profiling revealed bottleneck, where ensuring directory exist
  was called too often, so now base directories are ensured only once per
  rcc invocation
- adding more structure to timeline printout by indentation of blocks

## v11.1.3 (date: 21.9.2021)

- bugfix: changing performance thru auto-scaling workers based on number
  of CPUs (minus one, but at least 4 workers)

## v11.1.2 (date: 20.9.2021)

- bugfix: removing duplicate file copy on holotree recording
- removed "new live" phrase from debug printouts
- made robot tests to check holotree integrity in some selected points

## v11.1.1 (date: 17.9.2021)

- bugfix: using rename in hololib file copy to make it more transactional
- progress indicator now has elapsed time since previous progress entry
- experimental upgrade to use go 1.17 on Github Actions

## v11.1.0 (date: 16.9.2021)

- BREAKING CHANGES, but now this may be considered stable(ish)
- micromamba update to version 0.15.3
- added more robot tests and improved `rcc holotree plan` command

## v11.0.8 (date: 15.9.2021) UNSTABLE

- BREAKING CHANGES (ongoing work, see v11.0.0 for more details)
- showing correct `rcc_plan.log` and `identity.yaml` files on log
- reorganizing some common code away from conda module
- rpaframework upgrade to version 11.1.3 in templates

## v11.0.7 (date: 14.9.2021) UNSTABLE

- BREAKING CHANGES (ongoing work, see v11.0.0 for more details)
- changed progress indication to match holotree flow
- made log and telemetry waiting visible in timeline

## v11.0.6 (date: 13.9.2021) UNSTABLE

- BREAKING CHANGES (ongoing work, see v11.0.0 for more details)
- removing options from cleanup commands, since those are base/live specific
  and not needed anymore (orphans, miniconda)
- removed dead code resulted from above

## v11.0.5 (date: 10.9.2021) UNSTABLE

- BREAKING CHANGES (ongoing work, see v11.0.0 for more details)
- removing conda environment build related code
- internal clone command was removed
- side note: there is trail of FIXME comments in code for future work

## v11.0.4 (date: 9.9.2021) UNSTABLE

- BREAKING CHANGES (ongoing work, see v11.0.0 for more details)
- replaced `rcc env plan` with new `rcc holotree plan`, which now shows
  installation plans from holotree spaces
- now all env commands are removed, so also toplevel "env" command is gone
- added naive helper script, deadcode.py, to find dead code
- cleaned up some dead code branches

## v11.0.3 (date: 8.9.2021) UNSTABLE

- BREAKING CHANGES (ongoing work, see v11.0.0 for more details)
- removed commands "new", "delete", "list", and "variables" from `rcc env`
  command set
- replaced `rcc env hash` with new `rcc holotree hash`, which now calculates
  blueprint fingerprint hash similar way that env hash but differently
  because holotree uses siphash algorithm

## v11.0.2 (date: 8.9.2021) UNSTABLE

- BREAKING CHANGES (ongoing work, see v11.0.0 for more details)
- technical work: cherry-picking changes from v10.10.0 into series 11

## v11.0.1 (date: 7.9.2021) UNSTABLE

- BREAKING CHANGES (ongoing work, see v11.0.0 for more details)
- fixing robot tests

## v11.0.0 (date: 6.9.2021) UNSTABLE

- BREAKING CHANGES (ongoing work, small steps, considered unstable) and goal
  is to remove old base/live environment handling and make holotree default
  and only way to manage environments
- setting "user" as default space for all commands that need environments

## v10.10.0 (date: 7.9.2021)

- this is series 10 maitenance branch
- rcc config cleanup improvement, so that not partial cleanup is done on
  holotree structure (on Windows, respecting locked environments)

## v10.9.4 (date: 31.8.2021)

- invalidating hololib catalogs with broken files in hololib

## v10.9.3 (date: 31.8.2021)

- added diagnostic warnings on `PLAYWRIGHT_BROWSERS_PATH`, `NODE_OPTIONS`,
  and `NODE_PATH` environment variables when they are set

## v10.9.2 (date: 30.8.2021)

- bugfix: long running assistant run now updates access tokens correctly

## v10.9.1 (date: 27.8.2021)

- made problems in assistant heartbeats visible
- changed assistant heartbeat from 60s to 37s to prevent collision with
  DNS TTL value

## v10.9.0 (date: 25.8.2021)

- added --quick option to `rcc config cleanup` command to provide
  partial cleanup, but leave hololib and pkgs cache intact

## v10.8.1 (date: 24.8.2021)

- holotree check command now removes orphan hololib files
- environment creation metrics added on failure cases
- pip and micromamba exit codes now also in hex form
- minor error message fixes for Windows (colors)

## v10.8.0 (date: 19.8.2021)

- added holotree check command to verify holotree library integrity
- added "env cleanup" also as "config cleanup"
- minor go-routine schedule yield added (experiment)

## v10.7.1 (date: 18.8.2021)

- bugfix: trying to remove preformance hit on windows directory cleanup

## v10.7.0 (date: 16.8.2021)

- when environment creation is serialized, after short delay, rcc reports
  that it is waiting to be able to contiue
- added \_\_MACOSX as ignored files/directories

## v10.6.0 (date: 16.8.2021)

- added possibility to also delete holotree space by providing controller
  and space flags (for easier scripting)

## v10.5.2 (date: 12.8.2021)

- added once a day metric about timezone where rcc is executed

## v10.5.1 (date: 11.8.2021)

- improvements for detecting OS/architecture for multiple environment
  configurations

## v10.5.0 (date: 10.8.2021)

- supporting multiple environment configurations to enable operating system
  and architecture specific freeze files (within one robot project)

## v10.4.5 (date: 10.8.2021)

- bugfix: removing one more filesystem sync from holotree (Mac slowdown fix).

## v10.4.4 (date: 9.8.2021) broken

- bugfix: raising initial scaling factor to 16, so that there should always
  be workers waiting

## v10.4.3 (date: 9.8.2021) broken

- bugfix: trying to fix Mac related slowing by removing file syncs on
  holotree copy processes

## v10.4.2 (date: 5.8.2021) broken

- bugfix: scaling down holotree concurrency, since at least Mac file limits
  are hit by current concurrency limit

## v10.4.1 (date: 5.8.2021)

- taking micromamba 0.15.2 into use

## v10.4.0 (date: 5.8.2021)

- bug fix: `rcc_activate.sh` were failing, when path to rcc has spaces in it

## v10.3.3 (date: 29.6.2021)

- updated tips, tricks, and recipes

## v10.3.2 (date: 29.6.2021)

- fix for missing artifact directory on runs

## v10.3.1 (date: 29.6.2021) broken

- cleaning up `rcc robot dependencies` and related code now that freeze is
  actually implemented
- changed `--copy` to `--export` since it better describes the action
- removed `--bind` because copying freeze file from run is better way
- removed "ideal" conda.yaml printout, since runs now create artifact
  on every run in new envrionments
- removed those robot diagnostics that are misguiding now when dependencies
  are frozen
- updated rpaframework to version 10.3.0 in templates
- updated robot tests for rcc

## v10.3.0 (date: 28.6.2021)

- creating environment freeze YAML file into output directory on every run

## v10.2.4 (date: 24.6.2021)

- added `--bind` option to copy exact dependencies from `dependencies.yaml`
  into `conda.yaml`, so that `conda.yaml` represents fixed dependencies

## v10.2.3 (date: 24.6.2021)

- added `dependencies.yaml` into robot diagnostics
- show ideal `conda.yaml` that matches `dependencies.yaml`
- fixed `--force` install on base/live environments

## v10.2.2 (date: 23.6.2021)

- adding `rcc robot dependencies` command for viewing desired execution
  environment dependencies
- same view is now also shown in run context replacing `pip freeze` if
  golden-ee.yaml exists in execution environment

## v10.2.1 (date: 21.6.2021)

- showing dependencies listing from environment before runs

## v10.2.0 (date: 21.6.2021)

- adding golden-ee.yaml document into holotree space (listing of components)

## v10.1.1 (date: 18.6.2021)

- taking micromamba 0.14.0 into use

## v10.1.0 (date: 17.6.2021)

- adding pager for `rcc man xxx` documents
- more trace printing on workflow setup
- added [D] and [T] markers for debug and trace level log entries
- when debug and trace log level is on, normal log entries are prefixed with [N]
- fixed rights problem in file `rcc_plan.log`

## v10.0.0 (date: 15.6.2021)

- removed lease support, this is major breaking change (if someone was using it)

## v9.20.0 (date: 10.6.2021)

- added `rcc task script` command for running anything inside robot environment

## v9.19.4 (date: 10.6.2021)

- added json format to `rcc holotree export` output formats
- added docs/recipes.md and also new command `rcc docs recipes`
- added links to README.md to internal documentation

## v9.19.3 (date: 10.6.2021)

- added support for getting list of events out
- fix: moved holotree changes messages to trace level

## v9.19.2 (date: 9.6.2021)

- added locking of holotree into environment restoring

## v9.19.1 (date: 8.6.2021)

- added locking of holotree into new environment building and recording

## v9.19.0 (date: 8.6.2021)

- added event journaling support (no user visible yet)
- added first event "space-used" in holotree restore operations (this enables
  tracking of all places where environments are created)

## v9.18.0 (date: 3.6.2021)

- now using holotree location from catalog, so that catalog decides where
  holotree is created (defaults to `ROBOCORP_HOME` but can be different)
- if hololib.zip exist, then `--space` flag must be given or run fails
- hololib.zip is now reported in robot diagnostics
- environment difference print is now (mostly) behind `--trace` flag
- if rcc is not interactive, color toggling on Windows is skipped
- micromamba download is now done "on demand" only
- added robot tests for hololib.zip workflow

## v9.17.2 (date: 2.6.2021)

- fixing broken tests, and taking account changed specifications

## v9.17.1 (date: 2.6.2021) broken

- adding supporting structures for zip based holotree runs [experimental]

## v9.17.0 (date: 26.5.2021)

- added `export` command to holotree [experimental]

## v9.16.0 (date: 21.5.2021)

- catalog extension based on operating system, architecture and directory
  location

## v9.15.1 (date: 21.5.2021)

- added images as non-executable files
- run and testrun commands have new option `--no-outputs` which prevent
  capture of stderr/stdout into files
- separated `--trace` and `--debug` flags from `micromamba` and `pip` verbosity
  introduced in v9.12.0 (it is causing too much output and should be reserved
  only for `RCC_VERBOSE_ENVIRONMENT_BUILDING` variable

## v9.15.0 (date: 20.5.2021)

- for `task run` and `task testrun` there is now possibility to give additional
  arguments from commandline, by using `--` separator between normal rcc
  arguments and those intended for executed robot
- rcc now considers "http://127.0.0.1" as special case that does not require
  https

## v9.14.0 (date: 19.5.2021)

- added PYTHONPATH diagnostics validation
- added `--production` flag to diagnostics commands

## v9.13.0 (date: 18.5.2021)

- micromamba upgrade to version 0.13.1
- activation script fix for windows environment

## v9.12.1 (date: 18.5.2021)

- new environment variable `ROBOCORP_OVERRIDE_SYSTEM_REQUIREMENTS` to make
  skip those system requirements that some users are willing to try
- first such thing is "long path support" on some versions of Windows

## v9.12.0 (date: 18.5.2021)

- new environment variable `RCC_VERBOSE_ENVIRONMENT_BUILDING` to make
  environment building more verbose
- with above variable and `--trace` or `--debug` flags, both micromamba
  and pip are run with more verbosity

## v9.11.3 (date: 12.5.2021)

- adding error signaling on anywork background workers
- more work on improving slow parts of holotree
- fixed settings.yaml conda link (conda.anaconda.org reference)

## v9.11.2 (date: 11.5.2021)

- added query cache in front of slow "has blueprint" query (windows)
- more timeline entries added for timing purposes

## v9.11.1 (date: 7.5.2021)

- new get/robot capabilitySet added into rcc
- added User-Agent to rcc web requests

## v9.11.0 (date: 6.5.2021)

- started using new capabilitySet feature of cloud authorization
- added metric for run/robot authorization usage
- one minor typo fix with "terminal" word

## v9.10.2 (date: 5.5.2021)

- added metrics to see when there was catalog failure (pre-check related)
- added PYTHONDONTWRITEBYTECODE=x setting into rcc generated environments,
  since this will pollute the cache (every compilation produces different file)
  without much of benefits
- also added PYTHONPYCACHEPREFIX to point into temporary folder
- added `--space` flag to `rcc cloud prepare` command

## v9.10.1 (date: 5.5.2021)

- added check for all components owned by catalog, to verify that they all
  are actually there
- added debug level logging on environment restoration operations
- added possibility to have line numbers on rcc produced log output (stderr)
- rcc log output (stderr) is now synchronized thru a channel
- made holotree command tree visible on toplevel listing

## v9.10.0 (date: 4.5.2021)

- refactoring code so that runs can be converted to holotree
- added `--space` option to runs so that they can use holotree
- holotree blueprint should now be unified form (same hash everywhere)
- holotree now co-exists with old implementation in backward compatible way

## v9.9.21 (date: 4.5.2021)

- documentation fix for toplevel config flag, closes #18

## v9.9.20 (date: 3.5.2021)

- added blueprint subcommand to holotree hierarchy to query blueprint
  existence in hololib

## v9.9.19 (date: 29.4.2021)

- refactoring to enable virtual holotree for --liveonly functionality
- NOTE: leased environments functionality will go away when holotree
  goes mainstream (and plan for that is rcc series v10)

## v9.9.18 (date: 28.4.2021)

- some cleanup on code base
- changed autoupdate url for Robocorp Lab

## v9.9.17 (date: 20.4.2021)

- added environment, workspace, and robot support to holotree variables command
- also added some robot tests for holotree to verify functionality

## v9.9.16 (date: 20.4.2021)

- added support for deleting holotree controller spaces
- added holotree and hololib to full environment cleanup
- added required parameter to `rcc env delete` command also

## v9.9.15 (date: 19.4.2021)

- bugfix: locking while multiple rcc are doing parallel work should now
  work better, and not corrupt configuration (so much)

## v9.9.14 (date: 15.4.2021)

- environment variables conda.yaml ordering fix (from robot.yaml first)
- task shell does not need task specified anymore

## v9.9.13 (date: 15.4.2021)

- fixing environment variables bug from below

## v9.9.12 (date: 15.4.2021)

- updated rpaframework to version 9.5.0 in templates
- added more timeline entries around holotree
- minor performance related changes for holotree
- removed default PYTHONPATH settings from "taskless" environment
- known, remaining bug: on "env variables" command, with robot without default
  task and without task given in CLI, environment wont have PATH or PYTHONPATH
  or robot details setup correctly

## v9.9.11 (date: 13.4.2021)

- added support for listing holotree controller spaces

## v9.9.10 (date: 12.4.2021)

- removed index.py utility, since better place is on other repo, and it
  was mistake to put it here

## v9.9.9 (date: 9.4.2021)

- fixed index.py utility tool to work in correct repository

## v9.9.8 (date: 9.4.2021)

- skip environment bootstrap when there is no conda.yaml used
- added index.py utility tool for generating index.html for S3

## v9.9.7 (date: 8.4.2021)

- now `rcc holotree bootstrap` can only download templates with `--quick`
  flag, or otherwise also prepare environment based on that template

## v9.9.6 (date: 8.4.2021)

- holotree note: in this series 9, holotree will remain experimental and
  will not be used for production yet
- added separate `holotree` subtree in command structure (it is not internal
  anymore, but still hidden)
- partial implementations of holotree variables and bootstrap commands
- settings.yaml version 2021.04 update: now there is separate section
  for templates
- profiling option `--pprof` is now global level option
- improved error message when rcc is not configured yet

## v9.9.5 (date: 6.4.2021)

- micromamba upgrade to version 0.9.2

## v9.9.4 (date: 6.4.2021)

- fix for holotree change detection when switching blueprints

## v9.9.3 (date: 1.4.2021)

- added export/SET prefix to `rcc env variables` command
- updated README.md with patterns to version numbered releases
- known bug: holotree does not work correctly yet -- DO NOT USE

## v9.9.2 (date: 1.4.2021)

- more holotree integration work to get it more experimentable

## v9.9.1 (date: 31.3.2021)

- Github Actions upgrade to use Go 1.16 for rcc compilation

## v9.9.0 (date: 31.3.2021) broken

- added holotree as part of source code (but not as integrated part yet)
- added new internal command: holotree

## v9.8.11 (date: 30.3.2021)

- added Accept header to micromamba download command
- made some URL diagnostics optional, if they are left empty

## v9.8.10 (date: 30.3.2021)

- fix: no more panics when directly writing to settings.yaml

## v9.8.9 (date: 29.3.2021)

- added `cloud-ui` to settings.yaml

## v9.8.8 (date: 29.3.2021)

- mixed fixes and experiments edition
- ignoring empty variable names on environment dumps, closes #17
- added some missing content types to web requests
- added experimental ephemeral ECC implementation
- more common timeline markers added
- will not list pip dependencies on assistant runs
- will not ask cloud for runtime authorization (bug fix)

## v9.8.7 (date: 26.3.2021)

- more finalization of settings.yaml change
- made micromamba less quiet on environment building
- secrets now have write access enabled in rcc authorization requests
- if merged conda.yaml files do not have names, merge result wont have either

## v9.8.6 (date: 25.3.2021)

- settings.yaml cleanup
- fixed robot tests for 9.8.5 template changes

## v9.8.5 (date: 24.3.2021)

- Robot templates updated: Rpaframework updated to v9.1.0
- Robot templates updated: Improved task names
- Robot templates updated: Extended template has example of multiple tasks execution

## v9.8.4 (date: 24.3.2021)

- fix for pip made too silent on this v9.8.x series
- and also in failure cases, print out full installation plan

## v9.8.3 (date: 24.3.2021)

- can configure all rcc operations not to verify correct SSL certificate
  (please note, doing this is insecure and allows man-in-the-middle attacks)
- applied reviewed changes to what is actually in settings.yaml file

## v9.8.2 (date: 23.3.2021)

- ALPHA level pre-release (do not use, unless you know what you are doing)
- reorganizing some code to allow better use of settings.yaml
- more values from settings.yaml are now used

## v9.8.1 (date: 22.3.2021)

- ALPHA level pre-release (do not use, unless you know what you are doing)
- now some parts of settings are used from settings.yaml
- settings.yaml is now critical part of rcc, so diagnostics also contains it
- also from now, problems in settings.yaml may make rcc to fail
- changed ephemeral key size to 2048, which should be good enough

## v9.8.0 (date: 18.3.2021)

- ALPHA level pre-release with settings.yaml (do not use, unless you know
  what you are doing)
- started to moved some of hardcoded things into settings.yaml (not used yet)
- minor assistant upload fix, where one error case was not marked as error

## v9.7.4 (date: 17.3.2021)

- typo fix pull request from jaukia
- added micromamba --no-rc flag

## v9.7.3 (date: 16.3.2021)

- upgrading micromamba dependency to 0.8.2 version
- added .robot, .csv, .yaml, .yml, and .json in non-executable fileset
- also added "dot" files as non-executable
- added timestamp update to copyfile functionality
- added toplevel --tag option to allow semantic tagging for client
  applications to indicate meaning of rcc execution call

## v9.7.2 (date: 11.3.2021)

- adding visibility of installation plans in environment listing
- added --json support to environment listing including installation plan file
- added command `rcc env plan` to show installation plans for environment
- installation plan is now also part of robot diagnostics, if available

## v9.7.1 (date: 10.3.2021)

- fixes/improvements to activation and installation plan
- added missing content type to assistant requests
- micromamba upgrade to 0.8.0

## v9.7.0 (date: 10.3.2021)

- conda environments are now activated once on creation, and variables go
  with environment, as `rcc_activate.json`
- there is also now new "installation plan" file inside environment, called
  `rcc_plan.log` which contains events that lead to activation
- normal runs are now more silent, since details are moved into "plan" file

## v9.6.2 (date: 5.3.2021)

- fix for time formats used in timeline, some metrics, and stopwatch

## v9.6.1 (date: 3.3.2021)

- refactored code use common.When as consistent timestamp for current rcc run

## v9.6.0 (date: 3.3.2021)

- new command `rcc cloud prepare` to support installing assistants on
  local computer for faster startup time
- added more timeline entries on relevant parts

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
