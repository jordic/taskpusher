

### TinyTaskPusher

Just a small job Queue the unix way.. 

Process add added to the queue:

POST /api/add
	json: {
		"command":xxx
	}
	@return process_id

add process, returns a process_id. 

That it can be used to check status:

GET /api/status/id 

{
	status: "running"
}

{
	status: "ended",
	result: "xxxxx"
}

{
	status: "waiting",
}


A Task should be a webhook ( For the moment )succ


type Task struct {
	Cmd string
	Params string
	status int
	Result sting
	Priority
	started time.Time
	ended time.Time
	UID
}


type Manager struct {
	Tasks []*Task
	MaxTasksRuning int
	Backend // Where to store poending tasks
	Runing
}


func (m *Manager) Add( t *Task ) {
	// Store task to backend..
	Backend.Store( Task )
	if Running < MaxTasks
		Start..
}

// Load stored pending tasks
func (m *Manager) Load() {
	
}



func Loop() {


	if NewTasks and not Manager.full() and not Stoping()
		PickupNewTask().Start()
	
	delay 100
}