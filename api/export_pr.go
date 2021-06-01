package api

import (
	"reflect"
	"strings"
)

func (issue *Issue) ExportData(fields []string) *map[string]interface{} {
	v := reflect.ValueOf(issue).Elem()
	data := map[string]interface{}{}

	for _, f := range fields {
		switch f {
		case "comments":
			data[f] = issue.Comments.Nodes
		case "assignees":
			data[f] = issue.Assignees.Nodes
		case "labels":
			data[f] = issue.Labels.Nodes
		case "projectCards":
			data[f] = issue.ProjectCards.Nodes
		default:
			sf := fieldByName(v, f)
			data[f] = sf.Interface()
		}
	}

	return &data
}

func (pr *PullRequest) ExportData(fields []string) *map[string]interface{} {
	v := reflect.ValueOf(pr).Elem()
	data := map[string]interface{}{}

	for _, f := range fields {
		switch f {
		case "headRepository":
			data[f] = pr.HeadRepository
		case "statusCheckRollup":
			if n := pr.StatusCheckRollup.Nodes; len(n) > 0 {
				data[f] = n[0].Commit.StatusCheckRollup.Contexts.Nodes
			} else {
				data[f] = nil
			}
		case "commits":
			commits := make([]interface{}, 0, len(pr.Commits.Nodes))
			for _, c := range pr.Commits.Nodes {
				commit := c.Commit
				authors := make([]interface{}, 0, len(commit.Authors.Nodes))
				for _, author := range commit.Authors.Nodes {
					authors = append(authors, map[string]interface{}{
						"name":  author.Name,
						"email": author.Email,
						"id":    author.User.ID,
						"login": author.User.Login,
					})
				}
				commits = append(commits, map[string]interface{}{
					"oid":             commit.OID,
					"messageHeadline": commit.MessageHeadline,
					"messageBody":     commit.MessageBody,
					"committedDate":   commit.CommittedDate,
					"authoredDate":    commit.AuthoredDate,
					"authors":         authors,
				})
			}
			data[f] = commits
		case "comments":
			data[f] = pr.Comments.Nodes
		case "assignees":
			data[f] = pr.Assignees.Nodes
		case "labels":
			data[f] = pr.Labels.Nodes
		case "projectCards":
			data[f] = pr.ProjectCards.Nodes
		case "reviews":
			data[f] = pr.Reviews.Nodes
		case "files":
			data[f] = pr.Files.Nodes
		case "reviewRequests":
			requests := make([]interface{}, 0, len(pr.ReviewRequests.Nodes))
			for _, req := range pr.ReviewRequests.Nodes {
				if req.RequestedReviewer.TypeName == "" {
					continue
				}
				requests = append(requests, req.RequestedReviewer)
			}
			data[f] = &requests
		default:
			sf := fieldByName(v, f)
			data[f] = sf.Interface()
		}
	}

	return &data
}

func fieldByName(v reflect.Value, field string) reflect.Value {
	return v.FieldByNameFunc(func(s string) bool {
		return strings.EqualFold(field, s)
	})
}