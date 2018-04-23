Feature: get capacity
  In order to balance exchange with capacity
  As an administrator
  I need to get the capacity of exchanges

  Scenario Outline: Get capacity
    Given exchange exist "<exchangeExist>"
    And has "<deviceTypeOfFirstDevice>" "<portsOfFirstDevice>" ports
    And has "<deviceTypeOfSecondDevice>" "<portsOfSecondDevice>" ports
    When get capacity
    Then return error "<hasError>"
    And has "Fibre" capacity "<hasFibreCapacity>"
    And has "Adsl" capacity "<hasAdslCapacity>"

    Examples:
    | exchangeExist | deviceTypeOfFirstDevice | portsOfFirstDevice | deviceTypeOfSecondDevice | portsOfSecondDevice | hasError | hasFibreCapacity | hasAdslCapacity |
    |     false     |  Fibre                  |        0           | Adsl                     |           0         |   true   |      false       |      false      |
    |     true      |  Adsl                   |        1           | Adsl                     |           3         |   false  |      false       |      false      |
    |     true      |  Fibre                  |        1           | Fibre                    |           3         |   false  |      false       |      false      |
    |     true      |  Adsl                   |        2           | Adsl                     |           3         |   false  |      false       |      true       |
    |     true      |  Fibre                  |        2           | Fibre                    |           3         |   false  |      true        |      false      |
    |     true      |  Fibre                  |        0           | Adsl                     |           0         |   false  |      false       |      false      |
