package logic

import (
	"errors"
	"regexp"
)

func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 20 {
		return errors.New("用户名必须3~20字符")
	}
	return nil
}

func ValidatePassword(password string) error {
	pattern := `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d!@#$%^&*]{8,}$`
	matched, _ := regexp.MatchString(pattern, password)
	if !matched {
		return errors.New("密码必须至少8位，且包含字母和数字")
	}
	return nil
}

func ValidateStrongPassword(password string) error {
	return nil
}

// ValidateRegister 用户注册校验
func ValidateRegister(username, password, address string) error {
	//if len(username) < 3 || len(username) > 20 {
	//	return errors.New("用户名必须3~20字符")
	//}
	if err := ValidateUsername(username); err != nil {
		return err
	}

	//pattern := `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d!@#$%^&*]{8,}$`
	//matched, _ := regexp.MatchString(pattern, password)
	//if !matched {
	//	return errors.New("密码必须至少8位，且包含字母和数字")
	//}
	if err := ValidatePassword(password); err != nil {
		return err
	}

	if err := ValidateStrongPassword(password); err != nil {
		return err
	}

	return nil
}

func CheckAddress(address string) error {
	return nil
}

// ValidateLogin 创建订单校验
func ValidateLogin(username, password string) error {
	//if len(username) < 3 || len(username) > 20 {
	//	return errors.New("用户名必须3~20字符")
	//}
	if err := ValidateUsername(username); err != nil {
		return err
	}

	//pattern := `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d!@#$%^&*]{8,}$`
	//matched, _ := regexp.MatchString(pattern, password)
	//if !matched {
	//	return errors.New("密码必须至少8位，且包含字母和数字")
	//}
	if err := ValidatePassword(password); err != nil {
		return err
	}

	return nil
}
