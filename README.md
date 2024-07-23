# GitBeam

A service that pulls commits from github repo periodically. ( for this exercise it periodically pulls every x period of
time. )

Dependencies

- Go 1.24
- Sqlite DB
- Chi Router v5 ( as pure as net/http )

### Dependencies

> Gateway: https://github.com/Just4Ease/gitbeam
> Repo Manager Microservice: https://github.com/Just4Ease/gitbeam.repo.manager
> Commit Monitor Microservice: https://github.com/Just4Ease/gitbeam.commit.monitor

1. See the docs folder for more information on the workflow.
2. `go run main.go` it runs on port :8080
3. Import the postman collection to see and consume the endpoints.

Setup Process:

- Clone the three systems except the baselib.
- run the command below in each folder to pull baselib.

```bash
git submodule update --init --recursive --remote --checkout --force --rebase --recursive
```

- go run main.go 

Repo Manager Microservice runs on port 8001
Commit Monitor Microservice runs on port 8002
Gateway runs on Port 8080

I've attached the following system architecture and sequence diagrams.

I'm looking forward to hearing your feedback. 🚀
Kind Regards.