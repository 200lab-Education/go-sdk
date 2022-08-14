package sdkcm

type Requester interface {
	OAuth
	User
}

type User interface {
	UserID() uint32
	GetSystemRole() string
	GetUser() interface{}
}

type OAuth interface {
	OAuthID() string
}

type currentUser struct {
	OAuth
	User
}

func CurrentUser(t OAuth, u User) *currentUser {
	return &currentUser{t, u}
}
