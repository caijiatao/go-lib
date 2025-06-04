package logic

import (
	"errors"
	"regexp"
)

// 校验密码强度（正则表达式版）
func IsValidPassword(password string) bool {
	pattern := `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d!@#$%^&*]{8,}$`
	matched, _ := regexp.MatchString(pattern, password)
	return matched
}

// ValidateRegister 用户注册校验
func ValidateRegister(username, password, address string) error {
	if len(username) < 3 || len(username) > 20 {
		return errors.New("用户名必须3~20字符")
	}

	//pattern := `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d!@#$%^&*]{8,}$`
	//matched, _ := regexp.MatchString(pattern, password)
	//if !matched {
	//	return errors.New("密码必须至少8位，且包含字母和数字")
	//}
	if !IsValidPassword(password) {
		return errors.New("密码必须至少8位，且包含字母和数字")
	}

	// 注册还要验证地址是否正确
	if err := CheckAddress(address); err != nil {
		return errors.New("地址格式不正确")
	}

	return nil
}

func CheckAddress(address string) error {
	return nil
}

// ValidateLogin 创建订单校验
func ValidateLogin(username, password string) error {
	if len(username) < 3 || len(username) > 20 {
		return errors.New("用户名必须3~20字符")
	}

	//pattern := `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d!@#$%^&*]{8,}$`
	//matched, _ := regexp.MatchString(pattern, password)
	//if !matched {
	//	return errors.New("密码必须至少8位，且包含字母和数字")
	//}
	if !IsValidPassword(password) {
		return errors.New("密码必须至少8位，且包含字母和数字")
	}

	return nil
}
