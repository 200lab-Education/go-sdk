package fbprovider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type UserInfo interface {
	GetID() string
	GetEmail() string
	GetPhone() string
	GetPhonePrefix() string
	// Get full phone number joined from prefix + national phone number
	GetFullPhone() string
}

type userInfo struct {
	Id    string `json:"id"`
	Phone *struct {
		Number         string `json:"number"`
		CountryPrefix  string `json:"country_prefix"`
		NationalNumber string `json:"national_number"`
	} `json:"phone"`
	Email *struct {
		Address string `json:"address"`
	} `json:"email"`
}

func (u *userInfo) GetID() string {
	return u.Id
}

func (u *userInfo) GetEmail() string {
	if u.Email == nil {
		return ""
	}
	return u.Email.Address
}

// Get phone number
func (u *userInfo) GetPhone() string {
	if u.Phone == nil {
		return ""
	}
	return u.Phone.NationalNumber
}

func (u *userInfo) GetPhonePrefix() string {
	if u.Phone == nil {
		return ""
	}
	return u.Phone.CountryPrefix
}

// Get full phone number joined from prefix + national phone number
func (u *userInfo) GetFullPhone() string {
	if u.Phone == nil {
		return ""
	}
	return fmt.Sprintf("%s%s", u.Phone.CountryPrefix, u.Phone.NationalNumber)
}

type accountKit struct {
	akURI string
}

// Create AccountKit URI provider
// akURI: https://graph.accountkit.com/v1.3
// accessToken: Facebook accessToken when user logged in
func NewAccountKitProvider(akURI string) *accountKit {
	return &accountKit{akURI: akURI}
}

func (ak *accountKit) GetUserInfo(ctx context.Context, token string) (UserInfo, error) {
	resp, err := http.Get(fmt.Sprintf("%s/me/?access_token=%s", ak.akURI, token))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode >= http.StatusMultipleChoices {
		return nil, errors.New(string(data))
	}

	u := &userInfo{}

	err = json.Unmarshal(data, &u)
	if err != nil {
		return nil, err
	}

	return u, nil
}
