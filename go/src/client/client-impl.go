package client

import "fmt"

type TestClient struct {
	x, y int
}

func New() *TestClient {
	return &TestClient{1,2}
}

func (c *TestClient) Map() {
	fmt.Println("MAPPED")
}

func (c *TestClient) Reduce() {
	fmt.Println("REDUCED")
}