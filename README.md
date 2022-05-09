# jarvis
jarvis v2

# Docker Image

```bash

$ docker build -t jarvis ./server/ -f ./Dockerfile

$ docker run -p 4000:4000 --env-file example.env jarvis

[GIN-debug] GET    /health                   --> main.Health (3 handlers)
[GIN-debug] POST   /knock                     --> main.Knock (4 handlers)
2022/05/08 21:08:37 Bot authorized on account 'dev_bot'
Serving at http://localhost:4000
```