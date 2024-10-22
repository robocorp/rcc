*** Settings ***
Library             OperatingSystem
Library             supporting.py
Resource            resources.robot

Suite Setup         Sema4.ai setup
Suite Teardown      Sema4.ai teardown

Default Tags        wip


*** Test Cases ***
Goal: See rcc toplevel help for Sema4.ai
    Step    build/rcc --sema4ai --controller citests --help
    Must Have    SEMA4AI
    Must Have    Robocorp
    Must Have    --robocorp
    Must Have    --sema4ai
    Must Have    completion
    Must Have    robot
    Wont Have    ROBOCORP
    Wont Have    Robot
    Wont Have    assistant
    Wont Have    interactive
    Wont Have    community
    Wont Have    tutorial
    Wont Have    bash
    Wont Have    fish

Goal: See rcc commands for Sema4.ai
    Step    build/rcc --sema4ai --controller citests
    Use STDERR
    Must Have    SEMA4AI
    Wont Have    ROBOCORP
    Wont Have    Robocorp
    Wont Have    Robot
    Must Have    robot
    Wont Have    assistant
    Wont Have    interactive
    Wont Have    community
    Wont Have    tutorial
    Wont Have    completion
    Wont Have    bash
    Wont Have    fish

Goal: Default settings.yaml for Sema4.ai
    Step    build/rcc --sema4ai configuration settings --defaults --controller citests
    Must Have    Sema4.ai default settings.yaml
    Wont Have    assistant
    Wont Have    branding
    Wont Have    logo

Goal: Check holotree hash
    Step    build/rcc holotree hash --silent --controller citests robot_tests/bare_action/package.yaml
    Must Have    43139fade39b1952

    Step    build/rcc holotree hash --devdeps --silent --controller citests robot_tests/bare_action/package.yaml
    Must Have    ac2e488f48812dcf

Goal: Create package.yaml environment using uv
    Step    build/rcc --sema4ai ht vars --json -s sema4ai --controller citests robot_tests/bare_action/package.yaml
    Must Have    RCC_ENVIRONMENT_HASH
    Must Have    RCC_INSTALLATION_ID
    Must Have    SEMA4AI_HOME
    Wont Have    ROBOCORP_HOME
    Must Have    _4e67cd8_81359368
    Use STDERR
    Must Have    Progress: 01/15
    Must Have    Progress: 15/15
    # Must Have    Running uv install phase.
    Run With Env    python -c "import sys; print(sys.version)"    ${robot_stdout}    fail=${False}
    Run With Env    python -m pytest --version    ${robot_stdout}    fail=${True}

Goal: Create devenv package.yaml environment using uv
    Step
    ...    build/rcc --sema4ai ht vars --devdeps --json -s sema4ai --controller citests robot_tests/bare_action/package.yaml
    Must Have    RCC_ENVIRONMENT_HASH
    Must Have    RCC_INSTALLATION_ID
    Must Have    SEMA4AI_HOME
    Wont Have    ROBOCORP_HOME
    Must Have    _4e67cd8_81359368
    Use STDERR
    Must Have    Progress: 01/15
    Must Have    Progress: 15/15
    Must Have    Running uv install phase.
    Run With Env    python -m pytest --version    ${robot_stdout}    fail=${False}

Goal: Create venv with devdeps package.yaml environment using uv
    Remove Directory    tmp/venv-test    True
    Create Directory    tmp/venv-test

    ${cwd} =    Evaluate    os.path.abspath('tmp/venv-test')    modules=os
    ${package_yaml} =    Evaluate    os.path.abspath('robot_tests/bare_action/package.yaml')    modules=os
    IF    sys.platform == 'win32'
        ${RCC} =    Evaluate    os.path.abspath('build/rcc.exe')    modules=os
    ELSE
        ${RCC} =    Evaluate    os.path.abspath('build/rcc')    modules=os
    END

    Step
    ...    ${RCC} --sema4ai ht venv --devdeps --controller citests ${package_yaml}
    ...    cwd=${cwd}
    IF    sys.platform == 'win32'
        Run and return code output error    ${cwd}/venv/Scripts/python.exe -m pytest --version    check=${True}
    ELSE
        Run and return code output error    ${cwd}/venv/bin/python3 -m pytest --version    check=${True}
    END

    Remove Directory    tmp/venv-test    True


*** Keywords ***
Sema4.ai setup
    Remove Directory    tmp/sema4home    True
    Prepare Sema4.ai Home    tmp/sema4home

Sema4.ai teardown
    Remove Directory    tmp/sema4home    True
    Prepare Robocorp Home    tmp/robocorp
