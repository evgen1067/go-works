Feature: update event
  In order to be happy
  As a hungry gopher
  I need to be able to update an event

  Scenario: successful event update
    When I send "PUT" request to "http://localhost:8888/events/999" with "application/json" data:
    """
    {
		"title":       "Title",
		"description": "Description",
		"dateStart":   "2025-01-16T20:00:00Z",
		"dateEnd":     "2025-01-17T20:00:00Z",
		"notifyIn":    1,
		"ownerId":     1
    }
    """
    Then The response code should be 200
  Scenario: not successful event update (the date of the event is already occupied)
    When I send "PUT" request to "http://localhost:8888/events/1000" with "application/json" data:
    """
    {
		"title":       "Title",
		"description": "Description",
		"dateStart":   "2025-01-16T20:00:00Z",
		"dateEnd":     "2025-01-17T20:00:00Z",
		"notifyIn":    1,
		"ownerId":     1
    }
    """
    Then The response code should be 500
  Scenario: successful event update (the date of the event is already occupied, but the recipient is specified differently)
    When I send "PUT" request to "http://localhost:8888/events/1000" with "application/json" data:
    """
    {
		"title":       "Title",
		"description": "Description",
		"dateStart":   "2025-01-16T20:00:00Z",
		"dateEnd":     "2025-01-17T20:00:00Z",
		"notifyIn":    1,
		"ownerId":     2
    }
    """
    Then The response code should be 200