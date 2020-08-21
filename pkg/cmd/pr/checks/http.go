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
	path := fmt.Sprintf("repos/%s/%s/commits/%s/check-runs",
		repo.RepoOwner(), repo.RepoName(), pr.Commits.Nodes[0].Commit.Oid)
	var response struct {
		CheckRuns []struct {
			Name        string
			Status      string
			Conclusion  string
			StartedAt   time.Time `json:"started_at"`
			CompletedAt time.Time `json:"completed_at"`
			HtmlUrl     string    `json:"html_url"`
		} `json:"check_runs"`
	}

	err := client.REST(repo.RepoHost(), "GET", path, nil, &response)
	if err != nil {
		return list, err
	}

	for _, checkRun := range response.CheckRuns {
		elapsed := checkRun.CompletedAt.Sub(checkRun.StartedAt)
		fmt.Printf("%s %s %s %s\n", checkRun.Name, checkRun.Status, checkRun.Conclusion, elapsed)
	}

	return list, nil
}
