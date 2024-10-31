package good_util

import (
	"fmt"
	"github.com/samber/lo"
	lop "github.com/samber/lo/parallel"
	"math"
)

type Order struct {
	UserId int
}

func lopMap() {

	results := lop.Map([]Order{}, func(x Order, _ int) error {
		// update order
		return nil
	})

	fmt.Println(results)
}

type User struct {
	UserId int
	IsVip  bool
	Age    int
	Name   string
}

func loFilter() {
	users := []User{
		{UserId: 1, IsVip: true},
		{UserId: 2, IsVip: false},
	}

	filterVIPUser := func(u User, _ int) bool {
		return u.IsVip
	}
	vipUsers := lo.Filter(users, filterVIPUser)

	for _, user := range users {
		if user.IsVip {
			vipUsers = append(vipUsers, user)
		}
	}

	fmt.Println(vipUsers)
}

func loUniq() {
	users := []User{
		{UserId: 1, IsVip: true},
		{UserId: 2, IsVip: false},
		{UserId: 1, IsVip: true},
	}

	uniqUsers := lo.Uniq(users)

	fmt.Println(uniqUsers)
}

func lopPartitionBy() {
	users := []User{
		{UserId: 1, Age: 10},
		{UserId: 2, Age: 20},
		{UserId: 3, Age: 30},
		{UserId: 4, Age: 40},
		{UserId: 5, Age: 50},
		{UserId: 6, Age: 60},
	}

	agePartition := map[int]string{
		18:          "minor",
		35:          "youth",
		60:          "middle-aged",
		math.MaxInt: "senior",
	}

	partitions := lop.PartitionBy(users, func(user User) string {
		for age, partition := range agePartition {
			if user.Age < age {
				return partition
			}
		}
		return "unknown"
	})

	fmt.Println(partitions)
}

func loKeyBy() {
	users := []User{
		{UserId: 1, Age: 10},
		{UserId: 2, Age: 20},
		{UserId: 3, Age: 30},
		{UserId: 4, Age: 40},
		{UserId: 5, Age: 50},
		{UserId: 6, Age: 60},
	}

	userId2UserMap := lo.KeyBy(users, func(user User) int {
		return user.UserId
	})

	fmt.Println(userId2UserMap[1])
}

func loAssociate() {
	users := []User{
		{UserId: 1, Name: "10"},
		{UserId: 2, Name: "20"},
		{UserId: 3, Name: "30"},
		{UserId: 4, Name: "40"},
		{UserId: 5, Name: "50"},
		{UserId: 6, Name: "60"},
	}

	userId2Name := lo.Associate(users, func(user User) (userId int, name string) {
		return user.UserId, user.Name
	})

	fmt.Println(userId2Name[1])
}

func GetOrders() {
	orders := []Order{
		{UserId: 1},
		{UserId: 2},
		{UserId: 3},
		{UserId: 4},
		{UserId: 1},
	}

	userIds := lo.Map(orders, func(order Order, _ int) int {
		return order.UserId
	})
	fmt.Println(userIds)

	uniqUserIds := lo.Uniq(userIds)
	fmt.Println(uniqUserIds)

	getUserByIds(uniqUserIds)
}

func getUserByIds(userIds []int) {}
