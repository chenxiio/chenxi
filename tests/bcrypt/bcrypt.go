package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "myPassword123" // 要加密的密码
	// 生成密码的哈希值
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("密码加密失败:", err)
		return
	}
	fmt.Println("密码哈希值:", string(hashedPassword))
	// 验证密码
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		fmt.Println("密码验证失败:", err)
		return
	}
	fmt.Println("密码验证通过")
}
