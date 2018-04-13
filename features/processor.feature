Feature: Vortex processor
  In order to be able to deploy to different environments
  As a DevOps engineer
  I need a command that can inject the appropriate variables into manifest templates

  Scenario:
    Given a template file
    And a variable file
    When vortex is run for a template
    Then an output file should contain the interpolated variables

  Scenario:
    Given a template directory
    And a variable file
    When vortex is run for a directory
    Then an output directory should contain the output files
    And the output files should contain the interpolated variables
