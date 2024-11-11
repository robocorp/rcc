# Support for virtual environments

There is now experimental feature in rcc to create virtual environments on top
of rcc holotree environment.

## What does it do?

When you run command `rcc venv`, that will do following things:
- using given `conda.yaml` file, it will create that base environment
- this new environment will be "unmanaged" holotree space, so once created, rcc
  wont touch base environment (unless forced)
- it will also be "externally managed" PEP-668 environment, so no additional
  things should be pip installed in that environment
- then on top of that base environment, rcc will try automatically create local
  project specific `venv` directory and list available activation commands
- in addition to that, rcc also puts experimental `depxtraction.py` script in
  same directory (more on that later in this document)

## How to get started?

- first you need rcc version 17.17.0 or later available in your system
- then on some directory, you should have `conda.yaml` for base environment,
  something like this:

```
channels:
- conda-forge
dependencies:
- python=3.10.12
- pip=23.2.1
- robocorp-truststore=0.8.0
```

- then in that directory, run command `rcc venv conda.yaml`
- after that, you should see list of activation commands to use this new venv
- after activation, you can use normal pip commands to populate that venv as
  you wish

## Limitations of `rcc venv`:

- currently naming and location is fixed, so you cannot change those
- this venv is always build on top of holotree space, so that holotree space
  must always be there
- and that space is "unmanaged", so idea is, that once created, it is developer
  responsibility to delete or force update it if dependencies change
- also other things installed from conda-forget from underlying holotree space
  are hidden and only python environment is visible

## Dangers of using `--force` in `rcc venv` context.

- unmanged holotrees are not user specific, so be careful when using `--force`
  option to recreate those spaces, and recommendation is to use `--space` and
  `--controller` options to limit usage to your intentions
- be aware that `--force` makes three things to happen
- first it is needed if `venv` was already created (so rcc wont overwrite
  things in venv, unless you really force it)
- second it is used to tell rcc, that also underlying holotree space should
  be recreated (maybe with conflicting dependencies)
- third, it forces also full holotree space installation and updating caches

## What is this `depxtraction.py` thing?

It is "dependency extraction", with limitations (see below).

Idea behing `depxtraction.py` is, that when there is modified environment,
where additional dependencies are installed using tools like pip or poetry,
those dependencies can be extracted by tooling into simple `conda.yaml`
format.

## Limitations of `depxtraction.py`:

- no conda dependencies detected, and every dependency that python tooling
  reports are expected to be from PyPI (except "hardcoded" python, pip and
  robocorp-truststore that are defined as bootstrapping dependencies from
  conda-forge)
- if there are deeply recursive dependencies (X depends on Y depends on Z
  depends on X) then currently those dependencies will vanish, since "root"
  dependency is unclear (if you run these cases, please report those, so
  that better functionality can be implemented and tested)
- only top level dependencies are resolved and versioned, and listed as
  dependencies; subdependency resolving is left for pip resolver to figure out
- and because individual install commands can create inconsistent environment,
  it is possible, that once `conda.yaml` is generated out of such environment,
  recreation of such environment might actually fail to resolve correctly and
  in those cases, you have to adjust generated `conda.yaml` accordingly

## Ideas for usage

- start VS Code from CLI inside activated environment
- create rcc venv, install packages manually inside that environment,
  make your automation work, once automation is working so far, extract
  dependencies, recreate rcc venv and continue iterating ...
- try to run `depxtraction.py` on your system python setup and see what
  comes up there (this can be done by just using `depxtraction.py` with
  your system python without activating any virtual environments)
