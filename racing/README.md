## Race Service

Sport service is a tennis simulation where two players compete and support updates to the winner.


### ListRaces


`POST:/v1/list-races`

This endpoint returns a list of races, and then provides filters and sorting.

It accepts a JSON request body with optional `filter` and `order_by` properties as shown in the example below.

```
{
    "filter": {
        "visible": true
    },
    "order_by": {
        "parameter": "advertised_start_time",
        "direction": "ASC"
    }
}
```


`/v1/race/{id}`

his endpoint allows users to search a single race with the ID specified in the endpoint URL. 