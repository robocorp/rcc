*** Settings ***
Library  OperatingSystem
Library  supporting.py
Resource  resources.robot

*** Test cases ***

Using and running template example with shell file

  Goal        Show rcc version information.
  Step        build/rcc version --controller citests
  Must Have   v6.

  Goal        Show rcc license information.
  Step        build/rcc man license --controller citests
  Must Have   Apache License
  Must Have   Version 2.0
  Must Have   http://www.apache.org/licenses/LICENSE-2.0
  Must Have   Copyright 2020 Robocorp Technologies, Inc.
  Wont Have   EULA

  Goal        Telemetry tracking enabled by default.
  Step        build/rcc configure identity --controller citests
  Must Have   anonymous health tracking is: enabled
  Must Exist  %{ROBOCORP_HOME}/rcc.yaml
  Wont Exist  %{ROBOCORP_HOME}/rcccache.yaml

  Goal        Send telemetry data to cloud.
  Step        build/rcc feedback metric --controller citests -t test -n rcc.test -v robot.fullrun
  Must Have   OK

  Goal        Telemetry tracking can be disabled.
  Step        build/rcc configure identity --controller citests --do-not-track
  Must Have   anonymous health tracking is: disabled

  Goal        Show listing of rcc commands.
  Step        build/rcc --controller citests
  Must Have   rcc is environment manager
  Wont Have   missing

  Goal        Show toplevel help for rcc.
  Step        build/rcc -h
  Must Have   Available Commands:

  Goal        Show config help for rcc.
  Step        build/rcc config -h --controller citests
  Must Have   Available Commands:
  Must Have   credentials

  Goal        List available robot templates.
  Step        build/rcc robot init -l --controller citests
  Must Have   extended
  Must Have   python
  Must Have   standard
  Must Have   OK.

  Goal        Initialize new standard robot into tmp/fluffy folder using force.
  Step        build/rcc robot init --controller citests -t extended -d tmp/fluffy -f
  Must Have   OK.

  Goal        There should now be fluffy in robot listing
  Step        build/rcc robot list --controller citests -j
  Must Be Json Response
  Must Have   fluffy
  Must Have   "robot"

  Goal        Fail to initialize new standard robot into tmp/fluffy without force.
  Step        build/rcc robot init --controller citests -t extended -d tmp/fluffy  2
  Must Have   Error: Directory
  Must Have   fluffy is not empty

  Goal        Run task in place.
  Step        build/rcc task run --controller citests -r tmp/fluffy/robot.yaml
  Must Have   Progress: 0/4
  Must Have   Progress: 4/4
  Must Have   rpaframework
  Must Have   1 critical task, 1 passed, 0 failed
  Must Have   OK.
  Must Exist  %{ROBOCORP_HOME}/base/
  Must Exist  %{ROBOCORP_HOME}/live/
  Must Exist  %{ROBOCORP_HOME}/wheels/
  Must Exist  %{ROBOCORP_HOME}/pipcache/

  Goal        Run task in clean temporary directory.
  Step        build/rcc task testrun --controller citests -r tmp/fluffy/robot.yaml
  Must Have   Progress: 0/4
  Wont Have   Progress: 1/4
  Wont Have   Progress: 2/4
  Wont Have   Progress: 3/4
  Must Have   Progress: 4/4
  Must Have   rpaframework
  Must Have   1 critical task, 1 passed, 0 failed
  Must Have   OK.

  Goal        Merge two different conda.yaml files with conflict fails
  Step        build/rcc env new --controller citests conda/testdata/conda.yaml conda/testdata/other.yaml  1
  Must Have   robotframework=3.1 vs. robotframework=3.2

  Goal        Merge two different conda.yaml files with conflict fails
  Step        build/rcc env new --controller citests conda/testdata/other.yaml conda/testdata/third.yaml --silent
  Must Have   44d08c86724dd710b33228c3a1ea1a0b4bef1a44f01fe825d0e412044aaa7c0242030a

  Goal        See variables from specific environment without robot.yaml knowledge
  Step        build/rcc env variables --controller citests conda/testdata/conda.yaml
  Must Have   ROBOCORP_HOME=
  Must Have   PYTHON_EXE=
  Must Have   CONDA_DEFAULT_ENV=rcc
  Must Have   CONDA_EXE=
  Must Have   CONDA_PREFIX=
  Must Have   CONDA_PROMPT_MODIFIER=(rcc)
  Must Have   CONDA_PYTHON_EXE=
  Must Have   CONDA_SHLVL=1
  Must Have   PATH=
  Must Have   PYTHONPATH=
  Must Have   PYTHONHOME=
  Must Have   PYTHONEXECUTABLE=
  Must Have   PYTHONNOUSERSITE=1
  Wont Have   ROBOT_ROOT=
  Wont Have   ROBOT_ARTIFACTS=
  Must Have   ded08c86224cc710b22228d3a1aa1a074bdf1a44f01be819c0a816044eebb80242030a

  Goal        See variables from specific environment without robot.yaml knowledge in JSON form
  Step        build/rcc env variables --controller citests --json conda/testdata/conda.yaml
  Must Be Json Response

  Goal        See variables from specific environment with robot.yaml knowledge
  Step        build/rcc env variables --controller citests conda/testdata/conda.yaml --config tmp/alternative.yaml -r tmp/fluffy/robot.yaml -e tmp/fluffy/devdata/env.json
  Must Have   ROBOCORP_HOME=
  Must Have   PYTHON_EXE=
  Must Have   THE_ANSWER=42
  Must Have   CONDA_DEFAULT_ENV=rcc
  Must Have   CONDA_EXE=
  Must Have   CONDA_PREFIX=
  Must Have   CONDA_PROMPT_MODIFIER=(rcc)
  Must Have   CONDA_PYTHON_EXE=
  Must Have   CONDA_SHLVL=1
  Must Have   PATH=
  Must Have   PYTHONPATH=
  Must Have   PYTHONHOME=
  Must Have   PYTHONEXECUTABLE=
  Must Have   PYTHONNOUSERSITE=1
  Must Have   ROBOT_ROOT=
  Must Have   ROBOT_ARTIFACTS=
  Wont Have   RC_API_SECRET_HOST=
  Wont Have   RC_API_WORKITEM_HOST=
  Wont Have   RC_API_SECRET_TOKEN=
  Wont Have   RC_API_WORKITEM_TOKEN=
  Wont Have   RC_WORKSPACE_ID=

  Goal        See variables from specific environment with robot.yaml knowledge in JSON form
  Step        build/rcc env variables --controller citests --json conda/testdata/conda.yaml --config tmp/alternative.yaml -r tmp/fluffy/robot.yaml -e tmp/fluffy/devdata/env.json
  Must Be Json Response

