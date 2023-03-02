# go-http-echo

Just another echo HTTP server.

## Running

```bash
$ go run main.go
```

## Docker

```bash
$ docker run gkawamoto/go-http-echo
```

## Why?

I wanted a small, fast, dependency-free echo HTTP server to debug stuff, which I could minimally configure through flags or environment variables:
* Structured logs
* Response status code
* Port to listen on

None of which I could find somewhere else, so here's for whoever else might need it.