package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 要启动的进程和参数
	cmd := exec.Command("../build/testproj/testproj.exe", "arg1", "arg2")
	// 设置进程的输出和错误输出
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 启动进程
	err := cmd.Start()
	if err != nil {
		fmt.Printf("启动进程时遇到错误：%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("进程已启动，进程ID为：%d\n", cmd.Process.Pid)
	// 监视进程状态
	go func() {
		err := cmd.Wait()
		if err != nil {
			exitErr, ok := err.(*exec.ExitError)
			if ok {
				// 如果进程以非零状态退出，则在此处处理错误
				status := exitErr.Sys().(syscall.WaitStatus)
				fmt.Printf("进程以非零状态退出，退出状态码为：%d\n", status.ExitStatus())
			} else {
				fmt.Printf("等待进程退出时遇到错误：%v\n", err)
			}
		} else {
			fmt.Println("进程已成功退出")
		}
	}()
	// 持续运行，直到接收到终止信号
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt, syscall.SIGTERM)
	<-terminate
	// 终止进程
	err = cmd.Process.Signal(syscall.SIGTERM)
	if err != nil {
		fmt.Printf("终止进程时遇到错误：%v\n", err)
		os.Exit(1)
	}
	// 等待进程退出
	<-time.After(5 * time.Second)
	if cmd.ProcessState != nil && !cmd.ProcessState.Exited() {
		fmt.Println("进程无法正常退出，强制终止")
		err = cmd.Process.Kill()
		if err != nil {
			fmt.Printf("强制终止进程时遇到错误：%v\n", err)
			os.Exit(1)
		}
	}
	fmt.Println("程序已终止")
}
