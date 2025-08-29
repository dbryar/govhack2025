# Project Learnings

## Encore.dev Project Structure

**Date**: 2025-01-29  
**Issue**: Incorrect assumption about Encore project structure  
**Learning**: 
- Encore.dev projects can be organised with `/api` directory containing the backend services
- The `encore.app` configuration file should be placed at the root of the Encore application (in `/api` directory)
- Services should remain in `/api/services/` as per project architecture documentation
- Don't move code to match tooling expectations - configure tooling to work with intended structure
- Need to run Encore commands from the directory containing `encore.app` (the `/api` directory)

**Resolution**: 
- Keep `/api` directory structure as designed in architecture.md
- Run Encore commands from `/api` directory where `encore.app` is located
- Update commands in documentation to reflect correct working directory

## Testing with Encore

**Date**: 2025-01-29  
**Issue**: Tests must be run from correct directory  
**Learning**: 
- `encore test` command must be run from directory containing `encore.app`
- For this project structure, that means running from `/api` directory
- Root-level `go.mod` exists but Encore services are in `/api` subdirectory

**Resolution**: Use `cd api && encore test ./...` for testing