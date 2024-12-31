package models

type ProcessJob struct {
	PROCESS_JOB_ID        string
	PROCESS_JOB_NAME      string
	PROCESS_DEFINITION_ID string
	PARAMETER_LIST        map[string]string // recipe list,recipe1,recipe2
	MATERIAL_LIST         []string          // A1,B2,C3  A1~A17
	PRIORITY              string
	PROCESS_JOB_NOTES     string
	State                 string //QUEUED，EXECUTING，PAUSED,ABORTED,ENDING,COMPLETED
	Createtime            int64
	Starttime             int64
	Endtime               int64
}
