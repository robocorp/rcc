# Tips, tricks, and recipies


## How to see dependency changes?

Since version 10.2.2, rcc can show dependency listings using
`rcc robot dependencies` command. Listing always have two sided, "Wanted"
which is content from dependencies.yaml file, and "Available" which is from
actual environment command was run against. Listing is also shown during
robot runs.

### Why is this important?

- as time passes and world moves forward, new version of used components
  (dependencies) are released, and this may cause "configuration drift" on
  your robots, and without tooling in place, this drift might go unnoticed
- if your dependencies are not fixed, there will be configuration drift and
  your robot may change behaviour (become buggy) when dependency changes and
  goes against implemented robot
- even if you fix your dependencies in `conda.yaml`, some of those components
  or their components might have floating dependencies and they change your
  robots behaviour
- if your execution environment is different from your development environment
  then there might be different versions available for different operating
  systems
- if dependency resolution algorithm changes (pip for example) then you might
  get different environment with same `conda.yaml`
- when you upgrade one of your dependencies (for example, rpaframework) to new
  version, dependency resolution will now change, and now listing helps you
  understand what has changed and how you need to change your robot
  implementation because of that

### Example of dependencies listing from holotree environment

```sh
# first list dependencies from execution environment
rcc robot dependencies --space user

# if everything looks good, export it as wanted dependencies.yaml
rcc robot dependencies --space user --export

# and verify that everything looks `Same`
rcc robot dependencies --space user
```


## How to freeze dependencies?

Starting from rcc 10.3.2, there is now possibility to freeze dependencies.
This is how you can experiment with it.

### Steps

- have your `conda.yaml` to contain only those dependencies that your robot
  needs, either with exact versions or floating ones
- run robot in your target environment at least once, so that environment
  there gets created
- from that run's artifact directory, you should find file that has name
  something like `environment_xxx_yyy_freeze.yaml`
- copy that file back into your robot, right beside existing `conda.yaml`
  file (but do not overwrite it, you need that later)
- edit your `robot.yaml` file at `condaConfigFile` entry, and add your
  newly copied `environment_xxx_yyy_freeze.yaml` file there if it does not
  already exist there
- repackage your robot and now your environment should stay quite frozen

### Limitations

- this is new and experimental feature, and we don't know yet how well it
  works in all cases (but we love to get feedback)
- currently this freezing limits where robot can be run, since dependencies
  on different operating systems and architectures differ and freezing cannot
  be done in OS and architecture neutral way
- your robot will break, if some specific package is removed from pypi or
  conda repositories
- your robot might also break, if someone updates package (and it's dependencies)
  without changing its version number
- for better visibility on configuration drift, you should also have
  `dependencies.yaml` inside your robot (see other recipe for it)


## How pass arguments to robot from CLI?

Since version 9.15.0, rcc supports passing arguments from CLI to underlying
robot. For that, you need to have task in `robot.yaml` that co-operates with
additional arguments appended at the end of given `shell` command.

### Example robot.yaml with scripting task

```yaml
tasks:
  Run all tasks:
    shell: python -m robot --report NONE --outputdir output --logtitle "Task log" tasks.robot

  scripting:
    shell: python -m robot --report NONE --outputdir output --logtitle "Scripting log"

condaConfigFile: conda.yaml
artifactsDir: output
PATH:
  - .
PYTHONPATH:
  - .
ignoreFiles:
  - .gitignore
```

### Run it with `--` separator.

```sh
rcc task run --interactive --task scripting -- --loglevel TRACE --variable answer:42 tasks.robot
```


## How to run any command inside robot environment?

Since version 9.20.0, rcc now supports running any command inside robot space
using `rcc task script` command.

### Some example commands

Run following commands in same direcotry where your `robot.yaml` is. Or
otherwise you have to provide `--robot path/to/robot.yaml` in commandline.

```sh
# what python version we are running
rcc task script --silent -- python --version

# get pip list from this environment
rcc task script --silent -- pip list

# start interactive ipython session
rcc task script --interactive -- ipython
```


## How to convert existing python project to rcc?

### Basic workflow to get it up and running

1. Create a new robot using `rcc create` with `Basic Python template`.
2. Remove task.py and and copy files from your existing project to this new
   rcc/robot project.
3. Discover all your publicly available dependencies (including your python
   version) and try find as many as possible from https://anaconda.org/conda-forge/
   and take rest from https://pypi.org/ and put those dependencies
   into `conda.yaml`. And remove all those dependencies that you do not actually
   need in your project.
4. Do not add any private dependencies into `conda.yaml`, and also no passwords
   in that `conda.yaml` either (passwords belong to secure place, like Vault).
5. Modify your `robot.yaml` task definitions so, that it is how your python
   project should be executed.
6. If you have additional private libraries, put them inside robot directory
   structure (like under `libraries` or something similar) and edit PYTHONPATH
   settings in `robot.yaml` to include those paths (relative paths only).
7. If you have additional scripts/small binaries that your robot dependes on,
   add them inside robot directory structure (like under `scripts` directory)
   and edit PATH settings in `robot.yaml` to include that (relative) path.
8. If your python project needs external dependencies (like Word or Excel)
   then those dependencies must be present in machine where robot is executed
   and they are not part of this conversion.
9. Run robot and test if it works, and iterate to make needed changes.

### What next?

* Your python project is now converted to rcc and should be locally "runnable".
* Setup Assistant or Worker in your machine and create Assistant or Robot
  in Robocorp Control Room, and try to run it from there.
* If your robot is "headless", has all dependencies, and should be runnable
  in Linux, then you can try to run it in container from Control Room.
* If your project is python2 project, then consider converting it to python3.
* If you want to use `rpaframework` in your robot (like dialogs for example),
  then you have to start converting to use those features in your code.
* etc.


## Is rcc limited to Python and Robot Framework?

Absolutely not! Here is something completely different for you to think about.

Lets assume, that you are in almost empty Linux machine, and you have to
quickly build new micromamba in that machine. Hey, there is `bash`, `$EDITOR`,
and `curl` here.  But there are no compilers, git, or even python installed.

> Pop quiz, hot shot! Who you gonna call? MacGyver!

### This is what we are going to do ...

Here is set of commands we are going to execute in our trusty shell

```sh
mkdir -p builder/bin
cd builder
$EDITOR robot.yaml
$EDITOR conda.yaml
$EDITOR bin/builder.sh
curl -o rcc https://cdn.sema4.ai/rcc/releases/v18.5.0/linux64/rcc
chmod 755 rcc
./rcc run -s MacGyver
```

### Write a robot.yaml

So, for this to be a robot, we need to write heart of our robot, which is
`robot.yaml` of course.

```yaml
tasks:
  µmamba:
    shell: builder.sh
condaConfigFile: conda.yaml
artifactsDir: output
PATH:
- bin
```

### Write a conda.yaml

Next, we need to define what our robot needs to be able to do our mighty task.
This goes into `conda.yaml` file.

```yaml
channels:
- conda-forge
dependencies:
- git
- gmock
- cli11
- cmake
- compilers
- cxx-compiler
- pybind11
- libsolv
- libarchive
- libcurl
- gtest
- nlohmann_json
- cpp-filesystem
- yaml-cpp
- reproc-cpp
- python=3.8
- pip=20.1
```

### Write a bin/builder.sh

And finally, what does our robot do. And this time, this goes to our directory
bin, which is on our PATH, and name for this "robot" is actually `builder.sh`
and it is a bash script.

```sh
#!/bin/bash -ex

rm -rf target output/micromamba*
git clone https://github.com/mamba-org/mamba.git target
pushd target
version=$(git tag -l --sort='-creatordate' | head -1)
git checkout $version
mkdir -p build
pushd build
cmake .. -DCMAKE_INSTALL_PREFIX=/tmp/mamba -DENABLE_TESTS=ON -DBUILD_EXE=ON -DBUILD_BINDINGS=OFF
make
popd
popd
mkdir -p output
cp target/build/micromamba output/micromamba-$version
```


## Think what you can do with this conda.yaml?

```
channels:
  # Just using conda-forge, nothing else.
  - conda-forge

dependencies:
  # I'm not going to have python directly installed here ..
  # But let's go wild with conda-forge ...

  - nginx=1.21.6     # https://anaconda.org/conda-forge/nginx
  - php=8.1.5        # https://anaconda.org/conda-forge/php
  - go=1.17.8        # https://anaconda.org/conda-forge/go
  - postgresql=14.2  # https://anaconda.org/conda-forge/postgresql
  - terraform=1.1.9  # https://anaconda.org/conda-forge/terraform
  - awscli=1.23.9    # https://anaconda.org/conda-forge/awscli
  - firefox=100.0    # https://anaconda.org/conda-forge/firefox
```

## How to control holotree environments?

There is three controlling factors for where holotree spaces are created.

First is location of `ROBOCORP_HOME` at creation time of environment. This
decides general location for environment and it cannot be changed or relocated
afterwards.

Second controlling factor is given using `--controller` option and default for
this is value `user`. And when applications are calling rcc, they should
have their own "controller" identity, so that all spaces created for one
application are groupped together by prefix of their "space" identity name.

Third controlling factor is content of `--space` option and again default
value there is `user`. Here it is up to user or application to decide their
strategy of use of different names to separate environments to their logical
used partitions. If you choose to use just defaults (user/user) then there
is going to be only one real environment available.

But above three controls gives you good ways to control how you and your
applications manage their usage of different python environments for
different purposes. You can share environments if you want, but you can also
give a dedicated space for those things that need full control of their space.

So running following commands demonstrate different levels of control for
space creation.

```
export ROBOCORP_HOME=/tmp/rchome
rcc holotree variables simple.yaml
rcc holotree variables --space tips simple.yaml
rcc holotree variables --controller tricks --space tips simple.yaml
```

If you now run `rcc holotree list` it should list something like following.

```
Identity            Controller  Space  Blueprint         Full path
--------            ----------  -----  --------          ---------
5a1fac3c5_2daaa295  rcc.user    tips   c34ed96c2d8a459a  /tmp/rchome/holotree/5a1fac3c5_2daaa295
5a1fac3c5_9fcd2534  rcc.user    user   c34ed96c2d8a459a  /tmp/rchome/holotree/5a1fac3c5_9fcd2534
9e7018022_2daaa295  rcc.tricks  tips   c34ed96c2d8a459a  /tmp/rchome/holotree/9e7018022_2daaa295
```

### How to get understanding on holotree?

See: https://github.com/robocorp/rcc/blob/master/docs/environment-caching.md

### How to activate holotree environment?

On Linux/MacOSX:

```sh
# full robot environment
source <(rcc holotree variables --space mine --robot path/to/robot.yaml)

# or with just conda.yaml
source <(rcc holotree variables --space mine path/to/conda.yaml)
```

On Windows

```sh
rcc holotree variables --space mine --robot path/to/robot.yaml > mine_activate.bat
call mine_activate.bat
```

You can also try

```sh
rcc task shell --robot path/to/robot.yaml
```


## What is `ROBOCORP_HOME`?

It is environment variable level settings, that says where Robocorp tooling
can keep tooling specific files and configurations. It has default values,
and normal case is that defaults are fine. But if there is need to "relocate"
that somewhere else, then this environment variable does the trick.

### Are there some rules for `ROBOCORP_HOME` variable?

- go with defaults, unless you have very good reason to override it
- avoid using spaces or special characters in path that is `ROBOCORP_HOME`,
  so stick to basic english letters and numbers
- never use your "home" directory as `ROBOCORP_HOME`, it will cause conflicts
- never share `ROBOCORP_HOME` between two users, it should be unique to each
  different user account
- also keep it private and protected, other users should not have access
  to that directory
- never use `ROBOCORP_HOME` as working directory for user, or any other
  tools; this directory is only meant for Robocorp tooling to use, change,
  and operate on
- never put `ROBOCORP_HOME` on network drive, since those tend to be slow,
  and using those can cause real performance issues
- always make sure, that user owning that `ROBOCORP_HOME` directory has full
  control access and permissions to everything inside that directory structure


### When you might actually need to setup `ROBOCORP_HOME`?

- if your username contains spaces, or some special characters that can cause
  tooling to break
- if path to your home directory is very long, it might cause long path  issues,
  and one way to go around is have `ROBOCORP_HOME` on shorter path
- if you need to have `ROBOCORP_HOME` on some different disk than default
- if your home directory is on HDD drive (or even network drive), but you
  have fast SSD direve available, performance might be much better on SSD

## What is shared holotree?

Shared holotree is way to multiple users use same environment blueprint in
same machine, or even in different machines with same, once it is built or
imported into hololib.

## How to setup rcc to use shared holotree?

### One time setup

On each machine, where you want to use shared holotree, the shared location
needs to be enabled once. 
This depends on the operating system so the commands below are OS specific
and do require elevated rights from the user that runs them. 

The commands to enable the shared locations are:
* Windows: `rcc holotree shared --enable`
  * Shared location: `C:\ProgramData\robocorp`
* MacOS: `sudo rcc holotree shared --enable`
  * Shared location: `/Users/Shared/robocorp`
* Linux: `sudo rcc holotree shared --enable`
  * Shared location: `/opt/robocorp`

Note: On Windows the command below assumes the standard `BUILTIN\Users`
user group is present.
If your organization has replaced this you can grant the permission with:

```
icacls "C:\ProgramData\robocorp" /grant "*S-1-5-32-545:(OI)(CI)M" /T
```

To switch the user to using shared holotrees use the following command.

```sh
rcc holotree init
```

### Reverting back to private holotrees

If user wants to go back to private holotrees, they can run following command.

```sh
rcc holotree init --revoke
```

## What can be controlled using environment variables?

- `ROBOCORP_HOME` points to directory where rcc keeps most of Robocorp related
  files and directories are kept
- `ROBOCORP_OVERRIDE_SYSTEM_REQUIREMENTS` makes rcc more relaxed on system
  requirements (like long path support requirement on Windows) but it also
  means that if set, responsibility of resolving failures are on user side
- `RCC_VERBOSE_ENVIRONMENT_BUILDING` makes environment creation more verbose,
  so that failing environment creation can be seen with more details
- `RCC_CREDENTIALS_ID` is way to provide Control Room credentials using
  environment variables
- `RCC_NO_BUILD` with any non-empty value will prevent rcc for creating
  new environments (also available as `--no-build` CLI flag, and as
  an option in `settings.yaml` file)
- `RCC_VERBOSITY` controls how verbose rcc output will be. If this variable
  is not set, then verbosity is taken from `--silent`, `--debug`, and `--trace`
  CLI flags. Valid values for this variable are `silent`, `debug` and `trace`.
- `RCC_NO_TEMP_MANAGEMENT` with any non-empty value will prevent rcc for
  doing any management in relation to temporary directories; using this
  environment variable means, that something else is managing temporary
  directories life cycles (and this might also break environment isolation)
- `RCC_NO_PYC_MANAGEMENT` with any non-empty value will prevent rcc for
  doing any .pyc file management; using this environment variable means, that
  something else is doing that management (and using this makes rcc slower
  and hololibs become bigger and grow faster, since .pyc files are unfriendly
  to caching)


## How to troubleshoot rcc setup and robots?

```sh
# to get generic setup diagnostics
rcc configure diagnostics

# to get robot and environment setup diagnostics
rcc configure diagnostics --robot path/to/robot.yaml

# to see how well rcc performs in your machine
rcc configure speedtest
```

### Additional debugging options

- generic flag `--debug` shows debug messages during execution
- generic flag `--trace` shows more verbose debugging messages during execution
- flag `--timeline` can be used to see execution timeline and where time was spent
- with option `--pprof <filename>` enable profiling if performance is problem,
  and want to help improve it (by submitting that profile file to developers)

## Advanced network diagnostics

When using custom endpoints or just needing more control over what network
checks are done, command `rcc configure netdiagnostics` may become helpful.

```sh
# to test advanced network diagnostics with defaults
rcc configure netdiagnostics

# to capture advanced network diagnostics defaults to new configuration file
rcc configure netdiagnostics --show > path/to/modified.yaml

# to test advanced network diagnostics with custom tests
rcc configure netdiagnostics --checks path/to/modified.yaml
```

### Configuration

- get example configuration out using `--show` option (as seen above)
- configuration file format is YAML
- add or remove points to DNS, HTTP HEAD and GET methods
- `url:` and `codes:` are required fields for HEAD and GET checks
- `codes:` field is list of acceptable HTTP response codes
- `content-sha256` is optional, and provides additional confidence when content
  is static and result content hash can be calculated (using sha256 algorithm)

## What is in `robot.yaml`?

### Example

```yaml
tasks:
  Just a task:
    robotTaskName: Just a task
  Version command:
    shell: python -m robot --version
  Multiline command:
    command:
      - python
      - -m
      - robot
      - --report
      - NONE
      - -d
      - output
      - --logtitle
      - Task log
      - tasks.robot

devTasks:
  Editor setup:
    shell: python scripts/editor_setup.py
  Repository update:
    shell: python scripts/repository_update.py

condaConfigFile: conda.yaml

environmentConfigs:
- environment_linux_amd64_freeze.yaml
- environment_windows_amd64_freeze.yaml
- common_linux_amd64.yaml
- common_windows_amd64.yaml
- common_linux.yaml
- common_windows.yaml
- conda.yaml

preRunScripts:
- privatePipInstall.sh
- initializeKeystore.sh

artifactsDir: output

ignoreFiles:
- .gitignore

PATH:
- .
- bin

PYTHONPATH:
- .
- libraries
```

### What is this `robot.yaml` thing?

It is declarative description in [YAML format](https://en.wikipedia.org/wiki/YAML)
of what robot is and what it can do.

It is also a pointer to "a robot center of a universe" for directory it resides.
So it is marker of "current working folder" when robot starts to execute and
that will be indicated in `ROBOT_ROOT` environment variable. All declarations
inside `robot.yaml` should be relative to and inside of this location, so do
not use absolute paths here, or relative references to any parent directory.

It also marks root location that gets wrapped into `robot.zip` when either
wrapping locally or pushing to Control Room. Nothing above directory holding
`robot.yaml` gets wrapped into that zip file.

Also note that `robot.yaml` is just a name of a file. Other names can be used
and then given to commands using `--robot othername.yaml` CLI option. But
in Robocorp tooling, this default name `robot.yaml` is used to have common
ground without additional configuration needs.

### Why "the center of the universe"?

Firstly, it is not "the center", it is just "a center of a universe" for
specific robot. So it only applies to that specific robot, when operations
are done around that one specific robot. Other robots have their own centers.

And reason for thinking this way is, that it is "convention over configuration",
meaning that when we have this concept, there is much less configuration to do.
It gives following things automatically, without additional configuration:

- what is "root" folder, when wrapping robot into deliverable package
- what is starting working directory when robot is executed (robot itself can
  of course change its working directory freely while running)
- it gives solid starting point for relative paths inside robot, so that
  PATH, PYTHONPATH, artifactsDir, and other relative references can be
  converted absolute ones
- it allows robot location to be different for different users and on different
  machines, and still have everything declared with known (but relative)
  locations

### What are `tasks:`?

One robot can do multiple tasks. Each task is a single declaration of named
task that robot can do.

There are three types of task declarations:

1. The `robotTaskName` form, which is simplest and there only name of a task
   is given. In above example `Just a task` is a such thing. This is Robot
   Framework specific form.
2. The `shell` form, where full CLI command is given as oneliner. In above
   example, `Version command` is example of this.
3. The `command` form is oldest. It is given as list of command and its
   arguments, and it is most accurate way to declare CLI form, but it is also
   most spacious form.

### What are `devTasks:`?

They are tasks like above `tasks:` define. But they have two major differences
compared to normal `tasks:` definitions:

1. They are for developers at development machines, for doing development time
   activities and tasks. They should never be available in cloud containers,
   Assistants or Agents. Developer tools can provide support for them, but
   their semantics should be only valid in development context.
2. They can be run like normal tasks, by providing `--dev` flag. But during
   their run, all `preRunScripts:` are ignored. Otherwise environment is
   created and managed as with normal tasks, but without pre-run scripts
   applied.

The `devTasks:` primary goal is to provide developers a way to use the same
tooling to automate their development process as normal `tasks:` provide ways
to automate robot actions. Some examples could be: common editor setups,
version control repository updates.

Currently `--dev` option is only available for `rcc run` and `rcc task run`
commands. With the `--dev` option the only available tasks for execution will
be the `devTasks:`. The normal `tasks:` will be skipped/missing. If the `--dev`
option is missing, the `devTasks:` will be skipped/missing, and the normal
`tasks:` will be the ones available for execution.

### What is `condaConfigFile:`?

> Use of this is deprecated, please use `environmentConfigs:` instead.

This is actual name used as `conda.yaml` environment configuration file.
See next topic about details of `conda.yaml` file.
This is just single file that describes dependencies for all operating systems.
For more versatile selection, see `environmentConfigs` below. If that
`environmentConfigs` exists and one of those files matches machine running
rcc, then this config is ignored.

### What are `environmentConfigs:`?

These are like condaConfigFile above, but as priority list form. First matching
and existing item from that list is used as environment configuration file.

These files are matched by operating system (windows/darwin/linux) and by
architecture (amd64/arm64). If filename contains word "freeze", it must
match OS and architecture exactly. Other variations allow just some or none
of those parts.

And if there is no such file, then those entries are just ignored. And if
none of files match or exist, then as final resort, `condaConfigFile` value
is used if present.

### What are `preRunScripts:`?

This is set of scripts or commands that are run before actual robot task
execution. Idea with these scripts is that they can be used to customize
runtime environment right after it has been restored from hololib, and just
before actual robot execution is done.

If script names contains some of "amd64", "arm64", "darwin", "windows" and/or
"linux" words (like `script_for_amd64_linux.sh`) then other architectures
and operating systems will skip those scripts, and only amd64 linux systems
will execute them.

All these scripts are run in "robot" context with all same environment
variables available as in robot run.

These scripts can pollute the environment, and it is ok. Next rcc operation
on same holotree space will first do the cleanup though.

All scripts must be executed successfully or otherwise full robot run is
considered failure and not even tried. Scripts should use exit code zero
to indicate success and everything else is failure.

Some ideas for scripts could be:

- install custom packages from private pip repository
- use Vault secrets to prepare system for actual robot run
- setup and customize used tools with secret or other private details that
  should not be visible inside hololib catalogs (public caches etc)

### What is `artifactsDir:`?

This is location of technical artifacts, like log and freezefiles, that are
created during robot execution and which can be used to find out technical
details about run afterwards. Do not confuse these with work-item data, which
are more business related and do not belong here.

During robot run, this locations is available using `ROBOT_ARTIFACTS`
environment variable, if you want to store some additional artifacts there.

### What are `ignoreFiles:`?

This is a list of configuration files that rcc uses as locations for ignore
patterns used while wrapping robot into a robot.zip file. But note, that once
filename is on this list, it must also be present on directory structure, this
is part of a contract.

Content of those files should be similar to what is used normally as version
control systems as ignore files (like .gitignore file in git context).
Here rcc implements only subset of functionality, and allows just mostly
globbing patterns or exact names of files and directories.

Note: do not put file or directory names that you want to be ignored directly
in this list. They all should reside in one of those configurations listed in
this configuration list.

Tip: using `.gitignore` as one of those `ignoreFiles:` entries helps you to
remove duplication of maintenance pressures. But if you want ignore different
things in git and in robot.zip, or if there are conflicts between those,
feel free use different filenames as you see fit.

### What are `PATH:`?

This allows adding entries into `PATH` environment variable. Intention
is to allow something like `bin` directory inside robot, where custom
scripts and binaries can be located and available for execution during
robot run.

### What are `PYTHONPATH:`?

This allows adding entries into `PYTHONPATH` environment variable. Intention
is to allow something like `libraries` directory inside robot, where custom
libraries can be located and automatically loaded by python and robot.


## What is in `conda.yaml`?

### Example

```yaml
channels:
- conda-forge

dependencies:
- python=3.9.13
- nodejs=16.14.2
- pip=22.1.2
- pip:
  - robotframework-browser==12.3.0
  - rpaframework==15.6.0

rccPostInstall:
  - rfbrowser init
```

### What is this `conda.yaml` thing?

It is declarative description in [YAML format](https://en.wikipedia.org/wiki/YAML)
of environment that should be set up.

### What are `channels:`?

Channels are conda sources where to get packages to be used in setting up
environment. It is recommended to use `conda-forge` channel, but there are
others also. Other recommendation is that only one channel is used, to get
consistently build environments.

Channels should be in priority order, where first one has highest priority.

Example above uses `conda-forge` as its only channel.
For more details about conda-forge, see this [link.](https://anaconda.org/conda-forge)

### What are `dependencies:`?

These are libraries that are needed to be installed in environment that is
declared in this `conda.yaml` file. By default they come from locations
setup in `channels:` part of file.

But there is also `- pip:` part and those dependenies come from
[PyPI](https://pypi.org/) and they are installed after dependencies from
`channels:` have been installed.

In above example, `python=3.9.13` comes from `conda-forge` channel.
And `rpaframework==15.6.0` comes from [PyPI](https://pypi.org/project/rpaframework/).

### What are `rccPostInstall:` scripts?

Once environment dependencies have been installed, but before it is frozen as
hololib catalog, there is option to run some additional commands to customize
that environment. It is list of "shell" commands that are executed in order,
and if any of those fail, environment creation will fail.

All those scripts must come from package declared in `dependencies:` section,
and should not use any "local" knowledge outside of environment under
construction. This makes environment creation repeatable and cacheable.

Do not use any private or sensitive information in those post install scripts,
since result of environment build could be cached and visible to everybody
who has access to that cache. If you need to have private or sensitive packages
in your environment, see `preRunScripts` in `robot.yaml` file.


## How to do "old-school" CI/CD pipeline integration with rcc?

If you have CI/CD pipeline and want to updated your robots from there, this
recipe should give you ideas how to do it. This example works in linux, and
you probably have to modify it to work on Mac or Windows, but idea will be same.

Basic requirements are:
- have well formed robot in version control
- have rcc command available or possibility to fetch it
- possibility on CI/CD pipeline to run just simple CLI commands

### The oldschoolci.sh script

```sh
#!/bin/sh -ex

curl -o rcc https://cdn.sema4.ai/rcc/releases/v18.5.0/linux64/rcc
chmod 755 rcc
./rcc cloud push --account ${ACCOUNT_ID} --directory ${ROBOT_DIR} --workspace ${WORKSPACE_ID} --robot ${ROBOT_ID}
```

So above script uses `curl` command to download rcc from download site, and
makes it executable. And then it simply calls that `rcc` command, and expects
that CI system has provided few variables.

### A setup.sh script for simulating variable injection.

```sh
#!/bin/sh

export ACCOUNT_ID=4242:cafe9d9c0dadag00d37b9577babe1575b67bc1bbad3ce9484dead36a649c865beef26297e67c8d94f0f0057f0100ab64:https://api.eu1.robocorp.com
export WORKSPACE_ID=1717
export ROBOT_ID=2121
export ROBOT_DIR=$(pwd)/therobot
```

Expectations for above setup are:
- robot to be updated is in EU1 (behind https://api.eu1.robocorp.com API)
- Control Room account has "Access creadentials" 4242 available and active
- account has access to workspace 1717
- there exist previously created robot 2121 in that workspace
- robot is located in "therobot" directory directly under "current working
  directory" (centered around `robot.yaml` file)
- and account has suitable rights to actually push robot to Control Room

### Simulating actual CI/CD step in local machine.

```sh
#!/bin/sh -ex

source setup.sh
./oldschoolci.sh
```

Above script brings "setup" and "old school CI" together, but just for
demonstration purposes. For real life use, adapt and remember security (no
compromising variable content inside repository).

### Additional notes

- if CI/CD worker/container can be custom build, then it is recommended to
  download rcc just once and not on every run (like oldschoolci.sh script now
  does)
- that `ACCOUNT_ID` should be stored in credentials store/vault in CI system,
  because that is secret that you need to use to be able to push to cloud
- that `ACCOUNT_ID` is "ephemeral" account, and will not be saved in `rcc.yaml`
- also consider saving other variables in secure way
- in actual CI/CD pipeline, you might want to embed actual commands into
  CI step recipe and not have external scripts (but you decide that)


## How to setup custom templates?

Custom templates allows making your own templates that can be used when
new robot is created. So if you have your own standard way of doing things,
then custom template is good way to codify it.

You then need to do these steps:

- setup custom settings.yaml that point location where template configuration
  file is located (the templates.yaml file)
- if you are using profiles, then make above change in settings.yaml used there
- create that custom templates.yaml configuration file that lists available
  templates, and where template bundle can be found (the templates.zip file)
- and finally build that templates.zip to bundle together all those templates
  that were listed in configuration file
- and finally both templates.yaml and templates.zip must be somewhere behind
  URL that starts with https:

Note: templates are needed only on development context, and they are not used
or needed in Assistant or Worker context.

### Custom template configuration in `settings.yaml`.

In settings.yaml, there is `autoupdates:` section, and there is entry for
`templates:` where you should put exact name and location where active
templates configuration file is located.

Example:

```yaml
autoupdates:
  templates: https://special.acme.com/robot/templates-1.0.1.yaml
```

As above example shows, name is configurable, and can even contain some
versioning information, if so needed.

### Custom template configuration file as `templates.yaml`.

In that `templates.yaml` following things must be provided:

- `hash:` (sha256) of "templates.zip" file (so that integrity of templates.zip
  can be verified)
- `url:` to exact name and location where that templates.zip can be downloaded
- `date:` when this template.yaml file was last updated
- `templates:` as key/value pairs of templates and their "one liner"
  description seen in UIs
- so, if there is `shell.zip` inside templates.zip, then that should have
  `shell: Shell Robot Template` or something similar in that `templates:`
  section

Example:

```yaml
hash: c7b1ba0863d9f7559de599e3811e31ddd7bdb72ce862d1a033f5396c92c5c4ec
url: https://special.acme.com/robot/templates-1.0.1.zip
date: 2022-09-12
templates:
  shell: Simple Shell Robot template
  extended: Extended Robot Framework template
  playwright: Playwright template
  producer-consumer: Producer-consumer model template
```

### Custom template content in `templates.zip` file.

Then that `templates.zip` is zip-of-zips. So for each key from templates.yaml
`templates:` sections should have matching .zip file inside that master zip.

### Shared using `https:` protocol ...

Then both `templates.yaml` and `templates.zip` should be hosted somewhere
which can be accessed using https protocol. Names there should match those
defined in above steps.

And that `settings.yaml` should either be delivered standalone into those
developer machines that need to use those templates, or better yet, be part
of "profile" that developers can use to setup all of required configurations.


## Where can I find updates for rcc?

https://cdn.sema4.ai/rcc/releases/index.html

That is rcc download site with two categories of:
- tested versions (these are ones we ship with our tools)
- latest 20 versions (which are not battle tested yet, but are bleeding edge)


## What has changed on rcc?

### See changelog from git repo ...

https://github.com/robocorp/rcc/blob/master/docs/changelog.md

### See that from your version of rcc directly ...

```sh
rcc docs changelog
```


## Can I see these tips as web page?

Sure. See following URL.

https://github.com/robocorp/rcc/blob/master/docs/recipes.md

