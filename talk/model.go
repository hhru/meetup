package talk

import (
	"github.com/vecmezoni/gomeet/jira"
	"github.com/vecmezoni/gomeet/user"
)

type Talks struct {
	Talks []*Talk `json:"talks"`
}

func talksFromJira(j *jira.Talks, user *user.User) *Talks {
	var talks []*Talk

	for _, talk := range j.Talks {
		talks = append(talks, talkFromJira(&talk, user))
	}

	return &Talks{talks}
}

type Talk struct {
	Key         string       `json:"key"`
	Author      author       `json:"author"`
	Status      string       `json:"status"`
	Summary     string       `json:"summary"`
	Description string       `json:"description"`
	Duedate     string       `json:"duedate"`
	Votes       int          `json:"votes"`
	HasVoted    bool         `json:"hasVoted"`
	CanVote     bool         `json:"canVote"`
	Attachments []attachment `json:"attachments"`
}

type attachment struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type author struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func talkFromJira(j *jira.Talk, user *user.User) *Talk {
	author := author{
		Name:   j.Fields.Assignee.DisplayName,
		Avatar: j.Fields.Assignee.AvatarUrls.Big,
	}

	attachments := []attachment{}

	for _, item := range j.Fields.Attachment {
		attachments = append(attachments, attachment{item.FileName, item.Content})
	}

	if j.Fields.Presentation != "" {
		attachments = append(attachments, attachment{"presentation", j.Fields.Presentation})
	}

	if j.Fields.Video != "" {
		attachments = append(attachments, attachment{"video", j.Fields.Video})
	}

	canVote := true

	if j.Fields.Reporter.Name == user.Name {
		canVote = false
	}

	if j.Fields.Resolution.Name != "" {
		canVote = false
	}

	return &Talk{
		Key:         j.Key,
		Author:      author,
		Summary:     j.Fields.Summary,
		Duedate:     j.Fields.Duedate,
		Description: j.Fields.Description,
		Votes:       j.Fields.Votes.Votes,
		HasVoted:    j.Fields.Votes.HasVoted,
		Status:      j.Fields.Status.Name,
		Attachments: attachments,
		CanVote:     canVote,
	}
}
