package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPhoneNumber(t *testing.T) {
	tt := []struct {
		name   string
		input  string
		expect string
	}{
		{
			name:   "Only characters",
			input:  "adfmawfkjlawnfbjkaw",
			expect: "",
		},
		{
			name:   "Too short number",
			input:  "+3809549752",
			expect: "",
		},
		{
			name:   "Too long number",
			input:  "+3809549752222",
			expect: "",
		},
		{
			name:   "Valid with full code +38",
			input:  "+380954975222",
			expect: "+380954975222",
		},
		{
			name:   "Valid without plus 38",
			input:  "80954975222",
			expect: "+380954975222",
		},
		{
			name:   "Valid with non full code 8",
			input:  "0954975222",
			expect: "+380954975222",
		},
		{
			name:   "Valid with different characters",
			input:  "+(380)95-49-75-222",
			expect: "+380954975222",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expect, PhoneNumber(tc.input))
		})
	}
}
