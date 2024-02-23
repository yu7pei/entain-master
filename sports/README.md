## Sports Service

Sport service is a tennis simulation where two players compete and support updates to the winner.


### ListEvents


`POST:/v1/list-events`

This endpoint returns a list of tennis events, and then provides filters and sorting.

It accepts a JSON request body with optional `filter` and `order` properties as shown in the example below.

```
{
    "filter": {
        "visible": true,
        "player": "Bria Purdy",
        "arena": "Louisiana",
        "winner": "Bria Purdy"
    },
    "order": {
        "parameter": "advertised_start_time",
        "direction": "ASC"
    }
}
```

`/v1/update-winner`

This endpoint allows user update a winner to an event.

```
{
  "id": 69,
  "winner": "Griffin Borer"
}
```

`/v1/single-event/{id}`

his endpoint allows users to search a single event with the ID specified in the endpoint URL. 