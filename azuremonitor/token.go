// To parse and unparse this JSON data, add this code to your project and do:
//
//    token, err := UnmarshalToken(bytes)
//    bytes, err = token.Marshal()

package azuremonitor

import (
	"encoding/json"
	"strconv"
	"time"
)

//UnmarshalToken parse string to Token structure
func UnmarshalToken(data []byte) (Token, error) {
	var r Token
	err := json.Unmarshal(data, &r)
	return r, err
}

//Marshal Parse Token Structure to string
func (r *Token) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

//Token define the token struct defined by azure active directory
type Token struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    string `json:"expires_in"`
	EXTExpiresIn string `json:"ext_expires_in"`
	ExpiresOn    string `json:"expires_on"`
	NotBefore    string `json:"not_before"`
	Resource     string `json:"resource"`
	AccessToken  string `json:"access_token"`
}

//IsExpired return true if the token is expired or false if is not.
func (r *Token) IsExpired() bool {
	var t = time.Now().UTC()
	i64, _ := strconv.ParseInt(r.ExpiresOn, 10, 32)
	var tokenDate = time.Unix(i64, 0).UTC()
	tokenDate.Add(-100)

	return !t.Before(tokenDate)

}
