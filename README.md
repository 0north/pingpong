# pingpong
Go ping-pong app used to test internal service communication in EKS.

This version of the app includes Datadog tracer and profiler, as well as the logger that automatically injects `trace_id` and `span_id` to correlate traces with logs.

## Run

If `PING_URL` environmental variable is set, the app runs in `ping` mode. In `ping` mode, the app will poll the `PING_URL` every 5 seconds.

Otherwise, it runs in `pong` mode. In `pong` mode, the app will an HTTP server on port 8080, and exposes a `/ping` endpoint. Pass the endpoint via `PING_URL` to another instance of `pingpong`.


```bash
go mod download

// Run `ping` mode
PING_URL=http://pong/ping go run src/*.go

// Run `pong` mode
go run src/*.go

```

## Deploy

```bash
aws ecr get-login-password --region eu-west-1 | docker login --username AWS --password-stdin 407087036459.dkr.ecr.eu-west-1.amazonaws.com
docker build -t 407087036459.dkr.ecr.eu-west-1.amazonaws.com/pingpong:datadog --platform linux/amd64 .
docker push 407087036459.dkr.ecr.eu-west-1.amazonaws.com/pingpong:datadog
```
