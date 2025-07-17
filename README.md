# Quick Cell Reload

A Go tool for updating the "reload-hash" meta tag in Nomad jobs to trigger job reloads.

## Quick Cell Reload Tool

The `quick-cell-reload` tool updates the "reload-hash" meta tag in Nomad jobs with a new random hash value, which triggers Nomad to reload the job. This is useful for forcing job restarts without changing the actual job specification.

### Features

- Update "reload-hash" meta tag for individual jobs
- Update "reload-hash" meta tag for multiple jobs using patterns
- List current meta tags for jobs
- Support for different Nomad namespaces
- Automatically generates a new random hash for each run

### Usage

#### Build the tool

```bash
cd cmd/quick-cell-reload
go build -o quick-cell-reload main.go
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

#### List current meta tags for a job

```bash
./quick-cell-reload -job app-100117 -list
```

#### Use different namespace

```bash
./quick-cell-reload -job app-100117 -namespace platform
```

### Command Line Options

- `-job`: Job ID to update (optional, defaults to all "app-" jobs)
- `-namespace`: Nomad namespace (default: "sites")
- `-list`: List current meta tags for the job
- `-pattern`: Job pattern for updating multiple jobs (default: "app-")
- `-timeout`: Timeout for operations (default: 30s)

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

#### Example 4: List all meta tags for debugging

```bash
./quick-cell-reload -job app-100117 -list
```

### How It Works

1. **Retrieves current job spec**: Uses Nomad API to get the complete job specification
2. **Generates new hash**: Creates a random 32-character hex string
3. **Updates meta tag**: Sets the "reload-hash" meta tag to the new hash value
4. **Submits update**: Registers the modified job spec with Nomad
5. **Triggers reload**: Nomad detects the meta tag change and reloads the job

### Environment Variables

The tool uses the Nomad API client which automatically detects:
- `NOMAD_ADDR`: Nomad server address
- `NOMAD_TOKEN`: Nomad authentication token
- `NOMAD_CACERT`: Path to CA certificate
- `NOMAD_CLIENT_CERT`: Path to client certificate
- `NOMAD_CLIENT_KEY`: Path to client key

### Notes

- The tool preserves all existing meta tags and only updates the "reload-hash" tag
- Each run generates a unique hash, ensuring the job will reload even if run multiple times
- Job updates are submitted as new job registrations, which trigger Nomad evaluations
- The tool supports both single job updates and bulk updates using patterns
- When no job is specified, it defaults to updating all jobs starting with "app-" in the "sites" namespace 