# Viscraft

Viscraft is an AI-powered product ad visual generator. Users can organize their work into projects, upload a product image, pick a scene style and mood, and let the AI generate a ready-to-use ad visual — all from a clean, focused workspace.

![Landing Page](viscraft-docs/apps-documentation/landing-page.png)

---

## Features

### Authentication
Register and log in with a JWT-based session. Your workspace and generated images are tied to your account.

![Register](viscraft-docs/apps-documentation/register-modal.png) ![Login](viscraft-docs/apps-documentation/login-modal.png)

### Project Management
Organize your work into campaigns/projects. Each project holds its own set of generated scenes.

![Create Campaign](viscraft-docs/apps-documentation/create-campaign-modal.png)

### AI Scene Generation
Upload a product photo, choose a scene style, set a mood and lighting, optionally add a reference image — and generate a product ad visual powered by Pollinations AI.

![Generate](viscraft-docs/apps-documentation/generate-modal.png)

### Workspace
Browse all your generated scenes in a grid layout. Each card shows the result with options to view details, regenerate, or delete.

![Workspace](viscraft-docs/apps-documentation/workspace-page.png)

### Scene Detail & Regenerate
View a full-size scene, inspect the prompt that was used, and regenerate with tweaked settings if needed.

![Detail](viscraft-docs/apps-documentation/detail-card-modal.png) ![Regenerate](viscraft-docs/apps-documentation/regenerate-modal.png)

### Onboarding Tour
First-time users get a guided tour of the workspace so they know exactly where to start.

![Tour 1](viscraft-docs/apps-documentation/apps-tour-1.png) ![Tour 2](viscraft-docs/apps-documentation/apps-tour-2.png) ![Tour 3](viscraft-docs/apps-documentation/apps-tour-3.png)

---

## Tech Stack

| Layer | Tech |
|---|---|
| Frontend | React, TypeScript, Vite, Chakra UI v3, Zustand, Axios |
| Backend | Go, Gin, PostgreSQL, JWT |
| AI | Pollinations AI (image generation) |
| Infrastructure | Docker, Docker Compose, Nginx |

---

## Project Structure

```
Viscraft/
├── viscraft-frontend/    # React app (Vite + TypeScript)
├── viscraft-backend/     # REST API (Go + Gin)
└── viscraft-docs/        # Screenshots and documentation
```
