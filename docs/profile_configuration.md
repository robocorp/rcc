# Profile Configuration

## What is profile?

Profile is way to capture configuration information related to specific
network location. System can have multiple profiles defined, but only one
can be active at any moment.

### When do you need profiles?

- if you are in restricted network where direct network access is not available
- if you are working in multiple locations with different access policies
  (for example switching between office, hotel, airport, or remote locations)
- if you want to share your working setup with others in same network

### What does it contain?

- information from `settings.yaml` (can be partial)
- configuration for micromamba ("micromambarc" is almost like "condarc")
- configuration for pip (pip.ini or piprc)
- root certificate bundle in pem format
- proxy settings (`HTTP_PROXY` and `HTTPS_PROXY`)
- options for `ssl-verify` and `ssl-no-revoke`

## Quick start guide

### Setup Utility -- user interface for this

More behind [this link](https://sema4.ai/docs/automation/control-room/setup-utility).

### Pure rcc workflow for handling existing profiles

```sh
# import that Office profile, so that it can be used
rcc configuration import --filename profile_office.yaml

# start using that Office profile
rcc configuration switch --profile Office

# verify that basic things work by doing diagnostics
rcc configuration diagnostics

# when basics work, see if full environment creation works
rcc configuration speedtest

# when you want to reset profile to "system default" state
# in practice this means that all settings files removed
rcc configuration switch --noprofile

# if you want to export profile and deliver to others
rcc configuration export --profile Office --filename shared.yaml
```

## What is needed?

- you need rcc 11.9.10 or later
- your existing `settings.yaml` file (optional)
- your existing `micromambarc` file (optional)
- your existing `pip.ini` file (optional)
- your existing `cabundle.pem` file (optional)
- knowledge about your network proxies and certificate policies

## Discovery process

1. You must be inside that network that you are targetting the configuration.
2. Run Setup Utility and use it to setup and verify your profile.
3. Export profile and share it with rest of your team/organization.
4. Create other profiles for different network locations (remote, VPN, ...)
