# Agents

## transliterate

ASCII Name Transliteration Service - Converts multicultural names to ASCII-compatible structured JSON for legacy systems.

### Project Setup

```bash
# Install dependencies
bun install

# Install Encore CLI (if not installed)
curl -L https://encore.dev/install.sh | bash

# Initialize Encore app
encore app create --example=hello-world
```

### Development Workflow

#### Feature Branch Development

```bash
# 1. Create feature branch from develop
git checkout develop
git pull origin develop
git checkout -b feature/your-feature-name

# 2. Develop locally with Encore
encore run

# 3. Test via dashboard at http://localhost:9400/

# 4. Merge to develop with squash
git checkout develop
git merge --squash feature/your-feature-name
git commit -m "feat: add your feature description"
git push origin develop
```

#### Deployment Commands

```bash
# Deploy to staging (Encore Cloud)
git checkout stage
git merge develop
git push encore stage

# Deploy to production (triggers Terraform)
git checkout main
git merge stage
git push origin main
```

#### Backend (Go/Encore)

```bash
# Run local development server
encore run

# Run backend tests
encore test ./...

# Run specific service tests
encore test ./api/services/transliterate

# Generate TypeScript client from Encore API
encore gen client ./frontend/src/services/client --lang=typescript

# Check Encore compilation
encore check

# View Encore dashboard
encore dashboard
```

#### Frontend (Hugo)

```bash
# Start Hugo development server
hugo server -D

# Build Hugo static site
hugo --minify

# Build TypeScript assets for Hugo
bun run build:ts
# or
esbuild hugo/assets/ts/main.ts --bundle --outdir=hugo/static/js
```

### Testing Conventions

#### Backend Testing (Go)

- **Unit tests**: `*_test.go` files alongside source
- **Test command**: `go test -v ./...`
- **Coverage**: `go test -cover ./...`
- **Benchmarks**: `go test -bench=. ./...`

Example:
```go
// userService_test.go
func TestGetUserByID(t *testing.T) {
    // Unit test
}

// userApi_test.go  
func TestUserAPIEndpoint(t *testing.T) {
    // Integration test
}
```

#### Frontend Testing (TypeScript)

- **Unit tests**: `fileName.spec.ts`
- **API validation tests**: `apiName.test.ts`
- **Test command**: `bunx vitest run`
- **Watch mode**: `bunx vitest watch`
- **Coverage**: `bunx vitest run --coverage`

Example:
```typescript
// userService.spec.ts - Unit test
// userApi.test.ts - API validation test
```

### Build & Deployment

#### Local Development
```bash
# Run local server with hot reload
encore run

# Test locally
encore test ./...

# Build for validation
encore build
```

#### Staging Deployment (Encore Cloud)
```bash
# Deploy to staging environment
git checkout stage
git merge develop
git push encore stage

# View staging deployment
open https://app.encore.dev
```

#### Production Deployment
```bash
# Export Terraform for production infrastructure
encore terraform generate --env=production ./terraform

# Deploy via GitHub Actions (automatic on main push)
git checkout main
git merge stage
git push origin main

# Manual Terraform deployment (if needed)
cd terraform
terraform init
terraform plan
terraform apply
```

### Project Structure

```
/
├── api/                    # Encore backend
│   ├── services/          # Individual services
│   │   └── user/         
│   │       ├── user.go    # Service implementation
│   │       └── user_test.go
│   ├── lib/              # Shared libraries
│   └── modules/          # Business logic modules
└── frontend/             # Hugo static site
    ├── assets/
    │   └── ts/          # TypeScript files
    ├── content/         # Markdown content
    ├── layouts/         # HTML templates
    └── static/          # Static assets
```

### API Development

#### Encore Service Pattern (Go)

```go
//encore:api public method=GET path=/transliterate
func Transliterate(ctx context.Context, input string) (*TransliterateResponse, error) {
    // Implementation
}
```

#### Encore Service Pattern (TypeScript)

```typescript
import { api } from "encore.dev/api";

export const transliterate = api<{input: string}, TransliterateResponse>(
    { method: "GET", path: "/transliterate" },
    async ({ input }) => {
        // Implementation
    }
);
```

### Database Operations

```bash
# Run migrations (Encore handles automatically)
# Place migration files in: api/services/[service]/migrations/

# Access local database console
encore db shell transliterate --env=local --superuser

# Access staging database console
encore db shell transliterate --env=staging

# View database in dashboard
open http://localhost:9400/
```

### Environment & Secrets

```bash
# Set secrets for staging
encore secret set --env=staging API_KEY

# Set secrets for production
encore secret set --env=production API_KEY
encore secret set --env=production DB_PASSWORD

# List secrets
encore secret list --env=staging
encore secret list --env=production
```

### Code Quality

```bash
# Lint Go code
golangci-lint run

# Format Go code
go fmt ./...

# Lint TypeScript
bunx eslint ./src/frontend

# Type check TypeScript
bunx tsc --noEmit
```

### Debugging

```bash
# View local logs
encore logs --env=local

# View staging logs
encore logs --env=staging

# View specific service logs
encore logs --env=local --service=transliterate

# Debug with Encore dashboard
open http://localhost:9400/  # Local
open https://app.encore.dev  # Cloud environments
```

### Common Tasks

#### Add a new service
```bash
# Create service directory
mkdir -p api/services/newservice

# Create service file and migrations
touch api/services/newservice/newservice.go
mkdir -p api/services/newservice/migrations

# Add Encore annotations and implement
# Test locally with: encore run
```

#### Add a new Hugo page
```bash
# Create content file
hugo new content/page-name.md

# Create corresponding layout if needed
touch layouts/page-name.html
```

#### Generate API documentation
```bash
# Encore automatically generates OpenAPI spec
encore api docs
```

### Troubleshooting

- **Port conflicts**: Encore uses 4000 (API) and 9400 (dashboard)
- **Build failures**: Check `encore check` output
- **Test failures**: Ensure database migrations are current
- **Client generation**: Run after any API changes

### Performance Profiling

```bash
# Go profiling
go test -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof

# Frontend bundle analysis
bunx vite-bundle-visualizer
```