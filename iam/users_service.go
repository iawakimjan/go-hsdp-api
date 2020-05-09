package iam

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Jeffail/gabs"
	validator "github.com/go-playground/validator/v10"
)

const (
	userAPIVersion = "2"
)

// GetUserOptions describes search criteria for looking up users
type GetUserOptions struct {
	ID             *string `url:"_id,omitempty"`
	OrganizationID *string `url:"Id,omitempty"`
	Name           *string `url:"name,omitempty"`
	LoginID        *string `url:"loginId,omitempty"`
	GroupID        *string `url:"groupId,omitempty"`
	PageSize       *string `url:"pageSize,omitempty"`
	PageNumber     *string `url:"pageNumber,omitempty"`
}

// UsersService provides operations on IAM User resources
type UsersService struct {
	client *Client

	validate *validator.Validate
}

// Parameters holds parameters
type Parameters struct {
	ResourceType string  `json:"resourceType"`
	Parameter    []Param `json:"parameter"`
}

// ChangeLoginIDRequest
type ChangeLoginIDRequest struct {
	LoginID string `json:"loginId"`
}

// Param describes a resource
type Param struct {
	Name     string   `json:"name"`
	Resource Resource `json:"resource"`
}

// Resource holds a resource
type Resource struct {
	LoginID          string `json:"loginId,omitempty"`
	ConfirmationCode string `json:"confirmationCode,omitempty"`
	OldPassword      string `json:"oldPassword,omitempty"`
	NewPassword      string `json:"newPassword,omitempty"`
	Context          string `json:"context,omitempty"`
}

// UserList holds a paginated lists of users
type UserList struct {
	Users       []Person
	PageNumber  int
	PageSize    int
	HasNextPage bool
}

// CreateUser creates a new IAM user.
func (u *UsersService) CreateUser(person Person) (bool, *Response, error) {
	if err := u.validate.Struct(person); err != nil {
		return false, nil, err
	}
	req, err := u.client.NewRequest(IDM, "POST", "authorize/identity/User", &person, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", userAPIVersion)

	var bundleResponse interface{}

	doFunc := u.client.Do
	if person.ManagingOrganization == "" { // Self registration
		doFunc = u.client.DoSigned
	}
	resp, err := doFunc(req, &bundleResponse)

	if err != nil {
		return false, resp, err
	}
	ok := resp != nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated)
	return ok, resp, err
}

// DeleteUser deletes the  IAM user.
func (u *UsersService) DeleteUser(person Person) (bool, *Response, error) {
	req, err := u.client.NewRequest(IDM, "DELETE", "authorize/identity/User/"+person.ID, nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", "1")

	var bundleResponse interface{}

	doFunc := u.client.DoSigned
	if !u.client.validSigner() {
		doFunc = u.client.Do
	}
	resp, err := doFunc(req, &bundleResponse)

	if err != nil {
		return false, resp, err
	}
	ok := resp != nil && (resp.StatusCode == http.StatusNoContent)
	return ok, resp, err
}

// RecoverPassword triggers the recovery flow for the given user
//
// Deprecated: Support end date is 1 Augustus 2020
func (u *UsersService) RecoverPassword(loginID string) (bool, *Response, error) {
	body := &Parameters{
		ResourceType: "Parameters",
		Parameter: []Param{
			{
				Name: "recoverPassword",
				Resource: Resource{
					LoginID: loginID,
				},
			},
		},
	}
	return u.userActionV(body, "$recover-password", "1")
}

// ChangeLoginID changes the loginID
func (u *UsersService) ChangeLoginID(user Person, newLoginID string) (bool, *Response, error) {
	body := &ChangeLoginIDRequest{
		LoginID: newLoginID,
	}
	req, err := u.client.NewRequest(IDM, "POST", "authorize/identity/User/"+user.ID+"/$change-loginid", body, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", userAPIVersion)

	var bundleResponse interface{}
	doFunc := u.client.DoSigned
	if !u.client.validSigner() {
		doFunc = u.client.Do
	}
	resp, err := doFunc(req, &bundleResponse)
	if err != nil {
		return false, resp, err
	}
	ok := resp != nil && (resp.StatusCode == http.StatusNoContent)
	return ok, resp, err
}

// ResendActivation re-sends an activation email to the given user
func (u *UsersService) ResendActivation(loginID string) (bool, *Response, error) {
	body := &Parameters{
		ResourceType: "Parameters",
		Parameter: []Param{
			{
				Name: "resendOTP",
				Resource: Resource{
					LoginID: loginID,
				},
			},
		},
	}
	return u.userActionV(body, "$resend-activation", "2")
}

func (u *UsersService) userActionV(body *Parameters, action, apiVersion string) (bool, *Response, error) {
	req, err := u.client.NewRequest(IDM, "POST", "authorize/identity/User/"+action, body, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", apiVersion)

	var bundleResponse interface{}

	resp, err := u.client.DoSigned(req, &bundleResponse)
	if err != nil {
		return false, resp, err
	}
	ok := resp != nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted)
	return ok, resp, err
}

// SetPassword sets the password of a user given a correct confirmation code
func (u *UsersService) SetPassword(loginID, confirmationCode, newPassword, context string) (bool, *Response, error) {
	body := &Parameters{
		ResourceType: "Parameters",
		Parameter: []Param{
			{
				Name: "setPassword",
				Resource: Resource{
					LoginID:          loginID,
					ConfirmationCode: confirmationCode,
					NewPassword:      newPassword,
					Context:          context,
				},
			},
		},
	}
	return u.userActionV(body, "$set-password", "2")
}

// ChangePassword changes the password. The current pasword must be provided as well.
func (u *UsersService) ChangePassword(loginID, oldPassword, newPassword string) (bool, *Response, error) {
	body := &Parameters{
		ResourceType: "Parameters",
		Parameter: []Param{
			{
				Name: "changePassword",
				Resource: Resource{
					LoginID:     loginID,
					OldPassword: oldPassword,
					NewPassword: newPassword,
				},
			},
		},
	}
	return u.userActionV(body, "$change-password", "1")
}

// GetUsers looks up users by search criteria specified in GetUserOptions
func (u *UsersService) GetUsers(opts *GetUserOptions, options ...OptionFunc) (*UserList, *Response, error) {
	req, err := u.client.NewRequest(IDM, "GET", "security/users", opts, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("api-version", userAPIVersion)
	var bundleResponse struct {
		Exchange struct {
			Users []struct {
				UserUUID string `json:"userUUID"`
			}
			NextPageExists bool `json:"nextPageExists"`
		}
		ResponseCode    string `json:"responseCode"`
		ResponseMessage string `json:"responseMessage"`
	}

	resp, err := u.client.DoSigned(req, &bundleResponse)
	if err != nil {
		return nil, resp, err
	}
	var list UserList

	list.HasNextPage = bundleResponse.Exchange.NextPageExists
	list.Users = make([]Person, len(bundleResponse.Exchange.Users))
	for i, u := range bundleResponse.Exchange.Users {
		list.Users[i] = Person{ID: u.UserUUID}
	}

	return &list, resp, err
}

// GetUserByID looks up a user by UUID
func (u *UsersService) GetUserByID(uuid string) (*Person, *Response, error) {
	req, _ := u.client.NewRequest(IDM, "GET", "security/users/"+uuid, nil, nil)
	var user interface{}

	resp, err := u.client.Do(req, &user)
	if err != nil {
		return nil, resp, err
	}
	m, err := json.Marshal(user)
	if err != nil {
		return nil, resp, fmt.Errorf("error parsing json response: %w", err)
	}

	jsonParsed, err := gabs.ParseJSON(m)
	if err != nil {
		return nil, resp, fmt.Errorf("Eror decoding JSON: %w", err)
	}
	if err = checkResponseCode200(jsonParsed); err != nil {
		return nil, resp, &UserError{User: uuid, Err: err}
	}
	email, ok := jsonParsed.Path("exchange.loginId").Data().(string)
	if !ok {
		return nil, resp, fmt.Errorf("Invalid response")
	}
	r := jsonParsed.Path("exchange.profile")
	first := r.Path("givenName").Data().(string)
	last := r.Path("familyName").Data().(string)
	disabled := r.Path("disabled").Data().(bool)
	// TODO use Profile here
	var foundUser Person
	foundUser.Name.Family = last
	foundUser.Name.Given = first
	foundUser.Disabled = disabled
	foundUser.Telecom = append(foundUser.Telecom, TelecomEntry{
		System: "email",
		Value:  email,
	})
	return &foundUser, resp, nil
}

func checkResponseCode200(json *gabs.Container) error {
	if statusCode, ok := json.Path("responseCode").Data().(string); !ok || statusCode != "200" {
		responseMessage, ok := json.Path("responseMessage").Data().(string)
		if !ok {
			return fmt.Errorf("Unknown error")
		}
		switch responseMessage {
		case "Unauthorized Access":
			return ErrNotAuthorized
		}
		return fmt.Errorf(responseMessage)
	}
	return nil
}

// GetUserIDByLoginID looks up the UUID of a user by LoginID (email address)
func (u *UsersService) GetUserIDByLoginID(loginID string) (string, *Response, error) {
	req, _ := u.client.NewRequest(IDM, "GET", "security/users?loginId="+loginID, nil, nil)
	var d interface{}

	resp, err := u.client.Do(req, &d)
	if err != nil {
		return "", resp, err
	}
	m, err := json.Marshal(d)
	if err != nil {
		return "", resp, fmt.Errorf("error parsing json response: %w", err)
	}
	jsonParsed, err := gabs.ParseJSON(m)
	if err != nil {
		return "", resp, fmt.Errorf("Eror decoding JSON: %w", err)
	}
	if err = checkResponseCode200(jsonParsed); err != nil {
		return "", resp, &UserError{User: loginID, Err: err}
	}

	r := jsonParsed.Path("exchange.users").Index(0)
	userUUID, ok := r.Path("userUUID").Data().(string)
	if !ok {
		return "", resp, fmt.Errorf("lookup failed")
	}
	return userUUID, resp, nil

}

// SetMFA activate Multi-Factor-Authentication for the given UUID. See also SetMFAByLoginID.
func (u *UsersService) SetMFA(userID string, activate bool) (bool, *Response, error) {
	activateString := "true"
	if !activate {
		activateString = "false"
	}
	body := &struct {
		Activate string `json:"activate"`
	}{activateString}
	req, err := u.client.NewRequest(IDM, "POST", "authorize/identity/User/"+userID+"/$mfa", body, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", userAPIVersion)

	var bundleResponse interface{}

	resp, _ := u.client.Do(req, &bundleResponse)
	ok := resp != nil && (resp.StatusCode == http.StatusAccepted)
	return ok, resp, nil
}

// Unlock unlocks a user account with the given UserID
func (u *UsersService) Unlock(userID string) (bool, *Response, error) {
	req, err := u.client.NewRequest(IDM, "POST", "authorize/identity/User/"+userID+"/$unlock", nil, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("api-version", userAPIVersion)

	var bundleResponse interface{}

	resp, _ := u.client.Do(req, &bundleResponse)
	ok := resp != nil && (resp.StatusCode == http.StatusNoContent)
	return ok, resp, nil
}

// SetMFAByLoginID enabled Multi-Factor-Authentication for the given user. Only OrgAdmins can do this.
func (u *UsersService) SetMFAByLoginID(loginID string, activate bool) (bool, *Response, error) {
	userUUID, _, err := u.GetUserIDByLoginID(loginID)
	if err != nil {
		return false, nil, err
	}
	return u.SetMFA(userUUID, activate)
}
