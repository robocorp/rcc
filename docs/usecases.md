# Incomplete list of rcc use cases

* run robots in Robocorp Worker locally or in cloud containers
* run robots in Robocorp Assistant
* provide commands for Robocorp Code to develop robots locally and
  communicate to Robocorp Control Room
* provide commands that can be used in CI pipelines (Jenkins, Gitlab CI, ...)
  to push robots into Robocorp Control Room
* can also be used to run robot tests in CI/CD environments
* provide isolated environments to run python scripts and applications
* to use other scripting languages and tools available from conda-forge (or
  conda in general) with isolated and easily installed manner (see list below
  for ideas what is available)
* provide above things in computers, where internet access is restricted or
  prohibited (using pre-made hololib.zip environments, or importing prebuild
  environments build elsewhere)
* pull and run community created robots without Control Room requirement
* use rcc provided holotree environments as soft-containers (they are isolated
  environments, but also have access to rest of your machine resources)

## What is available from conda-forge?

* python and libraries
* ruby and libraries
* perl and libraries
* lua and libraries
* r and libraries
* julia and libraries
* make, cmake and compilers (C++, Fortran, ...)
* nodejs
* nginx
* rust
* php
* go
* gawk, sed, and emacs, vim
* ROS libraries (robot operating system)
* firefox
