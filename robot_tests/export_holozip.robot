*** Settings ***
Library             OperatingSystem
Library             supporting.py
Resource            resources.robot
Suite Setup         Export setup
Suite Teardown      Export teardown

*** Keywords ***
Export setup
  Remove Directory  tmp/developer  True
  Remove Directory  tmp/guest  True
  Remove Directory  tmp/standalone  True
  Prepare Robocorp Home    tmp/developer
  Fire And Forget   build/rcc ht delete 4e67cd8

Export teardown
  Prepare Robocorp Home    tmp/robocorp
  Remove Directory  tmp/developer  True
  Remove Directory  tmp/guest  True
  Remove Directory  tmp/standalone  True

*** Test cases ***

Goal: Create extended robot into tmp/standalone folder using force.
    Step    build/rcc robot init --controller citests -t extended -d tmp/standalone -f
    Use STDERR
    Must Have    OK.

    ${output}=    Capture Flat Output    build/rcc ht hash --silent tmp/standalone/conda.yaml
    Set Suite Variable    ${fingerprint}    ${output}

Goal: Create environment for standalone robot
    Step    build/rcc ht vars -s author --controller citests -r tmp/standalone/robot.yaml
    Must Have    RCC_ENVIRONMENT_HASH=
    Must Have    RCC_INSTALLATION_ID=
    Must Have    4e67cd8_fcb4b859
    Use STDERR
    Must Have    Progress: 01/15
    Must Have    Progress: 15/15

Goal: Must have author space visible
    Step    build/rcc ht ls
    Use STDERR
    Must Have    4e67cd8_fcb4b859
    Must Have    rcc.citests
    Must Have    author
    Must Have    ${fingerprint}
    Wont Have    guest

Goal: Show exportable environment list
    Step    build/rcc ht export
    Use STDERR
    Must Have    Selectable catalogs
    Must Have    - ${fingerprint}
    Must Have    OK.

Goal: Export environment for standalone robot
    Step    build/rcc ht export -z tmp/standalone/hololib.zip ${fingerprint}
    Use STDERR
    Wont Have    Selectable catalogs
    Must Have    OK.

Goal: Wrap the robot
    Step    build/rcc robot wrap -z tmp/full.zip -d tmp/standalone/
    Use STDERR
    Must Have    OK.

Goal: See contents of that robot
    Step    unzip -v tmp/full.zip
    Must Have    robot.yaml
    Must Have    conda.yaml
    Must Have    hololib.zip

Goal: Can delete author space
    Step    build/rcc ht delete 4e67cd8_fcb4b859
    Step    build/rcc ht ls
    Use STDERR
    Wont Have    4e67cd8_fcb4b859
    Wont Have    rcc.citests
    Wont Have    author
    Wont Have    ${fingerprint}
    Wont Have    guest

Goal: Can run as guest
    Fire And Forget    build/rcc ht delete 4e67cd8
    Prepare Robocorp Home    tmp/guest
  Step        build/rcc task run --controller citests -s guest -r tmp/standalone/robot.yaml -t "run example task"
    Use STDERR
    Must Have    point of view, "actual main robot run" was SUCCESS.
    Must Have    OK.

Goal: Space created under author for guest
    Prepare Robocorp Home    tmp/developer
    Step    build/rcc ht ls
    Use STDERR
    Wont Have    4e67cd8_fcb4b859
    Wont Have    author
    Must Have    rcc.citests
    Must Have    ${fingerprint}
    Must Have    4e67cd8_aacf1552
    Must Have    guest
