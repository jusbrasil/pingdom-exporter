package pingdom

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
	    "users": [
        {
            "id": 12,
            "name": "John Doe",
            "paused": "NO",
            "access_level": "contact",
            "use_severity_levels": false,
            "sms": [
                {
                    "id": 352,
                    "severity": "HIGH",
                    "country_code": "1",
                    "number": "6095555555",
                    "provider": "nexmo"
                }
            ]
        },
        {
            "id": 234,
            "name": "Jane Doe",
            "paused": "NO",
            "access_level": "default",
            "use_severity_levels": false,
            "email": [
                {
                    "id": 314,
                    "severity": "HIGH",
                    "address": "test@test.com"
                }
            ]
		}
		]
		}`)
	})
	want := []UsersResponse{
		{
			Id:       12,
			Paused:   "NO",
			Username: "John Doe",
			Sms: []UserSmsResponse{
				{
					Id:          352,
					Severity:    "HIGH",
					CountryCode: "1",
					Number:      "6095555555",
					Provider:    "nexmo",
				},
			},
			Email: nil,
		},
		{
			Id:       234,
			Paused:   "NO",
			Username: "Jane Doe",
			Sms:      nil,
			Email: []UserEmailResponse{
				{
					Id:       314,
					Severity: "HIGH",
					Address:  "test@test.com",
				},
			},
		},
	}

	users, err := client.Users.List()
	assert.NoError(t, err)
	assert.Equal(t, want, users, "Users.List() should return correct result")
}

func TestUserService_Read(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
	    "users": [
        {
            "id": 12,
            "name": "John Doe",
            "paused": "NO",
            "access_level": "contact",
            "use_severity_levels": false,
            "sms": [
                {
                    "id": 352,
                    "severity": "HIGH",
                    "country_code": "1",
                    "number": "6095555555",
                    "provider": "nexmo"
                }
            ]
        },
        {
            "id": 234,
            "name": "Jane Doe",
            "paused": "NO",
            "access_level": "default",
            "use_severity_levels": false,
            "email": [
                {
                    "id": 314,
                    "severity": "HIGH",
                    "address": "test@example.com"
                }
            ]
		}
		]
		}`)
	})
	want := UsersResponse{
		Id:       12,
		Paused:   "NO",
		Username: "John Doe",
		Sms: []UserSmsResponse{
			{
				Id:          352,
				Severity:    "HIGH",
				CountryCode: "1",
				Number:      "6095555555",
				Provider:    "nexmo",
			},
		},
		Email: nil,
	}

	users, err := client.Users.Read(12)
	assert.NoError(t, err)
	assert.Equal(t, &want, users, "Users.Read(12) should return a user")
}

func TestUserService_Read_Failure(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
	    "users": [
        {
            "id": 12,
            "name": "John Doe",
            "paused": "NO",
            "access_level": "contact",
            "use_severity_levels": false,
            "sms": [
                {
                    "id": 352,
                    "severity": "HIGH",
                    "country_code": "1",
                    "number": "6095555555",
                    "provider": "nexmo"
                }
            ]
        },
        {
            "id": 234,
            "name": "Jane Doe",
            "paused": "NO",
            "access_level": "default",
            "use_severity_levels": false,
            "email": [
                {
                    "id": 314,
                    "severity": "HIGH",
                    "address": "test@example.com"
                }
            ]
		}
		]
		}`)
	})

	want := fmt.Errorf("UserId: 24 not found")

	_, err := client.Users.Read(24)
	assert.Equal(t, want, err, "Read with an invalid user id should return an error.")
}

func TestUserService_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"user": {
				"id": 23439
			}
		}`)
	})

	want := &UsersResponse{
		Id: 23439,
	}

	u := User{
		Username: "testUser",
	}

	user, err := client.Users.Create(&u)
	assert.NoError(t, err)
	assert.Equal(t, want, user, "Users.Create() should return correct result")
}

func TestUserService_CreateContact(t *testing.T) {
	setup()
	defer teardown()

	userID := 12941

	mux.HandleFunc("/users/"+strconv.Itoa(userID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"contact_target": {
				"id": 23439
			}
		}`)
	})

	want := &CreateUserContactResponse{
		Id: 23439,
	}

	c := Contact{
		Email:       "test@example.com",
		CountryCode: "1",
		Number:      "5559995555",
	}

	contact, err := client.Users.CreateContact(userID, c)
	assert.NoError(t, err)
	assert.Equal(t, want, contact, "Users.CreateContact() should return contact_target.id")
}

func TestUserService_Delete(t *testing.T) {
	setup()
	defer teardown()

	userID := 12941

	mux.HandleFunc("/users/"+strconv.Itoa(userID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{
			"message":"Deletion of user was successful!"
		}`)
	})

	want := &PingdomResponse{
		Message: "Deletion of user was successful!",
	}

	response, err := client.Users.Delete(userID)
	assert.NoError(t, err)
	assert.Equal(t, want, response, "Users.Delete() should return PingdomResponse with message")

}

func TestUserService_DeleteContact(t *testing.T) {
	setup()
	defer teardown()

	userID := 12941
	contactID := 87655

	mux.HandleFunc("/users/"+strconv.Itoa(userID)+"/"+strconv.Itoa(contactID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{
			"message":"Deletion of contact target successful"
		}`)
	})

	want := &PingdomResponse{
		Message: "Deletion of contact target successful",
	}

	response, err := client.Users.DeleteContact(userID, contactID)
	assert.NoError(t, err)
	assert.Equal(t, want, response, "Users.DeleteContact() should return PingdomResponse with message")

}

func TestUserService_Update(t *testing.T) {
	setup()
	defer teardown()

	userID := 12941
	user := User{
		Username: "updatedUsername",
	}

	mux.HandleFunc("/users/"+strconv.Itoa(userID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
			"message":"Modification of user was successful!"
		}`)
	})

	want := &PingdomResponse{
		Message: "Modification of user was successful!",
	}

	response, err := client.Users.Update(userID, &user)
	assert.NoError(t, err)
	assert.Equal(t, want, response, "Users.Update() should return PingdomResponse with message")

}

func TestUserService_UpdateContact(t *testing.T) {
	setup()
	defer teardown()

	userID := 12941
	contactID := 87655
	contact := Contact{
		Email: "test@example.com",
	}

	mux.HandleFunc("/users/"+strconv.Itoa(userID)+"/"+strconv.Itoa(contactID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
			"message":"Modification of contact target was successful!"
		}`)
	})

	want := &PingdomResponse{
		Message: "Modification of contact target was successful!",
	}

	response, err := client.Users.UpdateContact(userID, contactID, contact)
	assert.NoError(t, err)
	assert.Equal(t, want, response, "Users.UpdateContact() should return PingdomResponse with message")

}
