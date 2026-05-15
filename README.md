# api-rest

REST API endpoint provider for MuxCore.

Registers HTTP API routes under the `/api/v1/` prefix on the core API server. Provides programmatic access to core functionality.

## Endpoints

- `GET /api/v1/modules` — List all registered modules with metadata and state
- `GET /api/v1/health` — Health check endpoint
- `GET /api/v1/status` — System status overview

## Capabilities

- `api.rest` — RESTful API endpoint provider
- `api.health` — Health check endpoint
- `api.status` — System status reporting

## Usage

```go
import "github.com/Muxcore-Media/api-rest"

mod := apirest.NewModule()
mgr.Register(mod, nil)
```
