# structiterator

## Validator

Validator is a library to provide a way to put validation logic for struct as tag. This tag can be seen as a metadata for a struct.

Validator is using tag valid to inject the validation logic. Like the following example:
```
type App struct {
	Id             uint   `valid:"funcVal:Required"`
	Name           string `valid:"funcVal:Required"`
	Status         string
	ApprovalReason string `valid:"funcVal:CondRequired,compareKey:Status,compareValue:approved|rejected"`
}
```

For actually executing the validation logic, put the following lines in your code:

1. First you need to instantiate the validator object.

```
validtr := NewValidStruct()
```

2. Then execute ```validtr.Valid``` to run the intended validation logic.

```
app := App{}
app.Id = 2
app.Name = "First Test"
app.Status = "rejected"
errors := validtr.Valid(app)
```

The Valid function gives a list of errors according to validation logic not meet.
##Validation Function
Validator currently consists of the following validation logic. This validation logic is named ```funcVal```

###funcVal: Required

funcVal: Required will force the object field to be filled with a value. If the field is empty, ```Valid```
function will generate and error.

Example of usage:

```
package main

import (
	"fmt"
	"github.com/zibilal/structiterator/validator"
)

type User struct {
	Name    string `valid:"funcVal:Required"`
	Address string
	Email   string
	Phone   string
}

func main() {
	user := User{}
	validtr := validator.NewValidStruct()
	errors := validtr.Valid(user)
	fmt.Println("Errors", errors)
}
```

If code above is executed, this will display ```Errors [Name is required]```.

### funcVal: Email
funcVal: Email is used to validate email. Validator default email regular expression is
```
`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
```
This regex will force the email to have format as the following example:
>user.test@example.com

Example usage is the following:
```
type User struct {
	Name         string `valid:"funcVal:Required"`
	Address      string
	DefaultEmail string `valid:"funcVal:Email"`
	Phone        string
}
...

	user := User{}
	user.Name = "User Test"
	user.DefaultEmail = "user.test@example" // invalid email address format
	validtr := validator.NewValidStruct()
	errors := validtr.Valid(user)
```
Code above will result errors contained one error message ```DefaultEmail has invalid format value```

### funcVal: Phone
funcVal: Phone is used to validate phone number. Validator default phone regular expression is the following:
```
`^([62]|[0])[0-9]$`
```
This regex will force the phone number to have format as the following example:
> 081718188 or 6281718188

Example of usage is the following:
```
type User struct {
	Name         string `valid:"funcVal:Required"`
	Address      string
	DefaultEmail string `valid:"funcVal:Email"`
	MobilePhone  string `valid:"funcVal:Phone"`
}
...
	user := User{}
	user.Name = "User Test"
	user.DefaultEmail = "user.test@example.com"
	user.MobilePhone = "2812777"
	validtr := validator.NewValidStruct()
	errors := validtr.Valid(user)
```
Code above will result errors containe one error message Phone 