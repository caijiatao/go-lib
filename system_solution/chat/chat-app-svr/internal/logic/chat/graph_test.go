package chat

import (
	"context"
	"testing"
)

func TestNeo4j(t *testing.T) {
	driver := connectNeo4j()
	defer driver.Close(context.Background())

	//prepareData(context.Background(), driver)
	queryFriendWhoUnderAge(context.Background(), driver)
	//query(context.Background(), driver)

}
