# flight_tracker

`flight_tracker` is service which esposes an API that can help us understand and track how a particular person's flight path may be queried. The API accepts a request that includes a `list of flights`, which are defined by a `source airport code` and `destination airport code`. These flights can be in any order, the result will be `[source_airport_code, destination_airport_code]` with a status code of `200 OK` for all valid flights or an error message with a proper status code explaining the error if the provided flight details are not valid.


## endpoint & usage
>
> ### 
>
> - POST `/track`
>   - body
>       - [["SFO",  "EWR"]]
>   - result
>       - ["SFO", "EWR"]
>   - body
>       - [["ATL", "EWR"], ["SFO", "ATL"]]
>   - result
>       - ["SFO", "EWR"]
>   - body
>       - [["IND", "EWR"], ["SFO", "ATL"], ["GSO", "IND"], ["ATL", "GSO"]]
>   - result
>       - ["SFO", "EWR"]
>   - body
>       - [["SFO", "EWR"], ["DEL", "BLR"]]
>   - result
>       - invalid flights



## prerequisites
<b>`go` or `docker`</b>

`flight_tracker` service at bare minimim needs `go` installed or `docker` installed/running

It will serve on port `8080` by default

## execution

<b>with `go` installed</b>

>`make run` will run the `flight_tracker` service

<b>with `docker` installed</b>

>`make docker_run` will start the `flight_tracker` service

>`make docker_stop` will teardown the `flight_tracker` service

<b>common</b>

>`make test` will run test cases & test `flight_tracker` service

>`make coverage` will check the code coverage of `flight_tracker` service
