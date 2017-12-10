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
Each *funcVal* is separated by *;* (semicolon)

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
Code above will result *errors* contained one error message ```DefaultEmail has invalid format value```

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
Code above will result *errors* contains one error message Phone 

### funcVal: Date
funcVal: Date is used to validate date data. The default validator for a date is a golang date layout of the following:
```
`01/02/2006` // this read as mm/dd/yyyy date formatting
```
This regex will force the date value to have format as the following example:
>12/29/2017  // read as december 29th 2017

Example of usage is the following:
```
type User struct {
	Name         string `valid:"funcVal:Required"`
	Address      string
	DefaultEmail string `valid:"funcVal:Email"`
	MobilePhone  string `valid:"funcVal:Phone"`
	BirthDate    string `valid:"funcVal:Date"`
}
...
	user := User{}
	user.Name = "User Test"
	user.DefaultEmail = "user.test@example.com"
	user.MobilePhone = "021812777"
	user.BirthDate = "29/12/2017"
	validtr := validator.NewValidStruct()
	errors := validtr.Valid(user)
```
Code above will result *errors* contains one error message BirthDate is expected of format mm/dd/yyyy

### funcVal: Match
funcVal: Match is used to match value with the regular expression in format attribute.

Example of usage is the following:

```
type User struct {
	Name         string `valid:"funcVal:Required"`
	Address      string
	DefaultEmail string `valid:"funcVal:Email"`
	MobilePhone  string `valid:"funcVal:Phone"`
	VoucherCode  string `valid:"funcVal:Match,format:[V|E]V-[1-9][0-9]{2}-[1-9][0-9]{4}"`
}
...
	user := User{}
	user.Name = "User Test"
	user.DefaultEmail = "user.test@example.com"
	user.MobilePhone = "0812777"
	user.VoucherCode = "VV-112-01234"
	validtr := validator.NewValidStruct()
	errors := validtr.Valid(user)
```
Code above will result *errors* contains one error message VourcherCode has invalid format value.

### funcVal: CondRequired
funcVal: CondRequired is a conditional required validation logic. This is for field that 
is become required if some other field have particular value.

*CondRequired* has another required attribute to be set. Those are *compareKey* and *compareValue*.
Compare key is attribute that define field on the object its value will affecting the field.
And *compareValue* is set of values that will affecting the field requirement status.

From previous example, if we want to make VoucherCode is required only if IsVoucher send field value
are "reg" or "elect". We can add tag, as shown below:

```
type User struct {
	Name         string `valid:"funcVal:Required"`
	Address      string
	DefaultEmail string `valid:"funcVal:Email"`
	MobilePhone  string `valid:"funcVal:Phone"`
	Type         string
	VoucherCode  string `valid:"funcVal:CondRequired,compareKey:Type,compareValue:e-voucher|voucher;funcVal:Match,format:[V|E]V-[1-9][0-9]{2}-[1-9][0-9]{4}"`
}
```

Valid tag on *VoucherCode* means if *Type* has value of *e-voucher* or *voucher* , *VoucherCode*
field becomes required. So, the following codes, will be generate an error *VoucherCode is required*

```
	user := User{}
	user.Name = "User Test"
	user.DefaultEmail = "user.test@example.com"
	user.MobilePhone = "0812777"
	user.Type = "e-voucher"
```

We need to also set *VoucherCode* to avoid the error.
```
	user.Name = "User Test"
	user.DefaultEmail = "user.test@example.com"
	user.MobilePhone = "0812777"
	user.Type = "e-voucher"
	user.VoucherCode = "VV-112-11234"
```

### funcVal: AfterDate
funcVal: AfterDate compares between two dates, and validate that one that should come after the other.
This validation function, has other required attribute, *compareKey*. This will compare this field date value, with value in the *compareKey*.
If the value on the *compareKey* is bigger on value on the field, this will generate error.

For example, if we want to force *DateApproved* always comes after *DateApplied*, we can write the user struect as follows:

```
type User struct {
	Name         string `valid:"funcVal:Required"`
	Address      string
	DefaultEmail string `valid:"funcVal:Email"`
	MobilePhone  string `valid:"funcVal:Phone"`
	Type         string
	VoucherCode  string `valid:"funcVal:CondRequired,compareKey:Type,compareValue:e-voucher|voucher;funcVal:Match,format:[V|E]V-[1-9][0-9]{2}-[1-9][0-9]{4}"`
	AppliedDate  string `valid:"funcVal:Required"`
	ApprovedDate string `valid:"funcVal:AfterDate,compareKey:AppliedDate"`
}
```

If we set both *AppliedDate* and *ApprovedDate* and *AppliedDate* value is after *ApprovedDate*, error will displayed.

```
	user := User{}
	user.Name = "User Test"
	user.DefaultEmail = "user.test@example.com"
	user.MobilePhone = "0812777"
	user.Type = "e-voucher"
	user.VoucherCode = "VV-112-11234"
	user.AppliedDate = "12/10/2017"
	user.ApprovedDate = "12/09/2017"
	validtr := validator.NewValidStruct()
	errors := validtr.Valid(user)
```
Code above will result *errors* contains one error message *invalid ApprovedDate should be after AppliedDate*.

### funcVal: AcceptedValues
funcVal: AcceptedValues is used to force a field to have some set of values or range of values. If we want to define a range
of values we separate value with **<->** operator. If we want to define set of values, we separate each values with **|** operator.
We define this accepted values inside *values* attribute.

For example, if want to force Age only have value between 16 to 25, and Role should have value *platinum*, *gold*, *silver*, or *fee*.
We can define the struct as follows:

```
type User struct {
	Name         string `valid:"funcVal:Required"`
	Address      string
	DefaultEmail string `valid:"funcVal:Email"`
	MobilePhone  string `valid:"funcVal:Phone"`
	Type         string
	VoucherCode  string `valid:"funcVal:CondRequired,compareKey:Type,compareValue:e-voucher|voucher;funcVal:Match,format:[V|E]V-[1-9][0-9]{2}-[1-9][0-9]{4}"`
	AppliedDate  string `valid:"funcVal:Required"`
	ApprovedDate string `valid:"funcVal:AfterDate,compareKey:AppliedDate"`
	Age          uint   `valid:"funcVal:AcceptedValues,values:16<->25"`
	Role         string `valid:"funcVal:AcceptedValues,values:platinum|gold|silver|free"`
}

...
	user := User{}
	user.Name = "User Test"
	user.DefaultEmail = "user.test@example.com"
	user.MobilePhone = "0812777"
	user.Type = "e-voucher"
	user.VoucherCode = "VV-112-11234"
	user.AppliedDate = "12/10/2017"
	user.ApprovedDate = "12/09/2017"
	validtr := validator.NewValidStruct()
	errors := validtr.Valid(user)
```
Code above will result *errors* contains two error message *15 is outside of range 16 - 25* and *wrong value ordinary, accepted values platinum|gold|silver|free*.

## Customizing error message

To customize error message, we can add *errorMessage* attribute inside struct definition;

```
type Appl struct {
	Id uint `valid:"funcVal:Required,errorMessage:id is required, pls provide the id"`
	DateTesting string `valid:"funcVal:Required;funcVal:Date,format:mm/dd/yyyy,dateLayout:01/02/2006,errorMessage:Wrong date format, pls check your format"`
}
```

We can also create error message map, and inject this map when instantiate the ValidStruct object.
```
    appll := Appll{}
	appll.Name="Testing"
	appll.RewardType = "The gifts"
	validtr := NewValidStructWithMap(map[string]string{
		"Required": "Fields is required",
	})
```
The map must have key the same with *funcVal* name, like the above example

## Customizing field name with json tag

If json tag is included in struct definition, it will treated as field name