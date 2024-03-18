package client

import (
	"context"
	"fmt"
	"github.com/Khan/genqlient/graphql"
	"net/http"
	"testing"
)

func Test_findTodos(t *testing.T) {
	ctx := context.Background()
	client := graphql.NewClient("http://localhost:8080/query", http.DefaultClient)
	createResp, err := createTodo(ctx, client)
	fmt.Println(createResp, err)
	resp, err := findTodos(ctx, client)
	fmt.Println(resp, err)
}
