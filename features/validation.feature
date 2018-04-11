Feature: Vortex validation
  In order ensure there are no typos or missed variables
  As a DevOps engineer
  I need a command that validates my templates against a variables file

  Scenario Outline: syntax validation
    Given <article> <template-validity> <input>
    And a <vars-validity> vars.yaml file
    When vortex is run with the -validate flag for a <input>
    Then validation should be <result>

    Examples:
      | article | template-validity | vars-validity | input          | result       |
      | a       | valid             | valid         | template       | successful   |
      | an      | invalid           | valid         | template       | unsuccessful |
      | a       | valid             | invalid       | template       | unsuccessful |
      | a       | valid             | valid         | directory      | successful   |
      | a       | mixed             | valid         | directory      | unsuccessful |
      | a       | valid             | valid         | directory tree | successful   |
      | a       | mixed             | valid         | directory tree | unsuccessful |

  Scenario Outline: variable availability validation
    Given a vars.yaml that <variables> all expected variables for a <input>
    When vortex is run with the -validate flag for a <input>
    Then validation should be <result>

    Examples:
      | variables       | input          | result       |
      | contains        | template       | successful   |
      | doesn't contain | template       | unsuccessful |
      | contains        | directory      | successful   |
      | doesn't contain | directory      | unsuccessful |
      | contains        | directory tree | successful   |
      | doesn't contain | directory tree | unsuccessful |
