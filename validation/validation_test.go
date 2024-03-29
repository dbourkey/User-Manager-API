package validation

import (
	"strings"
	"testing"
)

var valid_users = map[string]string{
	"bobman12": `{"username": "bobman12", "email": "bob@bobmail.com"}`,
	"jcdenton": `{"username": "jcdenton", "email": "jc@unatco.org"}`,
	"herpderp": `{"username": "herpderp", "email": "herp@derp.io"}`,
}

var invalid_users = map[string]string{
	"missingusername": `{"email": "missing@username.com"}`,
	"missingemail":    `{"username": "missingemail"}`,
}

var invalid_json = []string{
	`{"Json}`,
	`{"key": val}`,
	`{key: "val"}`,
	`{"key:" "val"}`,
	`{"key": {"subkey": "subval}`,
	``,
}

func Test_ValidateUser_invalid_json(t *testing.T) {
	for _, input := range invalid_json {
		_, err := ValidateUser(input)

		if err != nil {
			t.Logf("Invalid json %s was not rejected", input)
		}
	}
}

func Test_ValidateUser_UserMissingField(t *testing.T) {
	// Missing fullname
	user_missing_fullname := `{"username": "billy", "email": "Bob@bobmail.bob", "address": {"name": "Bob", "Line 1": "44 Bobstreet", "region": "Bobville", "country": "Bobland"}}`

	expected_err := "Fullname is a required field"
	returnval, actual_err := ValidateUser(user_missing_fullname)

	if actual_err.Error() != expected_err {
		t.Logf("Expected: %s\n Got: %s\n", expected_err, actual_err)
		t.Fail()
	}
	if returnval != "" {
		t.Logf("Return value should be empty string, Got: %s\n", returnval)
		t.Fail()
	}

	// Missing username
	user_missing_username := `{"fullname": "billy bobson", "email": "Bob@bobmail.bob", "address": {"name": "Bob", "Line 1": "44 Bobstreet", "region": "Bobville", "country": "Bobland"}}`

	expected_err = "Username is a required field"
	returnval, actual_err = ValidateUser(user_missing_username)

	if actual_err.Error() != expected_err {
		t.Logf("Expected: %s\n Got: %s\n", expected_err, actual_err)
		t.Fail()
	}
	if returnval != "" {
		t.Logf("Return value should be empty string, Got: %s\n", returnval)
		t.Fail()
	}

	// Missing email
	user_missing_email := `{"fullname": "billy bobson", "username": "bobman2000", "address": {"name": "Bob", "Line 1": "44 Bobstreet", "region": "Bobville", "country": "Bobland"}}`

	expected_err = "Email is a required field"
	returnval, actual_err = ValidateUser(user_missing_email)

	if actual_err.Error() != expected_err {
		t.Logf("Expected: %s\n Got: %s\n", expected_err, actual_err)
		t.Fail()
	}
	if returnval != "" {
		t.Logf("Return value should be empty string, Got: %s\n", returnval)
		t.Fail()
	}

}

func Test_ValidateUser_ReturnsLowercaseUsername(t *testing.T) {
	valid_users := []string{
		`{"username": "bIlLy2000", "fullname": "Billy Billy", "email": "Bob@bobmail.bob", "address": {"name": "Bob", "Line 1": "44 Bobstreet", "region": "Bobville", "country": "Bobland"}}`,
		`{"username": "BIlLy2000", "fullname": "Billy Billy", "email": "Bob@hotmail.com", "address": {"name": "Bob", "Line 1": "44 Bobstreet", "region": "Bobville", "country": "Bobland"}}`,
		`{"username": "BILLY2000", "fullname": "Billy Billy", "email": "Bob@gmail.nl", "address": {"name": "Null", "Line 1": "44 Edgevoid", "region": "Bobville", "country": "Bobland"}}`,
		`{"username": "bILLy2000", "fullname": "Billy Billy", "email": "Bob@edgcase.ie", "address": {"name": "Freddy", "Line 1": "Null", "region": "Bobville", "country": "Bobland"}}`,
	}

	for _, valid_user := range valid_users {
		returned_username, err := ValidateUser(valid_user)
		if returned_username != "billy2000" {
			t.Logf("Failed to return lowercase username: %s\n", returned_username)
			t.Fail()
		}
		if err != nil {
			t.Logf("Error %s, returned for valid user:\n%s\n", err, valid_user)
			t.Fail()
		}
	}
}

func Test_validateAddress(t *testing.T) {
	valid_addresses := []address{
		address{
			Name:     "Mr One",
			Line1:    "22 Woodroad",
			Line2:    "Ballyville",
			Region:   "Cork",
			Postcode: "11111",
			Country:  "Ireland",
		},
		address{
			Name:    "Mr Two",
			Line1:   "16 Fallsroad",
			Line2:   "Derry",
			Region:  "Belfast",
			Country: "Northern Ireland",
		},
		address{
			Name:    "The Occupant",
			Line1:   "22 Slumville",
			Region:  "Limerick",
			Country: "Something",
		},
		address{
			Name:     "Lord O'Fancy",
			Line1:    "Fancyland Manor",
			Line2:    "Use all lines",
			Line3:    "allfields",
			Region:   "Someregion",
			Postcode: "12123123",
			Country:  "Wakanda",
		},
		address{
			Name:    "名称",
			Line1:   "Somewhere",
			Region:  "عنوان",
			Country: "薛대한민국",
		},
	}

	for _, address := range valid_addresses {
		err := validateAddress(address)
		if err != nil {
			t.Logf("Expected nil, got: %s", err)
			t.Fail()
		}
	}

	missing_name := address{
		Line1:   "44 address blah blah",
		Region:  "Dublin 15",
		Country: "Ireland",
	}
	expected_err := "Address Name is a required field"
	err := validateAddress(missing_name)
	if err.Error() != expected_err {
		t.Logf("Expected: %s\nGot: %s\n", expected_err, err)
		t.Fail()
	}

}

func Test_validateFullname(t *testing.T) {
	expected_err := "Fullname is less than 3 characters"
	names_too_short := []string{"", "a", "bb"}

	for _, name := range names_too_short {
		err := validateFullname(name)
		if err.Error() != expected_err {
			t.Logf("Expected: %s\nGot: %s\n", expected_err, err)
			t.Fail()
		}
	}

	expected_err = "Fullname is greater than 128 characters"
	name_too_long := strings.Repeat("#", 129)
	err := validateFullname(name_too_long)
	if err.Error() != expected_err {
		t.Logf("Expected: %s\nGot: %s\n", expected_err, err)
		t.Fail()
	}
}

func Test_validateUsername(t *testing.T) {
	expected_err := "Username does not begin with a roman alphabetic character"
	first_char_non_alpha := []string{"1abcdefghijk", "0asdasdadad", "9asdadasaffafa"}
	for _, val := range first_char_non_alpha {
		err := validateUsername(val)
		if err.Error() != expected_err {
			t.Logf("Expected: %s\nGot: %s\n", expected_err, err)
			t.Fail()
		}
	}

	expected_err = "Username is not alphanumeric"
	non_alphanumeric := []string{
		"asdasdaasd@@", "sadad...adsads", "asdada////", "sadad'adsaad", "ewaew]asdada", "lllll[sadada", "asaæææææassas",
		"aϨϨϨϨϨasdaasd", "asdad日本語sdada", "asda⌘⌘⌘sadsa", "aaaa嗨嗨嗨嗨嗨", "sdadCześć", "AAAAAAaść",
	}
	for _, val := range non_alphanumeric {
		err := validateUsername(val)
		if err.Error() != expected_err {
			t.Logf("Failed for %s\nExpected: %s\nGot: %s\n", val, expected_err, err)
			t.Fail()
		}
	}

	expected_err = "Username is less than 8 characters"
	too_short := []string{"", "a", "aa", "aaa", "aaaa", "aaaaa", "aaaaaa", "aaaaaaa"}
	for _, val := range too_short {
		err := validateUsername(val)
		if err.Error() != expected_err {
			t.Logf("Failed for %s\nExpected: %s\nGot: %s\n", val, expected_err, err)
			t.Fail()
		}
	}

	expected_err = "Username is greater than 64 characters"
	name_too_long := strings.Repeat("a", 65)
	err := validateUsername(name_too_long)
	if err.Error() != expected_err {
		t.Logf("Expected: %s\nGot: %s\n", expected_err, err)
		t.Fail()
	}

	valid_usernames := []string{"bobman2000", "BOBman9000", "llll22222", "L33tBoaaa", "Joseph1111", "l34444444", "adadasdadadadasasfasfasff"}
	for _, val := range valid_usernames {
		if err := validateUsername(val); err != nil {
			t.Logf("Expected: nil\nGot: %s\n", err)
			t.Fail()
		}
	}
}

func Test_validateEmail(t *testing.T) {
	expected_err := "Invalid email format"
	invalid_emails := []string{"&&&&adad@gmailcom", "@@@", "ayyyy", "", "()()", "^^broseph.com", "%%ttt@gmailom", "%%sssss.com", "word@land"}

	for _, invalid_email := range invalid_emails {
		err := validateEmail(invalid_email)

		if err.Error() != expected_err {
			t.Logf("Failed input: %s\n Expected: %s\nGot %s\n", invalid_email, expected_err, err)
			t.Fail()
		}
	}

	should_pass := []string{"username@gmail.com", "username@hotmail.com", "username@protonmail.com", "username@protonmail.ch"}
	for _, email := range should_pass {
		err := validateEmail(email)
		if err != nil {
			t.Logf("Case %s failed to validate", email)
			t.Fail()
		}
	}
}
