*** Settings ***
Library             OperatingSystem
Library             supporting.py

Test Template       Verify exitcodes


*** Test Cases ***    EXITCODE    COMMAND
General failure of rcc command    1    build/rcc crapiti -h --controller citests
General output for rcc command    0    build/rcc --controller citests
Help for rcc command    0    build/rcc -h
Help for rcc assistant subcommand    0    build/rcc assistant -h --controller citests
Help for rcc cloud subcommand    0    build/rcc cloud -h --controller citests
Help for rcc community subcommand    0    build/rcc community -h --controller citests
Help for rcc configure subcommand    0    build/rcc configure -h --controller citests
Help for rcc create subcommand    0    build/rcc create -h --controller citests
Help for rcc feedback subcommand    0    build/rcc feedback -h --controller citests
Help for rcc holotree subcommand    0    build/rcc holotree -h --controller citests
Help for rcc help subcommand    0    build/rcc help -h --controller citests
Help for rcc interactive subcommand    0    build/rcc interactive -h --controller citests
Help for rcc internal subcommand    0    build/rcc internal -h --controller citests
Help for rcc man subcommand    0    build/rcc man -h --controller citests
Help for rcc pull subcommand    0    build/rcc pull -h --controller citests
Help for rcc robot subcommand    0    build/rcc robot -h --controller citests
Help for rcc run subcommand    0    build/rcc run -h --controller citests
Help for rcc task subcommand    0    build/rcc task -h --controller citests
Help for rcc tutorial subcommand    0    build/rcc tutorial -h --controller citests
Help for rcc version subcommand    0    build/rcc version -h --controller citests
Run rcc config settings    0    build/rcc config settings --controller citests
Run rcc docs changelog    0    build/rcc docs changelog --controller citests
Run rcc docs license    0    build/rcc docs license --controller citests
Run rcc docs recipes    0    build/rcc docs recipes --controller citests
Run rcc docs tutorial    0    build/rcc docs tutorial --controller citests
Run rcc holotree list    0    build/rcc holotree list --controller citests
Run rcc tutorial    0    build/rcc tutorial --controller citests
Run rcc version    0    build/rcc version --controller citests
Run rcc --version    0    build/rcc --version --controller citests


*** Keywords ***
Verify exitcodes
    [Arguments]    ${exitcode}    ${command}
    ${code}    ${output}    ${error}=    Run and return code output error    ${command}
    Log    <b>STDOUT</b><pre>${output}</pre>    html=yes
    Log    <b>STDERR</b><pre>${error}</pre>    html=yes
    Should be equal as strings    ${exitcode}    ${code}
