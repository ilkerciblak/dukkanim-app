package viladition_test

import (
	"dukkanim-api/pkg/viladition"
	"fmt"
	"strings"
	"testing"
)

func Test__String_Validator(t *testing.T) {
	cases := []struct {
		Name            string
		Value           string
		ExpectedResult  int
		DoesExpectError bool
	}{
		{
			"Empty String should Error",
			"",
			1,

			false,
		},
	}

	for _, tc := range cases {
		t.Run(
			tc.Name,
			func(t *testing.T) {
				got := viladition.String(tc.Value).Required().MinLength(3).MaxLength(255).GetErrors()

				if tc.ExpectedResult >= len(got) {
					s := strings.Join(got, " ")
					fmt.Printf("%s", s)
					t.Errorf("Excepted %d got %d", tc.ExpectedResult, len(got))
				}

			},
		)
	}

}

func Test__Email_Validator(t *testing.T) {
	cases := []struct {
		Name           string
		Value          string
		ExpectedResult int
	}{
		{
			"Correct Input",
			"example.a@gmail.com",
			0,
		},
		{
			"Empty Input Should Return Error",
			"",
			1,
		},
		{
			"Space Character in Email Should Error",
			"examp le.a@gmail.com",
			1,
		},
	}
	for _, tc := range cases {
		t.Run(
			tc.Name,
			func(t *testing.T) {
				got := viladition.Email(tc.Value).Required().Validate().GetErrors()

				if tc.ExpectedResult > len(got) {
					s := strings.Join(got, " ")
					fmt.Printf("%s", s)
					t.Errorf("Excepted %d got %d", tc.ExpectedResult, len(got))
				}

			},
		)
	}

}
