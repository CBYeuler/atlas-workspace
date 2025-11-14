# Atlas Workspace — v0.1 Foundations

This is the foundation phase of the Atlas Workspace monorepo.

## Structure
- backend: Go (Gin)
- frontend: Next.js 15 (TypeScript + Tailwind)
- devops: Docker Compose, infra scripts

## Commands

Start environment:

```bash
docker compose up -d --build
```

Health Check:
```bash
GET /healthz → {"status":"ok"}
```

Next.js Dev Server:
```bash
npm run dev
```