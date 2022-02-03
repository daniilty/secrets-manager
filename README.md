## Secrets manager
Inner service for managing app configs. `SOLID`, `12 factor`, bla-bla
## Routes
* GET /secrets?app_name=sample - get secrets for application
* POST /secrets - upload configuration for some app `{"app_name": "sample", "secret": "{\"k\": \"v\"}"}`

