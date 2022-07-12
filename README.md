# dp-maps-api
Maps API for the ONSWebsite. Provides OS maps.

### Getting started

* Run `make debug`

Note: In order to use the proxy, you will need an O/S access key defined in your environment or by running.

* `ORDNANCE_SURVEY_API_KEY=abcxyz123 make debug` (replacing `abcxyz123` with your own key)

### Dependencies

* No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable         | Default              | Description                                                                                                        |
|------------------------------|----------------------|--------------------------------------------------------------------------------------------------------------------|
| BIND_ADDR                    | localhost:27900      | The host and port to bind to                                                                                       |
| GRACEFUL_SHUTDOWN_TIMEOUT    | 5s                   | The graceful shutdown timeout in seconds (`time.Duration` format)                                                  |
| HEALTHCHECK_INTERVAL         | 30s                  | Time between self-healthchecks (`time.Duration` format)                                                            |
| HEALTHCHECK_CRITICAL_TIMEOUT | 90s                  | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format) |
| ORDNANCE_SURVEY_API_URL      | "https://api.os.uk/" | The URL for the Ordnance Survey API                                                                                | 
| ORDNANCE_SURVEY_API_KEY      |                      | The access key for the Ordnance Survey API                                                                         | 

Note: To obtain an Ordnance Survey access key, you will need access to a project on the [O/S Data Hub](https://osdatahub.os.uk/projects).

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2022, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

