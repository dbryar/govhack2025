# Transliterate - ASCII Name Transliteration Service

ASCII Name Transliteration Service - Converts multicultural names to ASCII-compatible structured JSON for legacy systems.

## Development Workflow

This project uses a feature branch workflow with squash merging and multiple deployment environments.

### Branch Strategy

- **`main`** - Production-ready code (deploys via Terraform to public cloud)
- **`stage`** - Staging environment (deploys to Encore Cloud for testing)  
- **`develop`** - Integration branch for features
- **`feature/*`** - Individual feature development branches

### Development Process

1. **Create feature branch** from `develop`:
   ```bash
   git checkout develop
   git pull origin develop
   git checkout -b feature/your-feature-name
   ```

2. **Develop locally** using Encore:
   ```bash
   encore run
   ```

3. **Test locally** via Encore dashboard at [http://localhost:9400/](http://localhost:9400/)

4. **Merge to develop** with squash:
   ```bash
   git checkout develop
   git merge --squash feature/your-feature-name
   git commit -m "feat: add your feature description"
   git push origin develop
   ```

5. **Deploy to staging** for remote testing:
   ```bash
   git checkout stage
   git merge develop
   git push encore stage
   ```

6. **Deploy to production** when ready:
   ```bash
   git checkout main  
   git merge stage
   git push origin main  # Triggers Terraform deployment
   ```

## Prerequisites 

**Install Encore:**
- **macOS:** `brew install encoredev/tap/encore`
- **Linux:** `curl -L https://encore.dev/install.sh | bash`
- **Windows:** `iwr https://encore.dev/install.ps1 | iex`
  
**Docker:**
1. [Install Docker](https://docker.com)
2. Start Docker

## Local Development

Start the development server:

```bash
encore run
```

Open the Encore developer dashboard at [http://localhost:9400/](http://localhost:9400/) to:
- View API traces and logs
- See architecture diagram  
- Browse API documentation
- Monitor database queries

## API Usage

### POST /transliterate — Transliterate text between scripts

```bash
curl 'http://localhost:4000/transliterate' \
  -H 'Content-Type: application/json' \
  -d '{
    "text": "Привет",
    "input_script": "cyrillic",
    "output_script": "latin",
    "input_locale": "ru-RU"
  }'
```

### GET /transliterate/:id — Retrieve stored transliteration

```bash
curl 'http://localhost:4000/transliterate/uuid-here'
```

### POST /transliterate/:id/feedback — Submit user feedback

```bash
curl 'http://localhost:4000/transliterate/uuid-here/feedback' \
  -H 'Content-Type: application/json' \
  -d '{
    "suggested_output": "Better transliteration",
    "feedback_type": "correction",
    "user_context": "More accurate pronunciation"
  }'
```

## Database Access

Connect to your local database:

```bash
encore db shell transliterate --env=local --superuser
```

View database schema and data through the Encore dashboard's database section.

## Testing

Run all tests:
```bash
encore test ./...
```

Run with coverage:
```bash
encore test -cover ./...
```

## Deployment Environments

### Local Development
- Use `encore run` for immediate feedback
- Database automatically provisioned
- Hot reload on code changes

### Staging (Encore Cloud)
- Deployed via `git push encore stage`
- Full cloud environment for integration testing
- Accessible via Encore Cloud dashboard

### Production (Public Cloud)
- Deployed via Terraform when pushing to `main`
- Infrastructure as code for scalability
- Monitoring and alerting configured

## Project Structure

```
/
├── api/                     # Encore backend
│   ├── encore.app          # App configuration
│   └── services/           # Service modules
│       └── transliterate/  # Transliteration service
├── frontend/               # Frontend application (future)
├── docs/                   # Documentation
├── scripts/                # Build and deployment scripts
└── terraform/              # Production infrastructure (generated)
```

## Getting Help

- **Encore Documentation**: https://encore.dev/docs
- **API Dashboard**: http://localhost:9400/ (when running locally)
- **Cloud Dashboard**: https://app.encore.dev