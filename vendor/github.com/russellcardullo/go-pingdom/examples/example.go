package main

import (
	"fmt"

	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/russellcardullo/go-pingdom/pingdom"
)

type credentials struct {
	User         string `json:"user"`
	Password     string `json:"password"`
	APIKey       string `json:"apikey"`
	AccountEmail string `json:"accountEmail"`
}

func getConfig() credentials {
	// Config Example
	// { "user" : "", "password" : "", "apikey" : "", "accountEmail" : "" }
	if len(os.Args) < 2 {
		fmt.Println("You must provide a configuration file.")
	}
	args := os.Args[1:] // get the config file from cli
	// Open our jsonFile
	jsonFile, err := os.Open(args[0])
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	fmt.Println("Successfully Opened: " + args[0])

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var config credentials

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &config)

	fmt.Println(config)

	return config
}

func userExamples() {
	config := getConfig()
	client := pingdom.NewMultiUserClient(config.User, config.Password, config.APIKey, config.AccountEmail)

	//Create User
	user := pingdom.User{
		Username: "exampleUser",
	}
	u, _ := client.Users.Create(&user)
	fmt.Println("User Id: " + strconv.Itoa(u.Id))

	// Create contact info
	contact := pingdom.Contact{
		Email: "test@example.com",
	}
	c, _ := client.Users.CreateContact(u.Id, contact)
	fmt.Println("Contact Id: " + strconv.Itoa(c.Id))

	//List all users and contacts
	users, _ := client.Users.List()
	fmt.Println("All users:", users)

	user.Username = "newExampleUser"
	uu, _ := client.Users.Update(u.Id, &user)
	fmt.Println(uu.Message)

	contact.Email = ""
	contact.Provider = "Nexmo"
	contact.Number = "5555555555"
	contact.CountryCode = "1"

	cc, _ := client.Users.UpdateContact(u.Id, c.Id, contact)
	fmt.Println(cc.Message)

	//Delete our example User Cpmtact
	rContact, _ := client.Users.DeleteContact(u.Id, c.Id)
	fmt.Println(rContact.Message)

	//Delete our example User
	rUser, _ := client.Users.Delete(u.Id)
	fmt.Println(rUser.Message)
}

func main() {
	client, err := pingdom.NewClientWithConfig(pingdom.ClientConfig{
		User:     "username",
		Password: "password",
		APIKey:   "api_key",
	})
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}

	// List all checks
	checks, _ := client.Checks.List()
	fmt.Println("All checks:", checks)

	// Create a new http check
	newCheck := pingdom.HttpCheck{Name: "Test Check", Hostname: "example.com", Resolution: 5}
	check, _ := client.Checks.Create(&newCheck)
	fmt.Println("Created check:", check) // {ID, Name}

	// Create a new ping check
	newPingCheck := pingdom.PingCheck{Name: "Test Ping", Hostname: "example.com", Resolution: 1}
	pingcheck, _ := client.Checks.Create(&newPingCheck)
	fmt.Println("Created check:", pingcheck) // {ID, Name}

	// Get details for a check
	details, _ := client.Checks.Read(check.ID)
	fmt.Println("Details:", details)

	// Update a check
	updatedCheck := pingdom.HttpCheck{Name: "Updated Check", Hostname: "example2.com", Resolution: 5}
	upMsg, _ := client.Checks.Update(check.ID, &updatedCheck)
	fmt.Println("Modified check, message:", upMsg)

	// Delete a check
	delMsg, _ := client.Checks.Delete(check.ID)
	fmt.Println("Deleted check, message:", delMsg)

	userExamples()
}
