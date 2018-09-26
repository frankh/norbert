package config

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

var query = `query {
  viewer {
    repositories(first: 100) {
      pageInfo {
        endCursor
        hasNextPage
      }
      edges {
        node {
          name

          history: object(expression: "HEAD") {
            ... on Commit {
              history(first: 1, path: ".norbert.yml") {
                edges {
                  node {
                    author {
                      name
                    }
                    committedDate
                    oid
                  }
                }
              }
            }
          }

          object(expression: "HEAD:.norbert.yml") {
            ... on Blob {
              text
            }
          }
        }
      }
    }
  }
}
`

type queryResult struct {
	Data *struct {
		Viewer struct {
			Repositories struct {
				PageInfo struct {
					EndCursor   string
					HasNextPage bool
				}

				Edges []struct {
					Node struct {
						Name string

						History struct {
							History struct {
								Edges []struct {
									Node struct {
										CommittedDate time.Time
										Oid           string
									}
								}
							}
						}

						Object *struct {
							Text string
						}
					}
				}
			}
		}
	}
	Errors []struct {
		Message string
	}
}

type queryRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type githubFetcher struct {
	client *http.Client
}

func NewGithubFetcher(oauthToken string, organisation string) Fetcher {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: oauthToken},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	return githubFetcher{httpClient}
}

func (f githubFetcher) Fetch() (config *Config, err error) {
	defer func() {
		if err != nil {
			log.Println("Error fetching:", err)
		}
	}()

	request, err := json.Marshal(queryRequest{
		Query:     query,
		Variables: nil,
	})

	if err != nil {
		return
	}

	resp, err := f.client.Post("https://api.github.com/graphql", "application/json", bytes.NewBuffer(request))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var result queryResult
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return
	}

	for _, repo := range result.Data.Viewer.Repositories.Edges {
		if repo.Node.Object != nil {
			return configFromYaml([]byte(repo.Node.Object.Text))
		}
	}
	config = &Config{}
	return
}

func init() {
}
