# jarvis (v2)

A home automation/virtual assistant.

## Deploy

Jarvis is hosted on Azure Kubernetes service via the [rollout spec](./rollout.yaml) in this repo.  The [deploy-to-cluster.yaml](./.github/workflows/deploy-to-cluster.yaml) workflow is triggers on pushes to `main`.

## Developing

This repo is set up with a `.devcontainer.json` configuration, for development in Codespaces or Remote-Containers.

Running the 'Start Server' vscode task will build and run the `go` project under `./server`.  The task expect a `dev.env` environment variable file.  See [example.env](./example.env) for an idea of what secrets are necessary.

### Build Docker Image

```bash

$ docker build -t jarvis ./server/ -f ./Dockerfile

$ docker run -p 4000:4000 --env-file example.env jarvis

[GIN-debug] GET    /health                   --> main.Health (3 handlers)
[GIN-debug] POST   /knock                    --> main.Knock (4 handlers)
2022/05/08 21:08:37 Bot authorized on account 'dev_bot'
Serving at http://localhost:4000
```