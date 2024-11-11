# Vocabulary

## Blueprint

Is unique identity calculated from `conda.yaml` or in general from some
`environment.yaml` file after it is formatted in canonical, unified form.
Currently `rcc` uses "siphash" for that. This is form of a fingerprint.

## Catalog

Catalog is description how final created environment should look like, after
it is expanded and relocated into target location. Catalog is metadata, and
is used to verify that created environment matches original specification.

## Controller

This is tool or context that is currently running `rcc` command.

## Diagnostics

Network, configuration, or robot diagnostics, that are executed to give
status of one or some of those aspects.

## Dirty environment

An environment, holotree space specially, become dirty when after it is
restored in pristine state, something adds, deletes, or modifies files or
directories inside that specific space. This can happen by preRunScripts
modifying something when robot start, robot itself doing something that
changes actual environment, or when someone intentionally tries to modify
or install something into environment manually.

When dirtyness is desired thing, like for developer purposes, use unmanaged
holotree spaces. But for normal automations and robots, it is good to start
from pristine, clean state.

## Environment

An environment here means either concrete holotree space, which contains
set of code and libraries (like python runtime environment) that are needed
for running specific robot or automation. Or it means that same environment
but as stashed away building blocks that are stored in hololib.

## Fingerprint

Fingerprint is normally a hash digest calculated from some content. Various
algorithms can be used for this, and some examples are Sha256 and siphash.

## Holotree

Is set of working areas, where concrete robots can run. Robots and processes
run inside one of these instances. These consume disk space. These are also
resetted into pristine state each time one of `rcc` run or environment related
subcommands are executes.

## Hololib

Is set of building blocks that are used to setup concreate holotree spaces.
Hololib contains both library and catalogs. Every unique content is stored
only once in library part. And catalogs refers to library fingerprints to
identify what parts they use.

## Identity

Identity is something that describes or identifies something uniquely.
For example `identity.yaml` is description that equals to `conda.yaml`.

## Platform

Platform refers to either Windows, MacOS, or Linux. And also either "amd64"
or "arm64" architectures of those.

## Prebuild environment

Prebuild environment is something that contains building blocks for full
holotree space in form of catalog + hololib parts. It is per operating system
and architecture, and can only be used in shared holotree context, where
parts are relocatable between different user accounts.

During building concrete holotree space from prebuild environment, there is
no need for internet connection. If something inside robot run needs internet,
then that is not prebuild environment concern.

## Pristine environment

Environment that is restored to match exactly original, specified state.
When environments are used and content inside is changed, then those are
dirty/corrupted environments. They can be restored back to pristine state
using `rcc` commands.

## Private holotree

This is state, where all environments are created for single user, and cannot
be shared between users. These must be build and managed privately and using
them normally requires internet access.

## Product family

This is reference to either Robocorp products or Sema4.ai products.

## Profile

Profile is set of settings that describes network and Robocorp configurations
so that cloud and Control Room can be used in robot context.

## Robot

Robot is automation or process, that will be running inside one of concrete
holotree space.

## Shared holotree

This is state, where created environment can be relocated and different users
can use same shared catalogs to quickly replicate environments with identical
specifications, but provided for each user as separate space.

## Space

Concrete created environment where processes and robot actually run. Each
holotree space is identified by three things: user, controller, and space
identifier. Each different combination of those values receives their own
separate directory. These will each separately consume diskspace.

## Unmanaged holotree space

This is holotree space, that is created by `rcc` but it is not managed by
`rcc` after it gets created. It is up to user or using tool to manage and
maintain that environment. It can get dirty, can have traditional tooling
adding dependencies there, and it can deviate from specification.

Note: unmanaged holotree spaces are not user specific, and managing access
to those spaces is left to tooling/users who use these unmanaged spaces.

## User

User account identity that is using `rcc`. Users wont share concrete holotrees
in shared holotree context.  Each user will get their own separate space.
