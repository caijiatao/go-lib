package main

import (
	"sort"
)

func sortUsersOpt(users []User, sortFunc UserSortFunc) {
	sort.Slice(users, func(i, j int) bool {
		return sortFunc(users[i], users[j])
	})
}
