*** Settings ***
Documentation   Template robot main suite.
Library         Collections

Library         MyLibrary
Resource        keywords.robot
Variables       MyVariables.py


*** Tasks ***
Example Task
    Example Keyword
    Example Python Keyword
    Log     ${TODAY}


