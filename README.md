# Moon FastCGI

Moon FastCGI is a small Go service that returns lunar phase information for a
requested date. It can render either JSON or a minimal HTML page.

The project includes:

- A FastCGI HTTP entry point in `moon.go`
- Moon phase, illumination, distance, and phase-date calculations in
  `moon_phase.go`
- A minimal HTML template in `phase_tenplates.go`
- Phase reference images in `images/`

## Requirements

- Go 1.17 or newer
- A FastCGI-capable web server or process manager

## Install

Download dependencies with:

```sh
go mod download
```

Build the FastCGI binary with:

```sh
go build
```

Run the binary in the environment where your FastCGI server will connect to it.
The process writes logs to `info.log` in the working directory.

## Endpoints

### JSON

```text
/moon.fcgi
/moon.fcgi/moon/{date}
```

When no date is provided, the service uses the current time.

The date path parameter is parsed with
[`dateparse.ParseAny`](https://pkg.go.dev/github.com/araddon/dateparse), so
inputs such as `2026-06-18` or `June 18, 2026` are accepted.

Example response fields include:

- `date`
- `phase`
- `Illumination`
- `Age`
- `Distance`
- `Diameter`
- `SunDistance`
- `SunDiameter`
- `NextNewMoon`
- `NextFullMoon`
- `ZodiacSign`

### HTML

Add an `html` query parameter to render a minimal HTML page:

```text
/moon.fcgi/moon/2026-06-18?html=1
```

## Development

Run formatting before committing changes:

```sh
gofmt -w *.go
```

Run the test command with:

```sh
go test ./...
```

There are no automated tests in the repository yet.

## License

This project is licensed under the MIT License. See `LICENSE` for details.
