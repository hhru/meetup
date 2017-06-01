package jira

type Talks struct {
	Talks []Talk `json:"issues"`
}

type Talk struct {
	Key    string `json:"key"`
	Fields struct {
		Summary      string `json:"summary"`
		Video        string `json:"customfield_22511"`
		Presentation string `json:"customfield_22510"`
		Description  string `json:"description"`
		Duedate      string `json:"duedate"`
		Reporter     struct {
			DisplayName string `json:"displayName"`
			Name        string `json:"name"`
			AvatarUrls  struct {
				Big string `json:"48x48"`
			} `json:"avatarUrls"`
		} `json:"reporter"`
		Assignee struct {
			DisplayName string `json:"displayName"`
			Name        string `json:"name"`
			AvatarUrls  struct {
				Big string `json:"48x48"`
			} `json:"avatarUrls"`
		} `json:"assignee"`
		Votes struct {
			Votes    int  `json:"votes"`
			HasVoted bool `json:"hasVoted"`
		} `json:"votes"`
		Status struct {
			Name string `json:"name"`
		} `json:"status"`
		Resolution struct {
			Name string `json:"name"`
		} `json:"resolution"`
		Attachment []struct {
			FileName string `json:"filename"`
			Content  string `json:"content"`
		} `json:"attachment"`
	} `json:"fields"`
}

type User struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	AvatarUrls  struct {
		Big string `json:"48x48"`
	} `json:"avatarUrls"`
}
