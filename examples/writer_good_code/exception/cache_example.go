package exception

func getUserFromCache() (*User, error) {
	return nil, nil
}

func getUserFromDB() (*User, error) {
	return nil, nil
}

func GerUser() (*User, error) {
	user, err := getUserFromCache()
	if err == nil {
		return user, nil
	}

	user, err = getUserFromDB()
	if err != nil {
		return nil, err
	}
	return user, nil
}
