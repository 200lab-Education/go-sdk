/*
 * @author           Viet Tran <viettranx@gmail.com>
 * @copyright        2019 Viet Tran <viettranx@gmail.com>
 * @license          Apache-2.0
 */

package oauthclient

import (
	"github.com/200Lab-Education/go-sdk/sdkcm"
	"net/url"
	"strings"
	"time"
)

type Gender string

type TokenIntrospect struct {
	Active   bool   `json:"active"`
	ClientId string `json:"client_id"`
	Scope    string `json:"scope"`
	Exp      uint32 `json:"exp"`
	Iat      uint32 `json:"iat"`
	Sub      string `json:"sub"`
	Username string `json:"username"`
	Email    string `json:"email"`
	UserId   string `json:"user_id"`
}

type Token struct {
	AccessToken         string `json:"access_token"`
	RefreshToken        string `json:"refresh_token"`
	OAuthId             string `json:"oauth_id"`
	Expiry              int    `json:"expires_in"`
	IsNew               bool   `json:"is_new"`
	HasUsernamePassword bool   `json:"has_username_password"`
}

type TokenResponse struct {
	*Token `json:",inline"`
	Error  string `json:"error_hint"`
}

func (t TokenIntrospect) OAuthID() string {
	return t.UserId
}

type AccountType string

const (
	// account logged in from other source: FB/Google/Github
	AccTypeExternal = "external"
	// account logged with username and password
	AccTypeInternal = "internal"
	// account external but has set username and password
	AccTypeBoth = "both"
)

type OAuthUser struct {
	Id          string          `json:"id" gorm:"id"`
	Username    string          `json:"username" gorm:"username"`
	Email       string          `json:"email"`
	PhonePrefix string          `json:"phone_prefix" bson:"phone_prefix,omitempty" gorm:"phone_prefix"`
	Phone       string          `json:"phone" bson:"phone,omitempty"`
	AccountType AccountType     `json:"account_type" bson:"account_type" gorm:"account_type"`
	FBId        string          `json:"fb_id" bson:"fb_id" gorm:"fb_id"`
	AKId        string          `json:"ak_id" bson:"account_kit_id" gorm:"column:account_kit_id"`
	AppleId     string          `json:"apple_id" bson:"apple_id" gorm:"apple_id"`
	ClientId    string          `json:"client_id" bson:"client_id"`
	Dob         *sdkcm.JSONDate `json:"dob,omitempty" form:"-" gorm:"dob" time_format:"2006-01-02"`
	Status      *int            `json:"-" form:"-" gorm:"status"`
	Gender      *Gender         `json:"gender,omitempty" form:"gender" gorm:"gender"`
	Address     *string         `json:"address,omitempty" form:"address" gorm:"address"`
}

type OAuthUserCreate struct {
	Username    *string      `json:"username" gorm:"username"`
	Password    *string      `json:"-"`
	Email       *string      `json:"email"`
	PhonePrefix *string      `json:"phone_prefix" bson:"phone_prefix,omitempty" gorm:"phone_prefix"`
	Phone       *string      `json:"phone" bson:"phone,omitempty"`
	AccountType *AccountType `json:"account_type" bson:"account_type" gorm:"account_type"`
	FBId        *string      `json:"fb_id" bson:"fb_id" gorm:"fb_id"`
	AKId        *string      `json:"ak_id" bson:"account_kit_id" gorm:"column:account_kit_id"`
	AppleId     *string      `json:"apple_id" bson:"apple_id" gorm:"apple_id"`
	ClientId    *string      `json:"client_id" bson:"client_id"`
}

type OAuthUserUpdate struct {
	Username    *string
	FirstName   *string `json:"first_name" form:"first_name" gorm:"first_name"`
	LastName    *string `json:"last_name" form:"last_name" gorm:"last_name"`
	Email       *string `json:"email,omitempty" form:"email" gorm:"email"`
	PhonePrefix *string `json:"phone_prefix,omitempty" form:"phone_prefix" gorm:"phone_prefix"`
	Phone       *string `json:"phone,omitempty" form:"phone" gorm:"phone"`
	Gender      *Gender `json:"gender,omitempty" form:"gender" gorm:"gender"`
	Address     *string `json:"address,omitempty" form:"address" gorm:"address"`

	Dob       *sdkcm.JSONDate `json:"dob,omitempty" form:"-" gorm:"dob" time_format:"2006-01-02"`
	DobString *string         `json:"dob" form:"dob" gorm:"-"`
	Status    *int            `json:"-" form:"-" gorm:"status"`

	Password             *string `json:"password"`
	PasswordConfirmation *string `json:"password_confirmation" gorm:"-"`

	FBId        *string      `json:"fb_id" bson:"fb_id" gorm:"fb_id"`
	AKId        *string      `json:"ak_id" bson:"account_kit_id" gorm:"column:account_kit_id"`
	AppleId     *string      `json:"apple_id" bson:"apple_id" gorm:"apple_id"`
	AccountType *AccountType `json:"account_type" bson:"account_type" gorm:"account_type"`
}

func (OAuthUserUpdate) TableName() string {
	return "users"
}

func (u *OAuthUserUpdate) ProcessData() error {
	if phone := u.Phone; phone != nil && strings.HasPrefix(*phone, "0") {
		p := strings.TrimPrefix(*phone, "0")
		u.Phone = &p
	}

	if u.DobString != nil {
		t, err := time.Parse("2006-01-02", *u.DobString)
		if err != nil {
			return err
		}

		dob := sdkcm.JSONDate(t)
		u.Dob = &dob
	}

	return nil
}

type OAuthUserFilter struct {
	UserId      *string `json:"id" bson:"-" gorm:"-"`
	Username    *string `json:"username" gorm:"username"`
	Email       *string `json:"email"`
	PhonePrefix *string `json:"phone_prefix" bson:"phone_prefix,omitempty" gorm:"phone_prefix"`
	Phone       *string `json:"phone" bson:"phone,omitempty"`
	FBId        *string `json:"fb_id" bson:"fb_id" gorm:"fb_id"`
	AKId        *string `json:"ak_id" bson:"account_kit_id" gorm:"column:account_kit_id"`
	AppleId     *string `json:"apple_id" bson:"apple_id" gorm:"apple_id"`
	ClientId    *string `json:"client_id" bson:"client_id"`
}

type OAuthUserLogin struct {
	Username *string `json:"username" gorm:"username"`
	Email    *string `json:"email" gorm:"email"`
	Phone    *string `json:"phone" gorm:"phone"`
	Password *string `json:"password" gorm:"phone"`
	OTPCode  *string `json:"otp_code" gorm:"otp_code"`
	ClientId *string `json:"client_id" gorm:"client"`
}

func (oul *OAuthUserLogin) URLValues() url.Values {
	payload := url.Values{}

	if oul.Username != nil {
		payload.Add("username", *oul.Username)
	}

	if oul.Password != nil {
		payload.Add("password", *oul.Password)
	}

	if oul.Phone != nil {
		payload.Add("phone", *oul.Phone)
	}

	if oul.Email != nil {
		payload.Add("email", *oul.Email)
	}

	if oul.OTPCode != nil {
		payload.Add("otp_code", *oul.OTPCode)
	}

	if oul.ClientId != nil {
		payload.Add("client_id", *oul.ClientId)
	}

	return payload
}
