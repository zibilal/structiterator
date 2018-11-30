package validator

import (
	"testing"
)

type User struct {
	Name    string `valid:"funcVal:Required,errorMessage:Please provide your name"`
	Address string `valid:"funcVal:Required"`
	Email   string `valid:"funcVal:Required;funcVal:Email"`
	Phone   string `valid:"funcVal:Required;funcVal:Phone"`
}

type Person struct {
	Name  string `valid:"funcVal:Required"`
	Email string `valid:"funcVal:Required;funcVal:Email"`
}

type Application struct {
	Id           uint   `json:"id" valid:"funcVal:Required"`
	Name         string `json:"application_name" valid:"funcVal:Required"`
	AppliedTime  string `json:"applied_time" valid:"funcVal:Required"`
	ApprovedTime string `json:"approved_time" valid:"funcVal:Required;funcVal:AfterDate,compareKey:AppliedTime"`
}

type ApplicationSecond struct {
	Id           uint   `json:"id" valid:"funcVal:Required"`
	Name         string `json:"application_name" valid:"funcVal:Required"`
	AppliedTime  string `json:"applied_time" valid:"funcVal:Required"`
	ApprovedTime string `json:"approved_time" valid:"funcVal:Required;funcVal:AfterDate,compareKey:applied_time"`
}

type App2 struct {
	Id             uint   `valid:"funcVal:Required"`
	Name           string `valid:"funcVal:Required"`
	Status         string
	ApprovalReason string `valid:"funcVal:CondRequired,compareKey:Status,compareValue:approved|rejected"`
}

type App3 struct {
	Id          uint   `valid:"funcVal:Required"`
	Name        string `valid:"funcVal:Required"`
	PhoneNumber string `valid:"funcVal:Required;funcVal:Match,format:^(62|0)([0-9]*)$"`
}

type Appl struct {
	Id          uint   `valid:"funcVal:Required"`
	Name        string `valid:"funcVal:Required"`
	DateJoin    string `valid:"funcVal:Required;funcVal:Date,format:mm/dd/yyyy,dateLayout:01/02/2006"`
	DateTesting string `valid:"funcVal:Required;funcVal:Date,format:mm/dd/yyyy,dateLayout:01/02/2006,errorMessage:Wrong date format, pls check your format"`
}

type Appll struct {
	Id         uint   `valid:"funcVal:Required"`
	Name       string `valid:"funcVal:Required"`
	RewardType string `valid:"funcVal:AcceptedValues,values:e-voucher|giftcard"`
}

func TestAppll(t *testing.T) {
	t.Log("\nTesting Appll struct")
	{
		appll := Appll{}
		appll.Name = "Testing"
		appll.RewardType = "The gifts"
		mapper := NewValidationMapper()
		validtr := NewValidStruct(mapper)
		errors := validtr.Valid(appll)
		t.Log("Errors", errors)
	}
}

func TestValidationWithErrorMap(t *testing.T) {
	t.Log("\nTesting Appll struct with error map")
	{
		appll := Appll{}
		appll.Name = "Testing"
		appll.RewardType = "The gifts"

		mapper := NewValidationMapper()
		validtr := NewValidStructWithMap(mapper, map[string]string{
			"Required": "Fields is required",
		})
		errors := validtr.Valid(appll)
		t.Log("Errors", errors)
	}
}

func TestApplStruct(t *testing.T) {
	mapper := NewValidationMapper()
	t.Log("\nTesting Appl struct")
	{
		appl := Appl{}
		validtr := NewValidStruct(mapper)
		errors := validtr.Valid(appl)
		t.Log("Errors", errors)
	}

	t.Log("\nTesting Appl struct")
	{
		appl := Appl{}
		appl.Id = 123
		appl.Name = "Test Appl Struct"
		validtr := NewValidStruct(mapper)
		errors := validtr.Valid(appl)
		t.Log("Errors", errors)
	}

	t.Log("\nTesting Appl struct, wrong date format")
	{
		appl := Appl{}
		appl.Id = 123
		appl.Name = "Test Appl Struct"
		appl.DateJoin = "29/11/2017"
		validtr := NewValidStruct(mapper)
		errors := validtr.Valid(appl)
		t.Log("Errors", errors)
	}

	t.Log("\nTesting Appl struct, wrong date format")
	{
		appl := Appl{}
		appl.Id = 123
		appl.Name = "Test Appl Struct"
		appl.DateTesting = "date 19/29/2017"
		validtr := NewValidStruct(mapper)
		errors := validtr.Valid(appl)
		t.Log("Errors", errors)
	}
}

func TestNewValidStruct(t *testing.T) {
	t.Log("\nTesting Using StructIterator:")
	{
		user := User{}
		mapper := NewValidationMapper()
		validtr := NewValidStruct(mapper)
		errors := validtr.Valid(user)
		t.Log("Errors", errors)
	}
}

func TestNewValidStruct2(t *testing.T) {
	mapper := NewValidationMapper()
	t.Log("\nTesting Using ValidStruc, with empty struct:")
	{
		person := Person{}
		validtr := NewValidStruct(mapper)
		errors := validtr.Valid(person)
		t.Log("Errors", errors)
	}

	t.Log("\nTesting unempty struct:")
	{
		person := Person{
			Name:  "Bilal Muhammad",
			Email: "Bilal.muhammad@exampl.com",
		}
		validtr := NewValidStruct(mapper)
		errors := validtr.Valid(person)
		t.Log("Errors2", errors)
	}
}

func TestNewValidStruct3(t *testing.T) {
	mapper := NewValidationMapper()
	t.Log("\nTesting for Application struct")
	{
		app := Application{}
		app.Id = 1
		app.Name = "AppName"
		app.AppliedTime = "09/20/2017"
		validtr := NewValidStruct(mapper)
		errors := validtr.Valid(app)
		t.Log("Errors3", errors)
	}

	t.Log("\nTesting for ApplicationSecond struct")
	{
		app := ApplicationSecond{}
		app.Id = 1
		app.AppliedTime = "09/20/2017"
		app.ApprovedTime = "09/19/2017"
		validtr := NewValidStruct(mapper)
		errors := validtr.Valid(app)
		t.Log("Errors3", errors)
	}
}

func TestNewValidStruct4(t *testing.T) {
	mapper := NewValidationMapper()
	t.Log("\nTesting for App2 struct")
	{
		app := App2{}
		app.Id = 2
		app.Name = "Bilal Muhammad"
		app.Status = "rejected"

		validtr := NewValidStruct(mapper)
		errors := validtr.Valid(app)
		t.Log("Errors4", errors)
	}
	t.Log("\nTesting for App3 struct")
	{
		app := App3{}
		app.Id = 3
		app.Name = "Bilal Muhammad"
		app.PhoneNumber = "+6281817800"

		validtr := NewValidStruct(mapper)
		errors := validtr.Valid(app)
		t.Log("Errors5", errors)
	}
}

func TestDataTag(t *testing.T) {
	t.Log("\nTesting data tag. Input 1\t")
	{
		input := "funcVal:Required;funcVal:Email;funcVal:Email"
		dataTags := []*dataTag{}
		dataTags = fetchDataTag(input, -1, dataTags)

		for _, tag := range dataTags {
			t.Logf("Tag %v\n", tag)
		}

		t.Log("ends")
	}

	t.Log("\nTesting data tag, Input2\n")
	{
		input2 := "funcVal:Required,errorMessage:Lagi Test Nih,key:saman_name"
		dataTags2 := []*dataTag{}
		dataTags2 = fetchDataTag(input2, -1, dataTags2)
		for _, tag := range dataTags2 {
			t.Logf("Tag %v\n", tag)
		}

		t.Log("ends")
	}

	t.Log("\nTesting data tag. Input3\n")
	{
		input3 := "funcVal:Required;funcVal:Email,key:email_address"
		dataTags3 := []*dataTag{}
		dataTags3 = fetchDataTag(input3, -1, dataTags3)
		for _, tag := range dataTags3 {
			t.Logf("Tag %v\n", tag)
		}

		t.Log("ends")
	}

	t.Log("\nTesting data tag.Input for compareKey and compareValue\n")
	{
		input4 := "funcVal:CondRequired,compareKey:update_status,compareValue:1"
		dataTag4 := []*dataTag{}
		dataTag4 = fetchDataTag(input4, -1, dataTag4)
		for _, tag := range dataTag4 {
			t.Logf("Tag %v\n", tag)
		}

		t.Log("ends")
	}
}
