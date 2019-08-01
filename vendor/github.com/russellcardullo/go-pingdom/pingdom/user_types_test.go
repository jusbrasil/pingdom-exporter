package pingdom

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_PostParams(t *testing.T) {
	name := "testUsername"

	user := User{
		Username: name,
	}
	params := user.PostParams()
	want := map[string]string{
		"name": name,
	}

	assert.Equal(t, want, params, "User.PostParams() should return correct map")
}

func TestUser_ValidUser_Positive(t *testing.T) {
	name := "testUsername"
	user := User{
		Username: name,
	}

	err := user.ValidUser()

	assert.Equal(t, nil, err, "User.ValidUser() should return nil")
}

func TestUser_ValidUser_Negative(t *testing.T) {
	user := User{
		Username: "",
	}

	want := fmt.Errorf("Invalid value for `Username`.  Must contain non-empty string")

	err := user.ValidUser()

	assert.Equal(t, want, err, "User.ValidUser() should return error")
}

func TestContact_ValidContact_Positive1(t *testing.T) {
	contact := Contact{
		Email:       "test@example.com",
		CountryCode: "1",
		Number:      "5559995555",
	}

	err := contact.ValidContact()
	assert.Equal(t, nil, err, "contact.ValidContact() should return nil")
}

func TestContact_ValidContact_Positive2(t *testing.T) {
	contact := Contact{
		Email: "test@example.com",
	}

	err := contact.ValidContact()
	assert.Equal(t, nil, err, "contact.ValidContact() should return nil")
}

func TestContact_ValidContact_Positive3(t *testing.T) {
	contact := Contact{
		CountryCode: "1",
		Number:      "5559995555",
	}

	err := contact.ValidContact()
	assert.Equal(t, nil, err, "contact.ValidContact() should return nil")
}

func TestContact_ValidContact_Negative1(t *testing.T) {
	contact := Contact{
		CountryCode: "1",
	}

	want := fmt.Errorf("you must provide either an Email or a Phone Number to create a contact target")

	err := contact.ValidContact()

	assert.Equal(t, want, err, "contact.ValidContact() should return error")
}

func TestContact_ValidContact_Negative2(t *testing.T) {
	contact := Contact{
		Number: "5559995555",
	}

	want := fmt.Errorf("you must provide a Country Code if providing a phone number")

	err := contact.ValidContact()

	assert.Equal(t, want, err, "contact.ValidContact() should return error")
}

func TestContact_PostContactParams(t *testing.T) {
	email := "test@example.com"
	countrycode := "1"
	number := "5559995555"

	contact := Contact{
		Email:       email,
		CountryCode: countrycode,
		Number:      number,
	}
	params := contact.PostContactParams()
	want := map[string]string{
		"email":       email,
		"number":      number,
		"countrycode": countrycode,
	}

	assert.Equal(t, want, params, "Contact.PostContactParams() should return correct map")
}

func TestContact_PutContactParams(t *testing.T) {
	email := "test@example.com"
	countrycode := "1"
	number := "5559995555"

	contact := Contact{
		Email:       email,
		CountryCode: countrycode,
		Number:      number,
	}
	params := contact.PutContactParams()
	want := map[string]string{
		"email":       email,
		"number":      number,
		"countrycode": countrycode,
	}

	assert.Equal(t, want, params, "Contact.PutContactParams() should return correct map")
}

func TestUser_PutParams(t *testing.T) {
	name := "myname"
	primary := "YES"
	paused := "NO"

	user := User{
		Username: name,
		Primary:  primary,
		Paused:   paused,
	}

	params := user.PutParams()
	want := map[string]string{
		"name":    name,
		"primary": primary,
		"paused":  paused,
	}

	assert.Equal(t, want, params, "User.PutParams() should return correct map")
}
