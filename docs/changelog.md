# rcc change log

## v11.12.2 (date: 5.5.2022)

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
- Old `rcc env cleanup` was renamed to `rcc configuration cleanup`  and
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
- added __MACOSX as ignored files/directories

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
