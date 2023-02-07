---
sidebar_position: 2
---

# Roadmap

This is our roadmap for supported features. Here you can track what features
did we already implemented and what we're planning to do next. Feel free to
reach out in our Slack or Github if you have a particular interest in any
planned feature or have a suggestion for a new one.

Legend:

✅ - Supported well
🟨 - Supported, but not implemented fully
⏰ - Planned

## Backend services

| Features                 | Go Template | Python Template (Beta) |
|--------------------------|-------------|------------------------|
| OpenAPI                  |✅           | ✅                     |
| Static Configs           |✅           | ✅                     |
| Dynamic Configs          |✅           | ✅                     |
| Postgres                 |✅           | ⏰                     |
| Structured logging       |✅           | ✅                     |
| OpenAPI server metrics   |✅           | ✅                     |
| OpenAPI clients metrics  |✅           | ✅                     |
| Graceful shutdown        |⏰           | ⏰                     |
| Authentication           |⏰           | ⏰                     |
| Periodic tasks           |⏰           | ⏰                     |
| API Gateway              |⏰           | ⏰                     |
| Message queues           |⏰           | ⏰                     |
| mify run helper          |🟨           | ⏰                     |

## Frontend services

| Features                 | React TS Template | NuxtJS Template  |
|--------------------------|-------------------|------------------|
| OpenAPI Clients          |✅                 | ✅               |
| Static Configs           |🟨                 | 🟨               |
| Structured logging       |⏰                 | ⏰               |
| Authentication           |⏰                 | ⏰               |

## Other

| Features                 | Status            |
|--------------------------|-------------------|
| Cron jobs                |⏰                 |

## Mify Cloud

| Features                         | Status |
|----------------------------------|--------|
| Cloud services deploy            |✅      |
| Postgres migrations              |✅      |
| Static configs                   |✅      |
| Consul for dynamic configs       |⏰      |
| ELK logs collection and querying |⏰      |
| Metrics collection               |⏰      |
| Multiple regions deploy          |⏰      |
| Canary rollout                   |⏰      |
| Quick rollback                   |⏰      |
