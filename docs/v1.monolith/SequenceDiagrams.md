## Getting initial set of commits on getting and saving the repository from Github.

```mermaid
sequenceDiagram;
autonumber;
participant world
participant api
participant core
participant db
participant github
participant eventStore



world ->> api: GET /repo/ownerName/repoName
api ->> core: service.GetRepo(ownerName, repoName)
core ->> db: db.getRepoByNames()
db -->> core: nil
core ->> github: gh.getRepository({...})
github -->> core: repo
core ->> db: storeRepo(repo)
core -->> eventStore: store.publish(repo.created, repo)
rect rgb(43,120,211)
	eventStore ->> core: service.FetchAndSaveCommits
	note over eventStore, core: This is called to mirror the first set of commits from day 1 on that repo.
end
core -->> api: repo
api -->> world: {success, message, data: repo}

```





## Start Beaming - To tell the system to periodically check for commits.

```mermaid
sequenceDiagram;
autonumber;
participant world
participant api
participant cronService
participant db
participant eventStore
participant core


world ->> api: POST /commits/start-beaming 
Note over world, api: ({ ownerName, repoName, fromDate?, toDate? })
api ->> cronService: cronStore.SaveCronTracker(data)
cronService ->> db: db.SaveCronTracker(data)
db -->> cronService: saved data
cronService -->> eventStore: store.pub("cron.tracker.created", data)
Note over cronService, eventStore: After saving the data in the cron tracker table
rect rgb(10, 10, 10)
par background activity
eventStore ->> core: service.FetchAndSaveCommits(data)
Note over eventStore, core: fetch commits instantly before cron job's next tick.
end
end
cronService -->> api: data 
rect rgb(100, 100, 10)
par cronService activity
cronService ->> db: listCronTrackers()
Note over cronService, db: periodically, check db for repos we need to get commits for.
db -->> cronService: [tasks]
cronService -->> core: service.FetchAndSaveCommits(task)
Note over cronService, core: trigger a go routine that will fetch the commits for each repo task.
end
end
api -->> world: {success, message, data}
```



## Stop Beaming - To stop tracking new commits for a repo

```mermaid
sequenceDiagram;
autonumber;
participant world
participant api
participant cronService
participant db
participant eventStore

rect rgb(10, 10, 10)
world ->> api: POST /commits/stop-beaming 
Note over world, api: ({ ownerName, repoName })
api ->> cronService: cronStore.DeleteCronTracker(data)
cronService ->> db: db.DeleteCronTracker(data)
db -->> cronService: deleted_data
cronService -->> eventStore: store.pub("cron.tracker.deleted", deleted_data)
Note over cronService, eventStore: This event can be stored in an audit/event log.
cronService -->> api: bool

api -->> world: {success, message: "Successfully stopped mirroring commits"}
end
```







# Repo Manager

```mermaid
sequenceDiagram;
autonumber;
participant gRPC
participant core
participant db
participant github

gRPC ->> core: rpc.GetRepoByOwnerAndRepoName()
core ->> db: getByOwnerAndRepoName
rect rgb(0,100,0)
alt if found
	db -->> core: repo
end
end

rect rgb(0,0,0)
alt if not found in db, and fallsback to github.
	core ->> github: getRepo(owner, repo)
	rect rgb(0,0,100)
	alt if not found on github
		github -->> core: not found
	end
	end
	
	rect rgb(100,0,100)
	alt if found on github
		github -->> core: repo
		core -->> db: save repo
	end
	end
end
end

core -->> gRPC: repo
```





# Commit Monitor ( Start Process )

```mermaid
sequenceDiagram;
autonumber;

participant gRPC
participant scheduler
participant db.cron_tracker
participant eventStore
participant core
participant github

gRPC ->> scheduler: startMonitoringRepo(monitor_payload)
scheduler ->> db.cron_tracker: save(monitor_payload)

rect rgb(90,50,0)
par save to job
scheduler ->> scheduler: registerJob(monitor_payload)
Note left of scheduler: Register job with the duration in monitor_payload (this can always be updated via the input)
scheduler ->> core: service.FetchAndSaveCommit(monitor_payload)
Note over scheduler, core: does not use the provided date range, so that it keeps tracking new changes from github.
end
end

rect rgb(0,0,0)
par event handling in parallel
scheduler ->> eventStore: monitor_payload
eventStore ->> core: service.FetchAndSaveCommit(monitor_payload)
Note over eventStore, core: uses the date range filter in the monitor payload to fetch if available.
end
end

```

