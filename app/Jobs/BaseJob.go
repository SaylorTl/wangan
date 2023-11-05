package Jobs

type BaseJobs struct {
}

func (b BaseJobs) Init() {
	LoopholeStatusCheckJob = &loopholestatuscheckjob{}
}
