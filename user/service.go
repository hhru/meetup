package user

import "github.com/hhru/meetup/jira"

func GetMyself(client *jira.AuthorizedClient) (*User, error) {
	myself, err := client.GetMyself()

	if err != nil {
		return nil, err
	}

	return fromJira(myself), nil
}
