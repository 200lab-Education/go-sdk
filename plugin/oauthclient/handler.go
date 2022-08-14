/*
 * @author           Viet Tran <viettranx@gmail.com>
 * @copyright        2019 Viet Tran <viettranx@gmail.com>
 * @license          Apache-2.0
 */

package oauthclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/200Lab-Education/go-sdk/sdkcm"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Return Token object when login with username and password
func (o *oauth) PasswordCredentialsToken(ctx context.Context, username, password string) (*Token, error) {
	var t TokenResponse

	req, _ := http.NewRequest("POST", o.clientConf.TokenURL, strings.NewReader(
		url.Values{
			"grant_type": {"password"},
			"username":   {username},
			"password":   {password},
			"scope":      o.clientConf.Scopes,
		}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(o.clientConf.ClientID, o.clientConf.ClientSecret)

	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req.WithContext(ctx))

	if err != nil {
		return nil, sdkcm.ErrInvalidRequest(err)
	}

	defer res.Body.Close()
	out, _ := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(out, &t); err != nil {
		return nil, sdkcm.ErrInvalidRequest(err)
	}

	if res.StatusCode != 200 {
		return nil, sdkcm.NewAppErr(errors.New(t.Error), res.StatusCode, t.Error).WithCode("wrong_username_password")
	}

	t.Token.HasUsernamePassword = true
	t.Token.IsNew = false

	return t.Token, nil
}

// Introspect return access token, refresh token, expired time and its data
func (o *oauth) Introspect(ctx context.Context, token string) (*TokenIntrospect, error) {
	var ti TokenIntrospect

	out, err := o.call(
		ctx,
		strings.Replace(o.clientConf.TokenURL, "token", "introspect", -1),
		url.Values{"token": []string{token}, "scope": []string{}},
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(out, &ti); err != nil {
		return nil, err
	}

	return &ti, nil
}

func (o *oauth) FindUserById(ctx context.Context, uid string) (*OAuthUser, error) {
	out, err := o.call(
		ctx,
		strings.Replace(
			o.clientConf.TokenURL,
			"token",
			fmt.Sprintf("users/%s", uid),
			-1,
		), nil)

	if err != nil {
		return nil, err
	}

	var user OAuthUser
	if err := json.Unmarshal(out, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (o *oauth) FindUser(ctx context.Context, filter *OAuthUserFilter) (*OAuthUser, error) {
	payload := url.Values{}

	if v := filter.Username; v != nil {
		payload.Add("username", *v)
	}

	if v := filter.Email; v != nil {
		payload.Add("email", *v)
	}

	if v := filter.FBId; v != nil {
		payload.Add("fb_id", *v)
	}

	if v := filter.AppleId; v != nil {
		payload.Add("apple_id", *v)
	}

	if v := filter.Phone; v != nil {
		payload.Add("phone", *v)
	}

	if v := filter.PhonePrefix; v != nil {
		payload.Add("phone_prefix", *v)
	}

	out, err := o.call(
		ctx,
		strings.Replace(
			o.clientConf.TokenURL,
			"token",
			fmt.Sprintf("find-user"),
			-1,
		), payload)

	if err != nil {
		return nil, err
	}

	data := struct {
		Code int       `json:"code"`
		User OAuthUser `json:"data"`
	}{}

	if err := json.Unmarshal(out, &data); err != nil {
		return nil, err
	}

	return &data.User, nil
}

func (o *oauth) CreateUser(ctx context.Context, user *OAuthUserCreate) (*Token, error) {
	var t Token

	payload := url.Values{}
	if user.Username != nil {
		payload.Add("username", *user.Username)
	}

	if user.Password != nil {
		payload.Add("password", *user.Password)
	}

	if user.Email != nil {
		payload.Add("email", *user.Email)
	}

	if user.PhonePrefix != nil {
		payload.Add("phone_prefix", *user.PhonePrefix)
	}

	if user.Phone != nil {
		payload.Add("phone", *user.Phone)
	}

	if user.ClientId != nil {
		payload.Add("client_id", *user.ClientId)
	}

	out, err := o.call(
		ctx,
		strings.Replace(
			o.clientConf.TokenURL,
			"token",
			"users",
			-1,
		), payload)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(out, &t); err != nil {
		return nil, err
	}

	return &t, nil
}

func (o *oauth) CreateUserWithEmail(ctx context.Context, email string) (*Token, error) {
	var t Token

	out, err := o.call(
		ctx,
		strings.Replace(
			o.clientConf.TokenURL,
			"token",
			"users?type=gmail",
			-1,
		), url.Values{
			"email": []string{email},
		})

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(out, &t); err != nil {
		return nil, err
	}

	return &t, nil
}

func (o *oauth) CreateUserWithPhone(ctx context.Context, phone string) (*Token, error) {
	var t Token

	out, err := o.call(
		ctx,
		strings.Replace(
			o.clientConf.TokenURL,
			"token",
			"users?type=phone",
			-1,
		), url.Values{
			"phone": []string{phone},
		})

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(out, &t); err != nil {
		return nil, err
	}

	return &t, nil
}

func (o *oauth) CreateUserWithFacebook(ctx context.Context, fbId, email string) (*Token, error) {
	var t Token

	out, err := o.call(
		ctx,
		strings.Replace(
			o.clientConf.TokenURL,
			"token",
			"users?type=facebook",
			-1,
		), url.Values{
			"fb_id": []string{fbId},
			"email": []string{email},
		})

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(out, &t); err != nil {
		return nil, err
	}

	return &t, nil
}

func (o *oauth) CreateUserWithAccountKit(ctx context.Context, akId, email, prefix, phone string) (*Token, error) {
	var t Token

	out, err := o.call(
		ctx,
		strings.Replace(
			o.clientConf.TokenURL,
			"token",
			"users?type=account-kit",
			-1,
		), url.Values{
			"ak_id":        []string{akId},
			"email":        []string{email},
			"phone_prefix": []string{prefix},
			"phone":        []string{phone},
		})

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(out, &t); err != nil {
		return nil, err
	}

	return &t, nil
}

func (o *oauth) CreateUserWithApple(ctx context.Context, appleId, email string) (*Token, error) {
	var t Token

	out, err := o.call(
		ctx,
		strings.Replace(
			o.clientConf.TokenURL,
			"token",
			"users?type=apple",
			-1,
		), url.Values{
			"apple_id": []string{appleId},
			"email":    []string{email},
		})

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(out, &t); err != nil {
		return nil, err
	}

	return &t, nil
}

func (o *oauth) UpdateUser(ctx context.Context, uid string, update *OAuthUserUpdate) error {
	payload := url.Values{
		"user_id": {uid},
	}

	if update.Username != nil {
		payload.Add("username", *update.Username)
	}
	if update.FirstName != nil {
		payload.Add("first_name", *update.FirstName)
	}
	if update.LastName != nil {
		payload.Add("last_name", *update.LastName)
	}
	if update.Gender != nil {
		payload.Add("gender", string(*update.Gender))
	}
	if update.Address != nil {
		payload.Add("address", *update.Address)
	}
	if update.Password != nil {
		payload.Add("password", *update.Password)
	}
	if update.Email != nil {
		payload.Add("email", *update.Email)
	}
	if update.PhonePrefix != nil {
		payload.Add("phone_prefix", *update.PhonePrefix)
	}
	if update.Phone != nil {
		payload.Add("phone", *update.Phone)
	}
	if update.Password != nil {
		payload.Add("password", *update.Password)
	}
	if update.PasswordConfirmation != nil {
		payload.Add("password_confirmation", string(*update.PasswordConfirmation))
	}
	if update.DobString != nil {
		payload.Add("dob", *update.DobString)
	}
	if update.FBId != nil {
		payload.Add("fb_id", *update.FBId)
	}
	if update.AKId != nil {
		payload.Add("ak_id", *update.AKId)
	}
	if update.AppleId != nil {
		payload.Add("apple_id", *update.FBId)
	}
	if update.AccountType != nil {
		payload.Add("account_type", string(*update.AccountType))
	}

	_, err := o.call(
		ctx,
		strings.Replace(
			o.clientConf.TokenURL,
			"token",
			fmt.Sprintf("users/%s/update", uid),
			-1,
		), payload)

	if err != nil {
		return err
	}

	return nil
}

func (o *oauth) ChangePassword(ctx context.Context, userId, oldPass, newPass string) error {
	_, err := o.call(
		ctx,
		strings.Replace(
			o.clientConf.TokenURL,
			"token",
			fmt.Sprintf("users/%s/change-password", userId),
			-1,
		), url.Values{
			"old_password": []string{oldPass},
			"new_password": []string{newPass},
		})

	if err != nil {
		return err
	}

	return nil
}

func (o *oauth) SetUsernamePassword(ctx context.Context, userId, username, password string) error {
	_, err := o.call(ctx,
		strings.Replace(
			o.clientConf.TokenURL,
			"token",
			fmt.Sprintf("users/%s/set-username-password", userId),
			-1,
		), url.Values{
			"username": []string{username},
			"password": []string{password},
		})

	if err != nil {
		return err
	}

	return nil
}

func (o *oauth) RevokeToken(ctx context.Context, token string) error {
	return nil
}

func (o *oauth) RefreshToken(ctx context.Context, refreshToken string) (*Token, error) {
	var t TokenResponse

	req, _ := http.NewRequest("POST", o.clientConf.TokenURL, strings.NewReader(
		url.Values{
			"grant_type":    {"refresh_token"},
			"refresh_token": {refreshToken},
			"scope":         o.clientConf.Scopes,
		}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(o.clientConf.ClientID, o.clientConf.ClientSecret)

	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req.WithContext(ctx))

	if err != nil {
		return nil, sdkcm.ErrInvalidRequest(err)
	}

	defer res.Body.Close()
	out, _ := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(out, &t); err != nil {
		return nil, sdkcm.ErrInvalidRequest(err)
	}

	if res.StatusCode != 200 {
		return nil, sdkcm.NewAppErr(errors.New(t.Error), res.StatusCode, t.Error)
	}

	return t.Token, nil
}

func (o *oauth) DeleteUser(ctx context.Context, userId string) error {
	_, err := o.call(ctx,
		strings.Replace(
			o.clientConf.TokenURL,
			"token",
			fmt.Sprintf("users/%s", userId),
			-1,
		), url.Values{})

	if err != nil {
		return err
	}

	return nil
}

func (o *oauth) GetUser(ctx context.Context, userId string) error {
	_, err := o.call(ctx,
		strings.Replace(
			o.clientConf.TokenURL,
			"token",
			fmt.Sprintf("users/%s", userId),
			-1,
		), url.Values{})

	if err != nil {
		return err
	}

	return nil
}

func (o *oauth) LoginOtherCredential(ctx context.Context, userLogin *OAuthUserLogin) (*Token, error) {
	payload := userLogin.URLValues()

	out, err := o.call(ctx,
		strings.Replace(
			o.clientConf.TokenURL,
			"token",
			"login",
			-1,
		), payload)

	if err != nil {
		return nil, err
	}

	var t Token

	if err := json.Unmarshal(out, &t); err != nil {
		return nil, err
	}

	return &t, nil
}

func (o *oauth) call(ctx context.Context, url string, params url.Values) ([]byte, error) {

	req, _ := http.NewRequest("POST", url, strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := o.client.Do(req.WithContext(ctx))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	out, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode >= 300 {
		var appErr sdkcm.AppError
		if err := json.Unmarshal(out, &appErr); err != nil {
			return nil, err
		}

		return nil, appErr
	}

	return out, nil
}
