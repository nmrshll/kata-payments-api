package fizzbuzzservice

import "testing"

func TestFizzBuzz(t *testing.T) {
	fbservice := new(FizzBuzzService)
	type test struct {
		got  string
		want string
	}

	tests := []test{
		{got: fbservice.FizzBuzz(16, 3, 5, "Fizz", "Buzz"), want: "1\n2\nFizz\n4\nBuzz\nFizz\n7\n8\nFizz\nBuzz\n11\nFizz\n13\n14\nFizzBuzz\n16\n"},
		{got: fbservice.FizzBuzz(16, 3, 5, "fizz", "buzz"), want: "1\n2\nfizz\n4\nbuzz\nfizz\n7\n8\nfizz\nbuzz\n11\nfizz\n13\n14\nfizzbuzz\n16\n"},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("FizzBuzz(...) \n got: \n%v \n want: \n%v", tt.got, tt.want)
		}
	}
}
