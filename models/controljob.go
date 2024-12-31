package models

type ControlJob struct {
	CONTROL_JOB_ID   string
	CARRIER_ID       []string
	PRIORITY         int
	PROCESS_JOB_LIST []string
	Mode             string // Auto ,Manual
	State            string //QUEUED，EXECUTING，PAUSED,ABORTED,ENDING,COMPLETED
	Createtime       int64
	Starttime        int64
	Endtime          int64
}
