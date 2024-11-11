# Incomplete list of rcc features

* supported operating systems are Windows, MacOS, and Linux
* support for product families `--robocorp` and `--sema4ai`
* supported sources for environment building are both conda and pypi
* provide repeatable, isolated, and clean environments for automations and
  robots to run on
* automatic environment creation based on declarative conda environment.yaml
  files
* easily run software robots (automations) based on declarative robot.yaml files
* also support environment creation from package.yaml files
* test robots in isolated environments before uploading them to Control Room
* provide commands for Robocorp runtime and developer tools (Worker, Assistant,
  VS Code, ...)
* provides commands to communicate with Robocorp Control Room from command line
* enable caching dormant environments in efficiently and activating them locally
  when required without need to reinstall anything
* diagnose robots and network settings to see if there is something that prevents
  using robots in specific environment
* support multiple configuration profiles for different network locations and
  conditions (remote, office, restricted networks, ...)
* running assistants from command line
* support prebuild environments, where that environment was build elsewhere
  and then just imported for local consumption
* allow "mass" prebuilding environments for delivery to those environments
  where it is not desired to build those locally
* support unmanaged environments, where rcc only initially build environment
  by the spec, but after that, does not do additional management of it
