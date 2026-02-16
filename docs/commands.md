# Commands

## Installation

```bash
pnpm install              # Install all dependencies
make dev                  # Alternative: install via make
make fresh                # Clean install (removes node_modules first)
```

## Development

```bash
# Run individual web apps (from their directories)
cd client/web/one && pnpm dev
cd client/web/compose && pnpm dev
cd client/web/taq && pnpm dev

# Or use make
cd client/web/one && make dev
```

## Testing

```bash
make test                                        # Test all (libs, clients, server)
pnpm --filter corteza-webapp-one test:unit       # Test specific app
pnpm --filter corteza-webapp-compose test:unit   # Test compose
pnpm --filter corteza-webapp-taq test:unit       # Test taq
cd lib/js && yarn test                           # Test js library
```

## Linting

```bash
make lint                                        # Lint all libs and clients
pnpm --filter corteza-webapp-one lint            # Lint specific app
cd lib/js && yarn lint                           # Lint js library
```

## Building

```bash
cd client/web/one && pnpm build       # Build specific app (outputs to dist/)
cd client/web/compose && pnpm build
cd client/web/taq && pnpm build
```
