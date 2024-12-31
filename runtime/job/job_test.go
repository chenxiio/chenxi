package job

import (
	"fmt"
	"testing"
	"time"

	"github.com/chenxiio/chenxi/models"
)

func TestAddAction(t *testing.T) {
	job := CreateJobInstance()
	act := &Unitjob{}
	job.AddAction(act)
	// 添加断言来验证 AddAction 方法是否正确执行
	if len(job.ujs) != 1 {
		t.Errorf("AddAction did not add the action to the job")
	}
}
func TestExecuting(t *testing.T) {
	job := CreateJobInstance()
	act := &Unitjob{}
	job.AddAction(act)
	go job.executing()
	// 添加断言来验证 Executing 方法是否正确执行
	time.Sleep(time.Second) // 等待执行一段时间
	if len(job.ujs) != 0 {
		t.Errorf("Executing did not remove the action from the job")
	}
}
func TestTestuj(t *testing.T) {
	//Init("./", slog.LevelDebug)
	job := CreateJobInstance()
	act := &Unitjob{}

	job.AddAction(act)
	go job.executing()
	// 添加断言来验证 Testuj 方法是否正确执行
	time.Sleep(time.Second) // 等待执行一段时间
	// if len(job.ujs) != 0 {
	// 	t.Errorf("Testuj did not remove the action from the job")
	// }
	for i := 0; i < 1000; i++ {
		job.AddAction(&Unitjob{Unitjob: models.Unitjob{Unit: fmt.Sprintf("unit%d", i)}})
		// time.Sleep(time.Second * 5)
	}
	select {}
}
