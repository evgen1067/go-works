Feature: create event
  In order to be happy
  As a hungry gopher
  I need to be able to add an event

  Scenario: successful event creation
    When I send "POST" request to "http://localhost:8888/events/new" with "application/json" data:
    """
    {
		"title":       "Title",
		"description": "Description",
		"dateStart":   "2024-01-16T20:00:00Z",
		"dateEnd":     "2024-01-17T20:00:00Z",
		"notifyIn":    1,
		"ownerId":     0
    }
    """
    Then The response code should be 201
  Scenario: not successful event creation (the date of the event is already occupied)
    When I send "POST" request to "http://localhost:8888/events/new" with "application/json" data:
    """
    {
		"title":       "Title",
		"description": "Description",
		"dateStart":   "2024-01-16T20:00:00Z",
		"dateEnd":     "2024-01-17T20:00:00Z",
		"notifyIn":    1,
		"ownerId":     0
    }
    """
    Then The response code should be 500
  Scenario: successful event creation (the date of the event is already occupied, but the recipient is specified differently)
    When I send "POST" request to "http://localhost:8888/events/new" with "application/json" data:
    """
    {
		"title":       "Title",
		"description": "Description",
		"dateStart":   "2024-01-16T20:00:00Z",
		"dateEnd":     "2024-01-17T20:00:00Z",
		"notifyIn":    1,
		"ownerId":     1
    }
    """
    Then The response code should be 201