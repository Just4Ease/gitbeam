# GitBeam

A service that pulls commits from github repo periodically. ( for this exercise it periodically pulls every x period of
time. )

Dependencies

- Go 1.24
- Sqlite DB
- Chi Router v5 ( as pure as net/http )
- gRPC




#### System codebase breakdown
> Gateway: https://github.com/Just4Ease/gitbeam

> Repo Manager Microservice: https://github.com/Just4Ease/gitbeam.repo.manager

> Commit Monitor Microservice: https://github.com/Just4Ease/gitbeam.commit.monitor

> Base Lib: https://github.com/Just4Ease/gitbeam.baselib

### How to start the service.
#### Clone this repository, and run the command below.
```shell
git clone https://github.com/Just4Ease/gitbeam
cd gitbeam
./setup.sh
```

### How to test the service
```shell
cd gitbeam
make test

# or
cd gitbeam
go test -v ./... -cover
```

- Repo Manager Microservice runs on port 8001
- Commit Monitor Microservice runs on port 8002
- Gateway runs on Port 8080

I've attached the following system architecture and sequence diagrams.

See docs/v2.microservices/

### System Architecture Diagram
![System Architecture](docs/v2.microservices/system_architecture.png)

---

### Repo Manager Workflow Diagram
![Repo Manager](docs/v2.microservices/repo_manager.png)

---

### Commit Monitor Workflow Diagram
![Commit Monitor](docs/v2.microservices/commit_monitor.png)


### Video Demo of how it works.
[Click Here To Watch Video Demo: https://drive.google.com/file/d/1R8E0pVdYpNkQ2dzXLC0y_zYEXCzRTUNS/view?usp=sharing](https://drive.google.com/file/d/1R8E0pVdYpNkQ2dzXLC0y_zYEXCzRTUNS/view?usp=sharing)


I'm looking forward to hearing your feedback. 🚀
Kind Regards.