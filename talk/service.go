package talk

import (
	"github.com/vecmezoni/gomeet/jira"
	"github.com/vecmezoni/gomeet/user"
)

func GetAllTalks(user *user.User, client *jira.AuthorizedClient) (*Talks, error) {
	result, err := client.GetTalks(`project = "R&D :: Meetups"`)

	if err != nil {
		return nil, err
	}

	return talksFromJira(result, user), nil
}

func Get(user *user.User, client *jira.AuthorizedClient, key string) (*Talk, error) {
	jiraTalk, err := client.GetTalk(key)

	if err != nil {
		return nil, err
	}

	result := talkFromJira(jiraTalk, user)

	return result, nil
}

func Like(user *user.User, client *jira.AuthorizedClient, key string) (*Talk, error) {
	err := client.Like(key)

	if err != nil {
		return nil, err
	}

	result, err := Get(user, client, key)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func Dislike(user *user.User, client *jira.AuthorizedClient, key string) (*Talk, error) {
	err := client.Dislike(key)

	if err != nil {
		return nil, err
	}

	result, err := Get(user, client, key)

	if err != nil {
		return nil, err
	}

	return result, nil
}
