package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"crypto/rand"
	"fmt"
	"gqlgen-todos/graph/model"
	"math/big"
)

// CreateTodo is the resolver for the createTodo field.
func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	randNumber, _ := rand.Int(rand.Reader, big.NewInt(100))
	todo := &model.Todo{
		Text:   input.Text,
		ID:     fmt.Sprintf("T%d", randNumber),
		UserID: input.UserID,
	}
	r.todos = append(r.todos, todo)
	return todo, nil
}

// CreateTime is the resolver for the CreateTime field.
func (r *orderResolver) CreateTime(ctx context.Context, obj *model.Order) (string, error) {
	panic(fmt.Errorf("not implemented: CreateTime - CreateTime"))
}

// PaymentTime is the resolver for the PaymentTime field.
func (r *orderResolver) PaymentTime(ctx context.Context, obj *model.Order) (string, error) {
	panic(fmt.Errorf("not implemented: PaymentTime - PaymentTime"))
}

// Todos is the resolver for the todos field.
func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	return r.todos, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	for _, todo := range r.todos {
		if todo.UserID == id {
			return &model.User{ID: todo.UserID, Name: "user " + todo.UserID}, nil
		}
	}
	return nil, nil
}

// Products is the resolver for the products field.
func (r *queryResolver) Products(ctx context.Context, productID *string) ([]*model.Product, error) {
	if productID != nil {
		return []*model.Product{
			{
				ProductId:   *productID,
				ProductName: "product " + *productID,
			},
		}, nil
	}
	return []*model.Product{}, nil
}

// Orders is the resolver for the orders field.
func (r *queryResolver) Orders(ctx context.Context) ([]*model.Order, error) {
	panic(fmt.Errorf("not implemented: Orders - orders"))
}

// User is the resolver for the user field.
func (r *todoResolver) User(ctx context.Context, obj *model.Todo) (*model.User, error) {
	return &model.User{ID: obj.UserID, Name: "user " + obj.UserID}, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Order returns OrderResolver implementation.
func (r *Resolver) Order() OrderResolver { return &orderResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Todo returns TodoResolver implementation.
func (r *Resolver) Todo() TodoResolver { return &todoResolver{r} }

type mutationResolver struct{ *Resolver }
type orderResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type todoResolver struct{ *Resolver }