package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// 全局变量：Redis 客户端、上下文
var (
	redisClient *redis.Client
	ctx         = context.Background()
	// Redis Key 前缀（支持多租户时可加租户ID，如 "tenant:100:role:btn:"）
	roleBtnKeyPrefix = "role:btn:"
	// 用户权限缓存 Key 前缀（存储用户合并后的按钮权限）
	userBtnKeyPrefix = "user:btn:"
	// 权限缓存过期时间（1小时，与之前的权限版本号方案对齐）
	cacheTTL = time.Hour * 1
)

// 初始化 Redis 客户端
func initRedis() error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // 无密码
		DB:       0,  // 使用第0个数据库
		PoolSize: 10, // 连接池大小
	})
	// 测试连接
	return redisClient.Ping(ctx).Err()
}

// ------------------------------
// 1. 角色按钮权限管理：给角色添加按钮权限
// roleID: 角色ID（如 1001）
// btnCodes: 按钮权限标识列表（如 ["btn:order_submit", "btn:order_list"]）
// ------------------------------
func addRoleBtnPermissions(roleID string, btnCodes []string) error {
	key := roleBtnKeyPrefix + roleID
	// 用 SADD 往 Redis Set 中添加按钮权限（自动去重）
	err := redisClient.SAdd(ctx, key, btnCodes).Err()
	if err != nil {
		return fmt.Errorf("add role btn permissions failed: %v", err)
	}
	// 设置过期时间（可选，避免Redis内存溢出）
	_ = redisClient.Expire(ctx, key, cacheTTL).Err()
	fmt.Printf("角色[%s]添加按钮权限成功：%v\n", roleID, btnCodes)
	return nil
}

// ------------------------------
// 2. 角色按钮权限管理：删除角色的某个按钮权限
// ------------------------------
func removeRoleBtnPermission(roleID string, btnCode string) error {
	key := roleBtnKeyPrefix + roleID
	err := redisClient.SRem(ctx, key, btnCode).Err()
	if err != nil {
		return fmt.Errorf("remove role btn permission failed: %v", err)
	}
	fmt.Printf("角色[%s]删除按钮权限成功：%s\n", roleID, btnCode)
	return nil
}

// ------------------------------
// 3. 用户登录：合并所有角色的按钮权限，返回最终权限列表
// userID: 用户ID（如 "u001"）
// roleIDs: 用户关联的角色ID列表（如 ["1001", "1002"]）
// ------------------------------
func userLoginMergeBtnPermissions(userID string, roleIDs []string) ([]string, error) {
	// 步骤1：构建所有角色的 Redis Key
	roleKeys := make([]string, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		roleKeys = append(roleKeys, roleBtnKeyPrefix+roleID)
	}

	// 步骤2：用 SUNION 合并所有角色的按钮权限（Redis 原生支持，高效去重）
	userBtnKey := userBtnKeyPrefix + userID
	err := redisClient.SUnionStore(ctx, userBtnKey, roleKeys...).Err()
	if err != nil {
		return nil, fmt.Errorf("merge role btn permissions failed: %v", err)
	}
	// 设置用户权限缓存过期时间
	_ = redisClient.Expire(ctx, userBtnKey, cacheTTL).Err()

	// 步骤3：获取合并后的按钮权限列表
	btnPermissions, err := redisClient.SMembers(ctx, userBtnKey).Result()
	if err != nil {
		return nil, fmt.Errorf("get user btn permissions failed: %v", err)
	}

	fmt.Printf("\n用户[%s]登录成功，关联角色：%v\n", userID, roleIDs)
	fmt.Printf("用户[%s]最终按钮权限：%v\n", userID, btnPermissions)
	return btnPermissions, nil
}

// ------------------------------
// 4. 已登录用户：刷新按钮权限（权限变更后调用）
// ------------------------------
func refreshUserBtnPermissions(userID string, roleIDs []string) ([]string, error) {
	fmt.Printf("\n用户[%s]刷新权限...\n", userID)
	// 先删除旧缓存
	_ = redisClient.Del(ctx, userBtnKeyPrefix+userID).Err()
	// 重新合并角色权限（复用登录逻辑）
	return userLoginMergeBtnPermissions(userID, roleIDs)
}

// ------------------------------
// 5. 前端渲染模拟：根据用户权限决定按钮是否显示
// ------------------------------
func mockFrontendRender(userID string, btnPermissions []string) {
	fmt.Printf("\n=== 前端渲染模拟（用户：%s）===\n", userID)
	// 系统所有按钮（实际场景从配置/数据库获取）
	allBtns := []struct {
		Code string // 按钮标识
		Name string // 按钮名称
	}{
		{"btn:order_submit", "提交订单"},
		{"btn:order_list", "订单列表"},
		{"btn:order_delete", "删除订单"},
		{"btn:order_export", "导出订单"},
	}

	// 遍历所有按钮，判断是否有权限
	for _, btn := range allBtns {
		hasPermission := false
		for _, p := range btnPermissions {
			if p == btn.Code {
				hasPermission = true
				break
			}
		}
		if hasPermission {
			fmt.Printf("[显示] 按钮：%s（标识：%s）\n", btn.Name, btn.Code)
		} else {
			fmt.Printf("[隐藏] 按钮：%s（标识：%s）\n", btn.Name, btn.Code)
		}
	}
	fmt.Println("===========================\n")
}

func main() {
	// 1. 初始化 Redis
	if err := initRedis(); err != nil {
		panic(fmt.Sprintf("Redis 初始化失败：%v", err))
	}
	fmt.Println("Redis 初始化成功！")

	// 2. 模拟配置角色按钮权限
	// 角色1001（普通用户）：提交订单、订单列表
	_ = addRoleBtnPermissions("1001", []string{"btn:order_submit", "btn:order_list"})
	// 角色1002（部门管理员）：提交订单、删除订单
	_ = addRoleBtnPermissions("1002", []string{"btn:order_submit", "btn:order_delete"})
	// 角色1003（超级管理员）：所有按钮
	_ = addRoleBtnPermissions("1003", []string{"btn:order_submit", "btn:order_list", "btn:order_delete", "btn:order_export"})

	// 3. 模拟用户登录（用户u001关联角色1001+1002）
	userID := "u001"
	userRoles := []string{"1001", "1002"} // 普通用户+部门管理员
	btnPermissions, _ := userLoginMergeBtnPermissions(userID, userRoles)
	// 模拟前端渲染
	mockFrontendRender(userID, btnPermissions)

	// 4. 模拟权限变更：给角色1002添加「导出订单」按钮
	_ = addRoleBtnPermissions("1002", []string{"btn:order_export"})
	// 已登录用户刷新权限
	updatedBtnPermissions, _ := refreshUserBtnPermissions(userID, userRoles)
	// 模拟前端重新渲染
	mockFrontendRender(userID, updatedBtnPermissions)

	// 5. 模拟权限变更：删除角色1001的「提交订单」按钮
	_ = removeRoleBtnPermission("1001", "btn:order_submit")
	// 已登录用户刷新权限
	finalBtnPermissions, _ := refreshUserBtnPermissions(userID, userRoles)
	// 模拟前端重新渲染
	mockFrontendRender(userID, finalBtnPermissions)

	// 6. 测试超级管理员用户（u002关联角色1003）
	superUserID := "u002"
	superUserRoles := []string{"1003"}
	superBtnPermissions, _ := userLoginMergeBtnPermissions(superUserID, superUserRoles)
	mockFrontendRender(superUserID, superBtnPermissions)
}
