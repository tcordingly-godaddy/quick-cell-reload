# Quick Cell Reload

A Go tool for updating the "reload-hash" meta tag in Nomad jobs to trigger job reloads with rate limiting.

## Quick Cell Reload Tool

The `quick-cell-reload` tool updates the "reload-hash" meta tag in Nomad jobs with a new random hash value, which triggers Nomad to reload the job. This is useful for forcing job restarts without changing the actual job specification. The tool includes rate limiting to prevent overwhelming the Nomad cluster.

### Features

- Update "reload-hash" meta tag for individual jobs
- Update "reload-hash" meta tag for multiple jobs using patterns
- **Rate limiting**: Control the rate of job updates using burst, limit, and interval parameters
- Support for different Nomad namespaces
- Automatically generates a new random hash for each run
- Graceful shutdown handling with SIGTERM

### Usage

#### Build the tool

```bash
cd quick-cell-reload
go build -o quick-cell-reload ./cmd/quick-cell-reload
```

#### Update a single job with new reload hash

```bash
./quick-cell-reload -job app-100117
```

#### Update multiple jobs using pattern

```bash
./quick-cell-reload -pattern app-
```

#### Update all app jobs (default behavior)

```bash
./quick-cell-reload
```

#### Update jobs with rate limiting

```bash
# Update jobs with rate limiting: 5 requests per second with burst of 20
./quick-cell-reload -burst 20 -limit 5 -interval 1s

# Update jobs with rate limiting: 1 request per 2 seconds with burst of 10
./quick-cell-reload -burst 10 -limit 1 -interval 2s

# Conservative rate limiting: 1 request per 5 seconds with burst of 5
./quick-cell-reload -burst 5 -limit 1 -interval 5s
```

#### Use different namespace

```bash
./quick-cell-reload -job app-100117 -namespace platform
```

### Command Line Options

- `-job`: Job ID to update (optional, defaults to all "app-" jobs)
- `-namespace`: Nomad namespace (default: "sites")
- `-pattern`: Job pattern for updating multiple jobs (default: "app-")
- `-burst`: Number of requests allowed in burst (default: 10)
- `-limit`: Number of requests allowed per interval (default: 1)
- `-interval`: Time interval for rate limiting (default: 1s)

### Examples

#### Example 1: Reload a specific app

```bash
./quick-cell-reload -job app-100117
```

This will:

1. Fetch the current job specification from Nomad
2. Generate a new random 32-character hex hash
3. Update the "reload-hash" meta tag with the new hash
4. Submit the job update to Nomad, triggering a reload

#### Example 2: Reload multiple apps with a pattern

```bash
./quick-cell-reload -pattern app-
```

This will update all jobs that start with "app-" in the namespace.

#### Example 3: Reload all app jobs (default)

```bash
./quick-cell-reload
```

This will update all jobs that start with "app-" in the "sites" namespace.

#### Example 4: Rate-limited updates

```bash
./quick-cell-reload -burst 20 -limit 5 -interval 1s
```

This will:

1. Find all jobs matching the pattern (default: "app-")
2. Use rate limiting to process 5 requests per second with a burst capacity of 20
3. Continue until all matching jobs are updated

#### Example 5: Conservative rate limiting

```bash
./quick-cell-reload -burst 5 -limit 1 -interval 5s
```

### Environment Variables

The tool currently depends on port forwarding to connect to the test V2 environment and uses the Nomad API client which automatically detects:

- `NOMAD_ADDR`: Nomad server address local port forwarding address
- `NOMAD_SKIP_VERIFY`: Nomad to skip the TLS certificate verification
