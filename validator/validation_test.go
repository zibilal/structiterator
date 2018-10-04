package validator

import (
	"testing"
	"time"
)

const (
	success = "\u2713"
	failed  = "\u2717"
)

func TestValidation_Email(t *testing.T) {
	validation := Validation{}

	t.Log("\nTesting Validation Email:")
	{
		err := validation.Email("samalekom", "email_address", "")
		t.Log("The error", err)

		err = validation.Email("kubilisme@example.com", "email_address2", "")
		t.Log("The error2", err)

		err = validation.Email(12345, "email_address", "")
		t.Log("The error3", err)
	}

	t.Log("\nTesting Validation Phone")
	{
		err := validation.Phone("79178877", "phone_number", "")
		t.Log("The error", err)

		err = validation.Phone("021791788777", "phone_number", "")
		t.Log("The error2", err)

		err = validation.Phone("62817180018", "cell_phone", "")
		t.Log("The error3", err)
	}
}

func TestValidation_AfterDate(t *testing.T) {
	validtn := Validation{}

	t.Log("\nTesting validation after date:")
	{
		app := Application{}
		app.Id = 1
		app.AppliedTime = "20/09/2017"
		app.ApprovedTime = "19/09/2017"
		err := validtn.AfterDate(app, "approved_time", "applied_time", "")
		t.Log("Error", err)
	}
}

func TestValidation_AcceptedValues(t *testing.T) {

	t.Log("\nTesting validation outside range, with error message")
	{
		validtn := Validation{}
		errorMessage := "Wrong range of values"
		err := validtn.AcceptedValues(120, "approved_value", "1-5", errorMessage)

		if err == nil {
			t.Logf("%s expected error Wrong range of values got nil", failed)
		} else {
			if err.Error() == errorMessage {
				t.Logf("%s expected error %s", success, err.Error())
			} else {
				t.Logf("%s expected error %s got %s", failed, errorMessage, err.Error())
			}
		}
	}

	t.Log("\nTesting validation outside range, without error message")
	{
		validtn := Validation{}
		errorMessage := "120 is outside of range 1 - 5"
		err := validtn.AcceptedValues(120, "approved_value", "1-5", "")

		if err == nil {
			t.Logf("%s expected error %s got nil", errorMessage, failed)
		} else {
			if err.Error() == errorMessage {
				t.Logf("%s expected error %s", success, err.Error())
			} else {
				t.Logf("%s expected error %s got %s", failed, errorMessage, err.Error())
			}
		}
	}

	t.Log("\nTesting validation inside range, without error message")
	{
		validtn := Validation{}
		err := validtn.AcceptedValues(4, "approved_value", "1-5", "")

		if err == nil {
			t.Logf("%s expected error nil", success)
		} else {
			t.Logf("%s expected error nil got %s", failed, err.Error())
		}
	}
}

func TestDate(t *testing.T) {
	t.Log("\nTesting parsing date:")
	{
		tim, err := time.Parse("01/02/2006", "09/20/2017")

		t.Log("The time", tim, "The error", err)

		value := "05/19/11"
		// Writing down the way the standard time would look like formatted our way
		layout := "01/02/06"
		tim2, _ := time.Parse(layout, value)
		t.Log(tim2)
	}
}

func TestUrl(t *testing.T) {
	t.Log("\nTesting parsing url:")
	{
		url1 := "abcdszzz"
		url2 := "http://exampl.com"
		url3 := "http://images.example.com/static/images/foo.jpg"

		validtn := Validation{}
		err := validtn.Url(url1, "url_satu", "url1 problem")

		if err != nil {
			t.Logf("%s expected error not nil, got %s", success, err.Error())
		} else {
			t.Fatalf("%s expected error not nil", failed)
		}

		err = validtn.Url(url2, "url_dua", "url2 problem")
		if err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)

		}

		err = validtn.Url(url3, "url_tiga", "url3 problem")
		if err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)

		}

		url4 := "http://www.google.com/test.png"
		err = validtn.Url(url4, "url_tiga", "")
		if err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)

		}

		url5 := "http://192.168.0.22/test.jpeg"
		err = validtn.Url(url5, "url_tiga", "url3 problem")
		if err != nil {
			t.Fatalf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)

		}
	}
}
