package map_demo

var (
	userIdMap = make(map[string]struct{})
)

func BuildUserIdMap(userIds []string) (userIdMap map[string]struct{}) {
	userIdMap = make(map[string]struct{})

	for _, userId := range userIds {
		userIdMap[userId] = struct{}{}
	}

	return userIdMap
}

func UserIdExists(userIdExists string) bool {
	_, ok := userIdMap[userIdExists]
	return ok
}
