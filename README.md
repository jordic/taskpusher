
### TaskPusher

Just a small, single binary, task pusher, job queue, based on http.

#### Status. 
	
	**Alpha. Incomplet.** 
	
	Finished: Job queue, and webTasks
	
	Pending:
		+ Backend persistence.
		- Log handling
		- Http API


#### Motivation

Lots of job queues exists, with workers, persistence and distributed jobs. But on the end, for small projects with php, python or ruby, to have a single binary task pusher where you can push your existing long run jobs ( email sending, video recode, ... ) and deattach from the request/response loop is a big win.

All Job Queues I found, follow the same model, of a master queue, with some workers. This is a good design, for large scale projects, but for small ones, I found it too complicated. ( You end up with dependencies, and at leas with two more services to mantain, the queue, and the worker). Also in this model, you end up having to write a worker in any language.. (Deamonize it, log it, etc... ) Dropping the job runner, too the main application could be beneficial for developers. ( You don't loose the web request context)


#### Installation and run

go get github.com/jordic/taskpusher/cmd/tasker

Setup your command flags: 

-- Basicauth
-- Backend store
-- Log level (debug, production)
-- Network address
-- Concurrency ( How many goroutines do you want to be started to parallelize work?). Default 4.

Run the binary:

./tasker --basicauth=test:test -address=localhost:9900

And you are ready for accepting jobs. There is also a small web ui, to check status and healthy. Access:

localhost:9900/tasks and you should be on the interface.


#### API

How the hell I start submitting jobs?

POST to /add with
	url: url of the service to be called
	Will return a uid, that you can use later, to check your task status.
	Currently only GET urls are handled. Perhaps in new versions it is extended to POST and other http methods. 
	All executed tasks should return a http 200 status. If not, task is considered erroneous.
	

	
GET /job/:uid
	Check status and results.
	

GET /jobs/waiting
	Returns the list of waiting jobs
	
	/jobs/runing
	Jobs actually in process..
	
	/jobs/completed
	List of completed jobs. The results are paginated at a maximum of 100 registers per page. Next page are handled with the last_uid returned from db.
	
	/jobs/errors
	Should return a list of task errors.

Client libraries for popular languages will not be written. Perhaps a small library for python ( My other language of use). But on the end is too easy to write a lib, just fire a request, and you have your job done.



#### Backend and Persistence.

At the moment, the backend persistence layer is based on boltdb. An embedable key/value store without dependenciess for go. On future, adapters for mysql, or whatever else could be written.


#### Api and operation.

At the moment, every task must be a url call. This way, is easy for developers on project to actually handly the test cases on context and you don't have to setup extra dependencies. On future, also a command (unix way) could be enqueued. The api should be simple.



#### Guarantees.

When you stop, or reload the project, it will wait till queue is drilled. During this phase, tasks pushed get stored to disk, and loaded next reload. The same way, if there is a electric cut, on reload, the app looks for inconsistence and rebuilds it job table.


### Extending

Task pusher, can be used as a golang lib, to add task processing to your current app, or to be extended with your workers needs. Think in it, and taskpusher could be used also as a distributed job runner... Just extend the lib, fire some instance in distincts services, and shard your petitons to them.

