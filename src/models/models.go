package models

type APIErrResponse struct {
	Status      string `json:"status"`
	Description string `json:"description,omitempty"`
}

type AccessTokenHeader struct {
	Sub          string `json:"sub"`
	TokenUse     string `json:"token_use"`
	Scope        string `json:"scope"`
	Auth_time    int    `json:"auth_time"`
	Iss          string `json:"iss"`
	Exp          int    `json:"exp"`
	Iat          int    `json:"iat"`
	Version      int    `json:"version"`
	Jti          string `json:"jti"`
	ClientId     string `json:"client_id"`
	UserName     string `json:"username"`
	Email        string `json:"email"`
	UserId       string `json:"custom:user_id"`
	IsAdmin      bool   `json:"custom:admin"`
	IsMaintainer bool   `json:"custom:maintainer"`
}

type AccessToken struct {
	Access_token string `json:"access_token"`
	Scope        string `json:"scope"`
	Token_type   string `json:"token_type"`
}

/*
Based on userdetails json struct of IDP the IDPUser model is changed.
For Example:
GitHub:
type IDPUser struct{
	Sub                int    `json:"id"`
	Name               string `json:"name"`
	Preferred_UserName string `json:"login"`
	Profile            string `json:"html_url"`
	Picture            string `json:"avatar_url"`
	Website            string `json:"blog"`
	Updated_At         string `json:"updated_at"`
}

Razer:
type IDPUser struct {
	Sub                string `json:"sub,omitempty"`
	Name               string `json:"name,omitempty"`
	Email              string `json:"email,omitempty"`
}
 */
type IDPUser struct {
	Sub                interface{}
	Id                 interface{}
	Name               string `json:"name,omitempty"`
	Preferred_UserName string `json:"login,omitempty"`
	Profile            string `json:"html_url,omitempty"`
	Picture            string `json:"avatar_url,omitempty"`
	Website            string `json:"blog,omitempty"`
	Updated_At         string `json:"updated_at,omitempty"`
	Email              string `json:"email,omitempty"`
}

type IDPUserEmailDetails struct {
	Email          string `json:"email"`
	Primary        bool   `json:"primary"`
	Email_Verified bool   `json:"verified"`
	Visibility     string `json:"visibility"`
}

type State struct {
	State     string `json:"state,omitempty"`
	Key       string `json:"key"`
	ExpiresOn int64  `json:"expires_on,omitempty"`
}

type TenantContext struct {
	CustomerId         string
	UserPoolId         string
	CognitoAppClientId string
}

type UserDetails struct {
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
}

type OpenIDUser struct {
	Sub                interface{}  `json:"sub"`
	Name               string `json:"username"`
	Preferred_UserName string `json:"login"`
	Updated_At         string `json:"updated_at"`
	Email              string `json:"email"`
	Email_Verified     bool   `json:"email_verified"`
	UserId             string `json:"custom:user_id"`
	Picture            string `json:"picture,omitempty"`
	IDPName            string `json:"name"`
}

type RainmakerOauth2Urls struct {
	RainmakerOauth2AuthorizeUrl string `json:"RainmakerOauth2AuthorizeUrl,omitempty"`
	RainmakerOauth2TokenUrl     string `json:"RainmakerOauth2TokenUrl,omitempty"`
	RainmakerOauth2UserinfoUrl  string `json:"RainmakerOauth2UserinfoUrl,omitempty"`
	RainmakerOauth2EmailUrl     string `json:"RainmakerOauth2EmailUrl,omitempty"`
}

type APIResponse struct {
	Status      string `json:"status"`
	Description string `json:"description,omitempty"`
}
