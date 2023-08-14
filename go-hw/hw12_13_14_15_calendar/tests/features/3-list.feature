Feature: list of events
  In order to be happy
  As a hungry gopher
  I need to be able to get an event list
  Scenario: not successful list (invalid event period)
    When I send "GET" request to "http://localhost:8888/events/list/fail"
    Then The response code should be 400
  Scenario: successful list
    When I send "GET" request to "http://localhost:8888/events/list/day?date=2025-01-16T19:00:00"
    Then The response code should be 200
    And I receive data:
    """
      {
      "code":200,
      "events":[
        {"id":999,"title":"Title","description":"Description","dateStart":"2025-01-16T20:00:00Z","dateEnd":"2025-01-17T20:00:00Z","notifyIn":1,"ownerId":1},
        {"id":1000,"title":"Title","description":"Description","dateStart":"2025-01-16T20:00:00Z","dateEnd":"2025-01-17T20:00:00Z","notifyIn":1,"ownerId":2}
      ]
      }
    """