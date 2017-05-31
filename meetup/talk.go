package meetup

import (
  "github.com/vecmezoni/gomeet/jira"
  "time"
)

type Talks struct {
	Talks []*Talk `json:"talks"`
}

func (t Talks) Len() int {
  return len(t.Talks)
}

func (t Talks) Swap(i, j int) {
  t.Talks[i], t.Talks[j] = t.Talks[j], t.Talks[i]
}

type ByDate struct {
  Talks
}

func (s ByDate) Less(i, j int) bool {
  first, _ := time.Parse("2006-01-02", s.Talks.Talks[i].Duedate)
  second, _ := time.Parse("2006-01-02", s.Talks.Talks[j].Duedate)
  return first.Before(second)
}

func TalksFromJira(j *jira.Talks) *Talks {
	var talks []*Talk

	for _, talk := range j.Talks {
		talks = append(talks, TalkFromJira(&talk))
	}

	return &Talks{talks}
}

type Talk struct {
	Key string `json:"key"`
  Author author `json:"author"`
	Status string `json:"status"`
	Summary string `json:"summary"`
  Description string `json:"description"`
	Duedate string `json:"duedate"`
	Votes int `json:"votes"`
	HasVoted bool `json:"hasVoted"`
  Attachments []attachment `json:"attachments"`
}

type attachment struct {
  Name string `json:"name"`
  URL string `json:"url"`
}

type author struct {
  Name string `json:"name"`
  Avatar string `json:"avatar"`
}

func TalkFromJira(j *jira.Talk) *Talk {
  author := author{
    Name: j.Fields.Assignee.DisplayName,
    Avatar: j.Fields.Assignee.AvatarUrls.Original,
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

	return &Talk{
		Key: j.Key,
    Author: author,
		Summary: j.Fields.Summary,
		Duedate: j.Fields.Duedate,
    Description: j.Fields.Description,
		Votes: j.Fields.Votes.Votes,
		HasVoted: j.Fields.Votes.HasVoted,
		Status: j.Fields.Status.Name,
    Attachments: attachments,
	}
}

type User struct {
  Name string `json:"name"`
  DisplayName string `json:"display_name"`
  Avatar string `json:"avatar"`
}
