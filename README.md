# VALIMONT

### Prerequisites:
- [Docker](https://www.docker.com/products/docker-desktop/)
- [Make](https://www.gnu.org/software/make/)
- [Go](https://go.dev/doc/install)

### Running:
To start the application, simply run `make` from the root directory where the Makefile is located. It will execute the following steps:
- Runs the [docker-compose](/build/docker-compose.yml) command, which contains Grafana, Prometheus, and Otel services
- Builds the binary
- Runs the binary

After that, navigate to Grafana (`localhost:3000`), open dashboards -> select **valimont**, and view the metrics.

## Addresses
- Grafana: `localhost:3000`
- Prometheus UI: `localhost:9090`
- OpenAPI spec: `localhost:8080/swagger/index.html`

## General Description
Prometheus supports both push and pull-based metrics. Our application uses pull-based metrics, meaning Prometheus queries the `/metrics` endpoint and fetches all metrics. This implies that it is a web application running a web host on port `8080`. Metrics are populated by a background daemon called **listener**. The listener is responsible for calling the Validators APIs to fetch attestation data and record relevant metrics.

The application also implements graceful shutdown, which propagates the Go cancellation context through the application layers. The **validator** has a rate limiter implemented, which is set from the application configuration and is equal to 10 requests per minute. The listener has a polling interval set to 30 seconds.

### Frameworks
- Logs - [slog](https://go.dev/blog/slog)
- Metrics - [prometheus](https://prometheus.io/)
- Traces - [opentelemetry](https://opentelemetry.io/)(with prometheus integration that scrapes from otel)

### Abstractions
I decided to abstract away the **web host** and **router**, which was not strictly necessary but makes the application entry point (`main.go`) look cleaner.

## Considerations
In general, implementing the Otel standard was a "good to have," so Prometheus SDKs were used from the start. In a production environment, I would prefer to drop the application dependency on Prometheus SDKs and instead use Otel libraries to maintain a unified standard for metrics, traces, and logs.