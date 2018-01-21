package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
)

type Client_V1 struct {
	URL           *url.URL
	TeamName      string
	HTTPClient    *http.Client
	Authorization string
	Logger        *log.Logger
}

type Categories struct {
	Categories []Category
}

type Category struct {
	Name     string     `json:"name"`
	Post     bool       `json:"post,omitempty"`
	Count    int        `json:"count,omitempty"`
	Children []Category `json:"children,omitempty"`
}

func (c Category) Tree() []string {
	if c.Post {
		return nil
	}
	var trees []string

	if len(c.Children) == 0 {
		trees = append(trees, c.Name+"/")
		return trees
	}
	for _, child := range c.Children {
		trees = append(trees, c.Append(child.Tree())...)
	}
	return trees
}

func (c Category) Append(trees []string) []string {
	var res []string
	for _, tree := range trees {
		res = append(res, c.Name+"/"+tree)
	}
	return res
}

func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}

func NewClient(auth string, teamName string) *Client_V1 {
	client := new(Client_V1)
	if os.Getenv("TEST") == "1" {
		u, _ := url.Parse(os.Getenv("TEST_URL"))
		client.URL = u
	} else {
		u, _ := url.Parse("https://api.esa.io/v1")
		client.URL = u
	}
	client.TeamName = teamName
	client.HTTPClient = &http.Client{}
	client.Authorization = auth
	client.Logger = log.New(os.Stdout, "logger: ", log.Lshortfile)
	return client
}

func (c *Client_V1) newRequest(method string, spath string, body io.Reader) (*http.Request, error) {
	u := *c.URL
	u.Path = path.Join(c.URL.Path, spath)
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.Authorization)
	return req, nil
}

func (c *Client_V1) GetCategories() (*Categories, error) {
	spath := fmt.Sprintf("/teams/%v/categories", c.TeamName)
	req, _ := c.newRequest("GET", spath, nil)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var categories Categories
	if err := decodeBody(res, &categories); err != nil {
		return nil, err
	}
	return &categories, nil
}

func main() {
	auth := os.Getenv("ESA_AUTH")
	team := os.Getenv("ESA_TEAM")
	client := NewClient(auth, team)
	categories, _ := client.GetCategories()
	for _, category := range categories.Categories {
		for _, tree := range category.Tree() {
			fmt.Println(tree)
		}
	}
}
