package fizzbuzzservice

import "fmt"

type POSTFizzBuzzParameters struct {
	Limit int    `json:"limit"`
	Int1  int    `json:"int1"`
	Int2  int    `json:"int2"`
	Str1  string `json:"str1"`
	Str2  string `json:"str2"`
}

type FizzBuzzService struct{}

func (fbservice *FizzBuzzService) FizzBuzz(limit, int1, int2 int, str1, str2 string) string {
	results := ""
	for i := 1; i <= limit; i++ {
		result := ""
		if i%int1 == 0 {
			result += str1
		}
		if i%int2 == 0 {
			result += str2
		}
		if result != "" {
			results += result + "\n"
			continue
		}
		results += fmt.Sprintf("%d\n", i)
	}
	return results
}
