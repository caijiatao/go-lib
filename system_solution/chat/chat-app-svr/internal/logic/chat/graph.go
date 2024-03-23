package chat

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func connectNeo4j() (driver neo4j.DriverWithContext) {
	ctx := context.Background()
	dbUri := "neo4j+s://74d2626d.databases.neo4j.io"
	dbUser := "neo4j"
	dbPassword := ""
	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))
	//defer driver.Close(ctx)

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}
	return driver
}

func prepareData(ctx context.Context, driver neo4j.DriverWithContext) {
	peoples := []map[string]any{
		{"name": "Alice", "age": 42, "friends": []string{"Bob", "Peter", "Anna"}},
		{"name": "Bob", "age": 19},
		{"name": "Peter", "age": 50},
		{"name": "Anna", "age": 30},
	}

	// Create some nodes
	for _, person := range peoples {
		_, err := neo4j.ExecuteQuery(ctx, driver,
			"MERGE (p:Person {name: $person.name, age: $person.age})",
			map[string]any{
				"person": person,
			}, neo4j.EagerResultTransformer,
			neo4j.ExecuteQueryWithDatabase("neo4j"))
		if err != nil {
			panic(err)
		}
	}

	// Create some relationships
	for _, person := range peoples {
		if person["friends"] != "" {
			_, err := neo4j.ExecuteQuery(ctx, driver, `
                MATCH (p:Person {name: $person.name})
                UNWIND $person.friends AS friend_name
                MATCH (friend:Person {name: friend_name})
                MERGE (p)-[:KNOWS]->(friend)
                `, map[string]any{
				"person": person,
			}, neo4j.EagerResultTransformer,
				neo4j.ExecuteQueryWithDatabase("neo4j"))
			if err != nil {
				panic(err)
			}
		}
	}
}

func queryFriendWhoUnderAge(ctx context.Context, driver neo4j.DriverWithContext) {
	// Retrieve Alice's friends who are under 40
	result, err := neo4j.ExecuteQuery(ctx, driver, `
        MATCH (p:Person {name: $name})-[:KNOWS]-(friend:Person)
        WHERE friend.age < $age
        RETURN friend
        `, map[string]any{
		"name": "Alice",
		"age":  40,
	}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))
	if err != nil {
		panic(err)
	}

	// Loop through results and do something with them
	for _, record := range result.Records {
		person, _ := record.Get("friend")
		fmt.Println(person)
		// or
		fmt.Println(record.AsMap())
	}
}

func query(ctx context.Context, driver neo4j.DriverWithContext) {
	// Get the name of all 42 year-olds
	result, _ := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (cust:Customer)-[:PURCHASED]->(:Order)-[o:ORDERS]->(p:Product),
      (p)-[:PART_OF]->(c:Category {categoryName:$categoryName})
RETURN cust.contactName as CustomerName,
       sum(o.quantity) AS TotalProductsPurchased`,
		map[string]any{
			"categoryName": "Produce",
		}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))

	// Loop through results and do something with them
	for _, record := range result.Records {
		fmt.Println(record.AsMap())
	}

	// Summary information
	fmt.Printf("The query `%v` returned %v records in %+v.\n",
		result.Summary.Query().Text(), len(result.Records),
		result.Summary.ResultAvailableAfter())
}
