# jarvis (v2)

An ever-changing assortment of microservices and home automation services/experiments. Wrapped into a friendly chatbot (named Jarvis) and hosted in an Azure Kubernetes cluster.

## Deploy

Jarvis is hosted on Azure Kubernetes service via the [rollout spec](./rollout.yaml) in this repo.  Deploy with `kubectl apply -f rollout.yaml`.

## Developing 🚀

This repo is set up with a `.devcontainer.json` configuration for development in Codespaces.

Running the `Start Server` vscode task will build and run the `go` project under `./server`.

The task requires either:

- A complete `dev.env` environment variable file. See [example.env](./example.env) for an idea of what secrets are necessary.

- [Codespace repo-scoped secrets](https://docs.github.com/en/enterprise-cloud@latest/rest/codespaces/repository-secrets) for [each required secret](./example.env).

### Build Docker Image

```bash

$ docker build -t jarvis ./server/ -f ./Dockerfile

$ docker run -p 4000:4000 --env-file example.env jarvis

[GIN-debug] GET    /health                   --> main.Health (3 handlers)
[GIN-debug] POST   /knock                    --> main.Knock (4 handlers)
2022/05/08 21:08:37 Bot authorized on account 'dev_bot'
Serving at http://localhost:4000
```

## Related projects

[jarvis-apple-watch](https://github.com/joshspicer/jarvis-apple-watch)
