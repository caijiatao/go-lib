package main

import (
	"sort"
	"time"
)

type SortKey string

const (
	Name                SortKey = "name"
	TransactionNum      SortKey = "transactionNum"
	LastTransactionDate SortKey = "lastTransactionDate"
)

type UserSortFunc func(userA, userB User) bool

func sortByName(userA, userB User) bool {
	return userA.Name < userB.Name
}

func sortByTransactionNum(userA, userB User) bool {
	return userA.TransactionNum < userB.TransactionNum
}

func sortByLastTransactionDate(userA, userB User) bool {
	return userA.LastTransactionDate.Before(userB.LastTransactionDate)
}

// 排序函数映射
var sortFunctions = map[SortKey]UserSortFunc{
	Name:                sortByName,
	TransactionNum:      sortByTransactionNum,
	LastTransactionDate: sortByLastTransactionDate,
}

type User struct {
	Name                string
	TransactionNum      int
	LastTransactionDate time.Time
}

func sortUsers(users []User, sortBy SortKey) []User {
	sortFunc, ok := sortFunctions[sortBy]
	// 这里可以根据自己的逻辑补充
	// 比如没有排序规则则返回原切片或者是使用默认排序均可
	if !ok {
		//...
	}

	sort.Slice(users, func(i, j int) bool {
		return sortFunc(users[i], users[j])
	})
	return users
}
