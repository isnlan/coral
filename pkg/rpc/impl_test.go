package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type (
	Data struct {
		Url   string
		Type  string
		Count int
	}

	Auth struct {
		ID       string
		Expire   time.Time
		DataList []*Data
	}

	Status string
)

type MockService struct {
}

func (m *MockService) Milt(ath *Auth, id string, idx int, data Data, mp map[string]interface{}, stat Status) int {
	fmt.Printf("auth -> %+v\n", ath)
	fmt.Printf("id -> %+v\n", id)
	fmt.Printf("idx -> %+v\n", idx)
	fmt.Printf("data -> %+v\n", data)
	fmt.Printf("map -> %+v\n", mp)
	fmt.Printf("status -> %+v\n", stat)
	return idx
}

type Cli struct {
	Id   int
	Name string
}

func (m *MockService) DefaultArgsWithArg(c *Cli, name string) *Cli {
	if c.Name != "s" {
		c.Name = name
	}
	return c
}

func (m *MockService) DefaultArgs(c Cli) *Cli {
	return &c
}

func (m *MockService) ReturnErr(c *Cli, err bool) (int, error) {
	if err {
		return 10, errors.New("set err")
	} else {
		return 1, nil
	}
}

func TestNewServer(t *testing.T) {
	server := New()
	err := server.Register(&MockService{})
	assert.NoError(t, err)
}

type Type1 struct{}

func (t *Type1) F1() int {
	return 0
}

func TestNewServer2(t *testing.T) {
	server := New()
	err := server.Register(Type1{})
	assert.Error(t, err)
}

func TestServer_Handler(t *testing.T) {
	server := New()
	err := server.Register(&MockService{})
	assert.NoError(t, err)

	req := &ClientRequest{
		Method: "MockService.Milt",
		Params: []interface{}{
			&Auth{
				ID:     "t1",
				Expire: time.Now(),
				DataList: []*Data{
					&Data{
						Url:   "www.baidu.com",
						Type:  "search",
						Count: 10,
					},
					&Data{
						Url:   "www.doc.com",
						Type:  "doc",
						Count: 9,
					},
				},
			},
			"sfs",
			12,
			Data{
				Url:   "www.google.com",
				Type:  "go",
				Count: 100,
			},
			map[string]interface{}{
				"d1": 12,
				"d2": Data{
					Url:   "www.kd.com",
					Type:  "dta",
					Count: 100,
				},
			},
			"statue",
		},
	}
	bytes, _ := json.Marshal(req)
	var req1 Request
	err = json.Unmarshal(bytes, &req1)
	assert.NoError(t, err)

	ret, err := server.Handler(&req1)

	assert.NoError(t, err)
	assert.Equal(t, ret, 12)
}

func TestServer_Handler2(t *testing.T) {
	server := New()
	err := server.Register(&MockService{})
	assert.NoError(t, err)

	req := &ClientRequest{
		Method: "MockService.DefaultArgsWithArg",
		Params: []interface{}{
			"lucy",
		},
	}
	bytes, _ := json.Marshal(req)
	var req1 Request
	err = json.Unmarshal(bytes, &req1)
	assert.NoError(t, err)

	ret, err := server.Handler(&req1, &Cli{
		Id:   10,
		Name: "s1",
	})

	assert.NoError(t, err)
	assert.Equal(t, ret, &Cli{
		Id:   10,
		Name: "lucy",
	})
}

func TestServer_Handler3(t *testing.T) {
	server := New()
	err := server.Register(&MockService{})
	assert.NoError(t, err)

	req := &ClientRequest{
		Method: "MockService.DefaultArgs",
		Params: []interface{}{},
	}
	bytes, _ := json.Marshal(req)
	var req1 Request
	err = json.Unmarshal(bytes, &req1)
	assert.NoError(t, err)

	ret, err := server.Handler(&req1, Cli{
		Id:   10,
		Name: "s1",
	})

	assert.NoError(t, err)
	assert.Equal(t, ret, &Cli{
		Id:   10,
		Name: "s1",
	})
}

func TestServer_Handler4(t *testing.T) {
	server := New()
	err := server.Register(&MockService{})
	assert.NoError(t, err)

	bytes, _ := json.Marshal(&ClientRequest{
		Method: "MockService.ReturnErr",
		Params: []interface{}{
			false,
		},
	})
	var req1 Request
	err = json.Unmarshal(bytes, &req1)
	assert.NoError(t, err)

	ret, err := server.Handler(&req1, &Cli{
		Id:   10,
		Name: "s1",
	})

	assert.NoError(t, err)
	assert.Equal(t, ret, 1)

	bytes, _ = json.Marshal(&ClientRequest{
		Method: "MockService.ReturnErr",
		Params: []interface{}{
			true,
		},
	})
	var req2 Request
	err = json.Unmarshal(bytes, &req2)
	assert.NoError(t, err)

	ret1, err1 := server.Handler(&req2, &Cli{
		Id:   10,
		Name: "s1",
	})

	assert.Error(t, err1)
	assert.Nil(t, ret1)
}
