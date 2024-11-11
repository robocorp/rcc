# Holotree and library maintenance

This documentation section gives you brief ideas how to maintain your
holotree/hololib setup.

## Why do maintenance?

There are number of reasons for doing maintenance, some of which are:

- running into problem with using holotree, and wanting to start over
- something breaks in holotree/lib and there is need to fix it
- running out of disk space, and wanting to reduced used foot print
- remove old/unused spaces from holotree
- remove old/unused catalogs from hololib
- just to keep things running smoothly on future robot invocations

## Shared holotree and maintenance

When doing maintenance in shared holotree, you should be aware, that it might
affect other user accounts in same machine. So when you are doing system wide
maintenance in shared holotree, make sure that nothing is working on those
environments and catalogs that your maintenance targets to.

## Maintenance vs. tools using holotrees

When doing maintenance on any holotree, you should be aware, that if Robocorp
tooling (Worker, Assistant, VS Code plugins, rcc, ...) is
also at same time using same holotree/hololib, your maintenance actions might
have negative effect on those tools.

Some of those effects might be:

- wiping environment under tool using it and causing automation, debugging,
  editing, or development tooling to crash or produce unexpected results
- if catalog or space was removed, and it is needed later, then that must
  be rebuild or downloaded, and that will slow down that use-case
- removing catalogs or hololib entries that will be needed by automations
  or tooling, might cause slowness when needed next time, or if builds are
  prevented it might even deny usage of those spaces

## Maintenace and product families

Since v18 of rcc, there are two different product families present. To
explicitely maintain specific product family holotree, then either
`--robocorp` or `--sema4ai` flag should be given. Both product families
have their separate holotree libraries and spaces.

## Deleting catalogs and spaces

Before you delete anything, you should be aware of those things and what is
there.

Catalogs can be listed using `rcc holotree catalogs` command, and
if you add `--identity` you can see what was their environment specification.

Then command `rcc holotree list` is used to list those concrete spaces that
are consuming your disk space. There you can also see how many times space
has been used, and when was last time it was used. (And using in this context
means, that rcc did create or refresh that specific space.)

Once you know what is there, and there are needs to remove catalogs, then
see `rcc holotree remove -h` for more about information on that. One good
option to use there is `--check 5` to also cleanup all released spare parts.

And to free disk space consumed by concrete holotrees, see command
`rcc holotree delete -h`, which can be used to delete those spaces that
are not needed anymore.

## Keeping hololib consistent

And in cases, where there are holotree restoration problems, or hololib
issues, it is good to run consistency checks against that hololib. This
can be done using `rcc holotree check -h` command. And good option there
is to add `--retries 5` option, to get more "garbage collection cycles"
to maintain used disk space.

Note that after running this command, and if there was something broken
inside hololib, then some of your catalogs have been removed, and in this
case it is good thing, since they were broken. And if they are needed in
future, those should be either build or imported.

## Summary of maintenance related commands

- `rcc holotree list -h` lists holotree spaces and their location
- `rcc holotree catalogs -h` list known catalogs, their blueprints, and stats
- `rcc configuration cleanup -h` for general cleanup procedures
- `rcc holotree delete -h` for deleting individual spaces
- `rcc holotree remove -h` for removing individual catalogs
- `rcc holotree check -h` for checking integrity of hololib
