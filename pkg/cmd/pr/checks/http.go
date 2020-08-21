package checks

import (
	"fmt"
	"time"

	"github.com/cli/cli/api"
	"github.com/cli/cli/internal/ghrepo"
)

type checkRun struct {
	Name    string
	Status  string
	Elapsed time.Duration
}

type checkRunList struct {
	Passing   int
	Failing   int
	Pending   int
	checkRuns []checkRun
}

func checkRuns(client *api.Client, repo ghrepo.Interface, pr *api.PullRequest) (checkRunList, error) {
	list := checkRunList{}

	type response struct {
		Repository struct {
			PullRequest struct {
				Commits struct {
					Nodes []struct {
						Commit struct {
							CheckSuites struct {
								Nodes []struct {
									CheckRuns struct {
										Nodes []struct {
											Title      string
											Status     string
											Name       string
											Conclusion string
											DetailsUrl string
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	query := `
	query PullRequestChecks($owner: String!, $repo: String!, $pr_number: Int!) {
		repository(owner: $owner, name: $repo) {
			pullRequest(number: $pr_number) {
				commits(last: 1) {
				  nodes {
				  	commit {
				  		checkSuites(first:100) {
				  			nodes {
				  			  checkRuns(first:100) {
				  			  	nodes {
				  			  		title
				  			  		status
				  			  		conclusion
				  			  		detailsUrl
				  			  	}
				  			  }
				  			}
				  		}
				  	}
				  }
				}
			}
		}
	}
	`

	variables := map[string]interface{}{
		"owner":     repo.RepoOwner(),
		"repo":      repo.RepoName(),
		"pr_number": pr.Number,
	}

	var resp response
	err := client.GraphQL(repo.RepoHost(), query, variables, &resp)
	if err != nil {
		return list, err
	}

	suites := resp.Repository.PullRequest.Commits.Nodes[0].Commit.CheckSuites.Nodes

	for _, suite := range suites {
		for _, runNode := range suite.CheckRuns.Nodes {
			fmt.Printf("DEBUG %#v\n", runNode)
		}
	}

	return list, nil
}
