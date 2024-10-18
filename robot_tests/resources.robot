*** Settings ***
Library     OperatingSystem
Library     supporting.py


*** Keywords ***
Clean Local
    Remove Directory    tmp/robocorp    True

Prepare Sema4.ai Home
    [Arguments]    ${location}
    Create Directory    ${location}
    Remove Environment Variable    ROBOCORP_HOME
    Set Environment Variable    SEMA4AI_HOME    ${location}
    Copy File    robot_tests/settings.yaml    ${location}/settings.yaml
    Fire And Forget    build/rcc --sema4ai ht init --revoke --controller citests

Prepare Robocorp Home
    [Arguments]    ${location}
    Create Directory    ${location}
    Remove Environment Variable    SEMA4AI_HOME
    Set Environment Variable    ROBOCORP_HOME    ${location}
    Copy File    robot_tests/settings.yaml    ${location}/settings.yaml

Prepare Local
    Remove Directory    tmp/fluffy    True
    Remove Directory    tmp/nodogs    True
    Remove Directory    tmp/robocorp    True
    Remove File    tmp/nodogs.zip
    Prepare Robocorp Home    tmp/robocorp

    Comment    Make sure that tests do not use shared holotree
    Fire And Forget    build/rcc ht init --revoke

    Fire And Forget    build/rcc ht delete 4e67cd8

    Comment    Verify micromamba is installed or download and install it.
    Step    build/rcc ht vars --controller citests robot_tests/conda.yaml
    Must Exist    %{ROBOCORP_HOME}/bin/
    Must Exist    %{ROBOCORP_HOME}/wheels/
    Must Exist    %{ROBOCORP_HOME}/pipcache/

Fire And Forget
    [Arguments]    ${command}
    ${code}    ${output}    ${error}=    Run and return code output error    ${command}
    Log    <b>STDOUT</b><pre>${output}</pre>    html=yes
    Log    <b>STDERR</b><pre>${error}</pre>    html=yes

Step
    [Arguments]    ${command}    ${expected}=0    ${cwd}=${None}
    ${code}    ${output}    ${error}=    Run and return code output error    ${command}    cwd=${cwd}
    Set Suite Variable    ${robot_stdout}    ${output}
    Set Suite Variable    ${robot_stderr}    ${error}
    Use Stdout
    Log    <b>STDOUT</b><pre>${output}</pre>    html=yes
    Log    <b>STDERR</b><pre>${error}</pre>    html=yes
    Should be equal as strings    ${expected}    ${code}
    Wont Have    Failure:

Use Stdout
    Set Suite Variable    ${robot_output}    ${robot_stdout}

Use Stderr
    Set Suite Variable    ${robot_output}    ${robot_stderr}

Must Be
    [Arguments]    ${content}
    Should Be Equal As Strings    ${robot_output}    ${content}

Wont Be
    [Arguments]    ${content}
    Should Not Be Equal As Strings    ${robot_output}    ${content}

Must Have
    [Arguments]    ${content}
    Should Contain    ${robot_output}    ${content}

Wont Have
    [Arguments]    ${content}
    Should Not Contain    ${robot_output}    ${content}

Must Exist
    [Arguments]    ${filepath}
    Should Exist    ${filepath}

Wont Exist
    [Arguments]    ${filepath}
    Should Not Exist    ${filepath}

Must Be Json Response
    Parse JSON    ${robot_output}
