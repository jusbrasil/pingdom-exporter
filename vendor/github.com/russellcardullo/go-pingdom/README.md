# go-pingdom #

[![Build Status](https://travis-ci.org/russellcardullo/go-pingdom.svg?branch=master)](https://travis-ci.org/russellcardullo/go-pingdom) [![Go Report Card](https://goreportcard.com/badge/github.com/russellcardullo/go-pingdom/pingdom)](https://goreportcard.com/report/github.com/russellcardullo/go-pingdom/pingdom) [![GoDoc](https://godoc.org/github.com/russellcardullo/go-pingdom/pingdom?status.svg)](https://godoc.org/github.com/russellcardullo/go-pingdom/pingdom)

go-pingdom is a Go client library for the Pingdom API.

This currently supports working with HTTP, ping checks, and TCP checks.

## Usage ##

### Client ###

Pingdom handles single-user and multi-user accounts differently.

Construct a new single-user Pingdom client:

```go
client, err := pingdom.NewClientWithConfig(pingdom.ClientConfig{
    Username: "pingdom_username",
    Password: "pingdom_password",
    APIKey: "pingdom_api_key",
})
```

Construct a multi-user Pingdom client:

```go
client, err := pingdom.NewClientWithConfig(pingdom.ClientConfig{
    Username: "pingdom_username",
    Password: "pingdom_password",
    APIKey: "pingdom_api_key",
    AccountEmail: "pingdom_account_email",
})
```

The `pingdom_account_email` variable is the email address of the owner of the multi-user account. This is passed in the `Account-Email` header to the Pingdom API.

Using a Pingdom client, you can access supported services.

You can override the timeout or other parameters by passing a custom http client:
```go
client, err := pingdom.NewClientWithConfig(pingdom.ClientConfig{
    Username: "pingdom_username",
    Password: "pingdom_password",
    APIKey: "pingdom_api_key",
    AccountEmail: "pingdom_account_email",
    HTTPClient: &http.Client{
        Timeout: time.Second * 10,
    },
})
```

### CheckService ###

This service manages pingdom Checks which are represented by the `Check` struct.
When creating or updating Checks you must specify at a minimum the `Name`, `Hostname`
and `Resolution`.  Other fields are optional but if not set will be given the zero
values for the underlying type.

More information on Checks from Pingdom: https://www.pingdom.com/features/api/documentation/#ResourceChecks

Get a list of all checks:

```go
checks, err := client.Checks.List()
fmt.Println("Checks:", checks) // [{ID Name} ...]
```

Create a new HTTP check:

```go
newCheck := pingdom.HttpCheck{Name: "Test Check", Hostname: "example.com", Resolution: 5}
check, err := client.Checks.Create(&newCheck)
fmt.Println("Created check:", check) // {ID, Name}
```

Create a new Ping check:
```go
newCheck := pingdom.PingCheck{Name: "Test Check", Hostname: "example.com", Resolution: 5}
check, err := client.Checks.Create(&newCheck)
fmt.Println("Created check:", check) // {ID, Name}
```

Create a new TCP check:
```go
newCheck := pingdom.TCPCheck{Name: "Test Check", Hostname: "example.com", Port: 25, StringToSend: "HELO foo.com", StringToExpect: "250 mail.test.com", Resolution: 5}
check, err := client.Checks.Create(&newCheck)
fmt.Println("Created check:", check) // {ID, Name}
```

Get details for a specific check:

```go
checkDetails, err := client.Checks.Read(12345)
```

For checks with detailed information, check the specific details in
the field `Type` (e.g. `checkDetails.Type.HTTP`).

Update a check:

```go
updatedCheck := pingdom.HttpCheck{Name: "Updated Check", Hostname: "example2.com", Resolution: 5}
msg, err := client.Checks.Update(12345, &updatedCheck)
```

Delete a check:

```go
msg, err := client.Checks.Delete(12345)
```

Create a check with basic alert notification to a user.

```go
newCheck := pingdom.HttpCheck{Name: "Test Check", Hostname: "example.com", Resolution: 5, SendNotificationWhenDown: 2, UserIds []int{12345}}
checkResponse, err := client.Checks.Create(&newCheck)
```

### MaintenanceService ###

This service manages pingdom Maintenances which are represented by the `Maintenance` struct.
When creating or updating Maintenances you must specify at a minimum the `Description`, `From`
and `To`.  Other fields are optional but if not set will be given the zero
values for the underlying type.

More information on Maintenances from Pingdom: https://www.pingdom.com/resources/api/2.1#ResourceMaintenance

Get a list of all maintenances:

```go
maintenances, err := client.Maintenances.List()
fmt.Println("Maintenances:", maintenances) // [{ID Description} ...]
```

Create a new Maintenance Window:

```go
m := pingdom.MaintenanceWindow{
    Description: "My Maintenance",
    From:        1,
    To:          1234567899,
}
maintenance, err := client.Maintenances.Create(&m)
fmt.Println("Created MaintenanceWindow:", maintenance) // {ID Description}
```

Get details for a specific maintenance:

```go
maintenance, err := client.Maintenances.Read(12345)
```

Update a maintenance: (Please note, that based on experience, you are allowed to modify only `Description`, `EffectiveTo` and `To`)

```go
updatedMaintenance := pingdom.MaintenanceWindow{
    Description: "My Maintenance",
    To:          1234567999,
}
msg, err := client.Maintenances.Update(12345, &updatedMaintenance)
```

Delete a maintenance:

Note: that only future maintenance window can be deleted. This means that both `To` and `From` should be in future.

```go
msg, err := client.Maintenances.Delete(12345)
```

After contacting Pingdom, the better approach would be to use update function and setting `To` and `EffectiveTo` to current time

```go
maintenance, _ := client.Maintenances.Read(12345)

m := pingdom.MaintenanceWindow{
    Description: maintenance.Description,
    From:        maintenance.From,
    To:          1,
    EffectiveTo: 1,
}

maintenanceUpdate, err := client.Maintenances.Update(12345, &m)
```

### ProbeService ###

This service gets pingdom Probes which are represented by the `Probes` struct.

More information on Probes from Pingdom: https://www.pingdom.com/resources/api/2.1#ResourceProbes
Several parameters are supported for filtering output. Please see them in Pingdom API documentation.

**NOTE:** Official documentation does not specify that `region` is returned for every probe entry, but it does and you can use it.

Get a list of all probes:

```go
params := make(map[string]string)

probes, err := client.Probes.List(params)
fmt.Println("Probes:", probes) // [{ID Name} ...]

for _, probe := range probes {
    fmt.Println("Probe region:", probe.Region)  // Probe region: EU
}
```

### TeamService ###

This service manages pingdom Teams which are represented by the `Team` struct.
When creating or updating Teams you must specify at a minimum the `Name`.
Other fields are optional but if not set will be given the zero
values for the underlying type.

More information on Teams from Pingdom: https://www.pingdom.com/resources/api/2.1#ResourceTeam

Get a list of all teams:

```go
teams, err := client.Teams.List()
fmt.Println("Teams:", teams) // [{ID Name Users} ...]
```

Create a new Team:

```go
t := pingdom.TeamData{
    Name: "Team",
}
team, err := client.Teams.Create(&t)
fmt.Println("Created Team:", team) // {ID Name Users}
```

Get details for a specific team:

```go
team, err := client.Teams.Read(12345)
```

Update a team:

```go
modifyTeam := pingdom.TeamData{
    Name:    "New Name"
    UserIDs: "123,678",
}
team, err := client.Teams.Update(12345, &modifyTeam)
```

Delete a team:

```go
team, err := client.Teams.Delete(12345)
```

### PublicReportService ###

This service manages the pingdom public report.  The public report is a
publicly-visible list of checks, and this API can be used to List those
checks, Publish checks to the public report or Withdrawl them.

More information on the Public Report from Pingdom: https://www.pingdom.com/resources/api/2.1/#ResourceReports.public

> Note: There is only one "public report", and it cannot be deleted.  To remove the public report, you must list all the checks and withdrawl them one-by-one.

Get a list of all published checks in the public report:

```go
checks, err := client.PublicReport.List()
fmt.Println("Checks:", checks) // [{ID Name ReportURL} ...]
```

Publish a check to the public report:

```go
_, err := client.PublicReport.PublishCheck(12345)
```

Withdrawl a check from the public report:

```go
_, err := client.PublicReport.WithdrawlCheck(12345)
```

### UserService ###

This service manages users and their contact information which is represented by the `User` struct. 
When creating or modifying users you must provide the `Username`. 
More information from Pingdom: https://www.pingdom.com/resources/api/2.1/#ResourceUsers

Get all users and contact info:

```go
users, err := client.Users.List()
fmt.Println(users)
```

Create a new user and contact:

```go
user := User{
    Username : "loginName",
    Paused : "NO",
} 
userId, err := client.Users.Create(user)
fmt.Println("New UserId: ", userId.Id)

contact := Contact{
    Number : "5555555555",
    CountryCode : "1",
    Provider : "Verizon",
}
contactId, err := client.Users.CreateContact(userId.Id, contact)
fmt.Println("New Contact Id: ", contactId.Id)
```

Update a user and contact

```go
userId := 1234
contactId := 90877

user := User{
    Username : "loginName",
    Paused : "NO",
} 
result, err := client.Users.Update(userId, user)
fmt.Println("result.Message)

contact := Contact{
    Number : "5555555555",
    CountryCode : "1",
    Provider : "Verizon",
}
result, err := client.Users.UpdateContact(userId, contactId, contact)
fmt.Println(result.Message)
```

Delete a user and contact (deleting a user will delete all contact info)

```go
userId := 1234
contactId := 90877

result, err := client.Users.DeleteContact(userId, contactId)
fmt.Println(result.Message)


result, err := client.Users.Delete(userId)
fmt.Println("result.Message)
```

## Development ##

### Acceptance Tests ###

You can run acceptance tests against the actual pingdom API to test any changes:
```
PINGDOM_USER=[username] \
  PINGDOM_PASSWORD=[password] \
  PINGDOM_API_KEY=[api key] make acceptance
```

Note that this will create actual resources in your Pingdom account.  The tests will make a best effort to clean up but these would
not be guaranteed on test failures depending on the nature of the failure.
