package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/graphql-go/graphql"
)

type Todo struct {
	Title string `json:"title"`
	Id string `json:"id"`
	Completed bool `json:"completed"`
}

var data []Todo = []Todo {
	Todo{
		Title:     "this is a todo",
		Id:        "0",
		Completed: false,
	},
	Todo{
		Title:     "this is a todo",
		Id:        "1",
		Completed: false,
	},
	Todo{
		Title:     "this is a todo",
		Id:        "2",
		Completed: false,
	},
	Todo{
		Title:     "this is a todo",
		Id:        "3",
		Completed: false,
	},
	Todo{
		Title:     "this is a todo",
		Id:        "4",
		Completed: false,
	},
	Todo{
		Title:     "this is a todo",
		Id:        "5",
		Completed: false,
	},
}

type Edge struct {
	Cursor string `json:"cursor"`
	Node Todo `json:"node"`
}

type PageInfo struct {
	EndCursor string `json:"endCursor"`
	HasNextPage bool `json:"hasNextPage"`
}

type TodosResultCursor struct {
	Edges []Edge `json:"edges"`
	PageInfo PageInfo `json:"pageInfo"`
	TotalCount int `json:"totalCount"`
}
func main() {

	todoType:= graphql.NewObject(graphql.ObjectConfig{
		Name:        "Todo",
		Fields:      graphql.Fields{
			"id": &graphql.Field{
				Type:              graphql.String,
			},
			"title": &graphql.Field{
				Type:              graphql.String,
			},
			"completed": &graphql.Field{
				Type:              graphql.Boolean,
			},
		},
	})

	edgeType:= graphql.NewObject(graphql.ObjectConfig{
		Name:        "Edge",
		Fields:      graphql.Fields{
			"cursor": &graphql.Field{
				Type:              graphql.String,
			},
			"node": &graphql.Field{
				Type:              todoType,
			},
		},
	})

	PageInfoType:= graphql.NewObject(graphql.ObjectConfig{
		Name:        "Edge",
		Fields:      graphql.Fields{
			"endCursor": &graphql.Field{
				Type:              graphql.String,
			},
			"hasNextPage": &graphql.Field{
				Type:              graphql.Boolean,
			},
		},
	})

	TodosResultCursorType:= graphql.NewObject(graphql.ObjectConfig{
		Name:        "Todo",
		Fields:      graphql.Fields{
			"edges": &graphql.Field{
				Type:              graphql.NewList(edgeType),
			},
			"totalCount": &graphql.Field{
				Type:              graphql.Int,
			},
			"pageInfo": &graphql.Field{
				Type:              PageInfoType,
			},
		},
	})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Query",
		Fields:      graphql.Fields{
			"data": &graphql.Field{
				Type:              graphql.NewList(todoType),
				Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {
					return data, nil
				},
			},
			"allTodosCursor": &graphql.Field{
				Type:              TodosResultCursorType,
				Args:              graphql.FieldConfigArgument{
					"after": &graphql.ArgumentConfig{
						Type:         graphql.String,
					},
					"first": &graphql.ArgumentConfig{
						Type:         graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {

					var after int
					if p.Args["after"] == nil {
						after,_ = strconv.Atoi(string(0))
					} else {
						temp, _ :=b64.StdEncoding.DecodeString(p.Args["after"].(string))
						after, _ = strconv.Atoi(string(temp))
					}

					var first int

					if p.Args["first"] == nil {
						first = 5
					} else {
						first = p.Args["first"].(int)
					}

					start :=after
					page := data[after:after+first]
					totalCount := len(data)
					hasNextPage := start + first < totalCount


					var nodes[]Edge
					var endCursor string
					for _, edge := range page {
						endCursor = b64.StdEncoding.EncodeToString([]byte(edge.Id))
						edge := Edge{Cursor:endCursor, Node:edge}
						nodes = append(nodes, edge)
					}

					pageInfo := PageInfo{EndCursor: endCursor, HasNextPage: hasNextPage}

					todosResultCursor := TodosResultCursor{ PageInfo:pageInfo, Edges:nodes, TotalCount: totalCount }

					return todosResultCursor, nil
				},
			},
		},
	})

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:         rootQuery,
	})

	query := `
		{
			allTodosCursor(first: 3, after: "Mw==") {
				edges {
				  node {
					id,
					title,
					completed
				  }
				  cursor
				}
				pageInfo {
				  endCursor,
				  hasNextPage
				}
				totalCount
		  	}
		}
	`
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON)
}