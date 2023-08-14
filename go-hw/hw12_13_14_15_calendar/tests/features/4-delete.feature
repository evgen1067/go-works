Feature: delete event
  In order to be happy
  As a hungry gopher
  I need to be able to delete an event
  Scenario: successful event delete
    When I send "DELETE" request to "http://localhost:8888/events/999"
    Then The response code should be 202
  Scenario: not successful event delete (the date of the event is already occupied)
    When I send "DELETE" request to "http://localhost:8888/events/999"
    Then The response code should be 500