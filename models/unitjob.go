package models

type Unitjob struct {
	Unitjob_id     int64
	PROCESS_JOB_ID string
	MATERIAL_LIST  []string
	PreUnit        string // 上一个单元
	Unit           string // place 成功之后 设置
	StepName       int    //
	//Recipe         string // 通过unit 和processrecipe 获得 unitrecipe ，通过 processrecipe 获取单元路径
	State      string
	Createtime int64
	Starttime  int64
	Endtime    int64
}
