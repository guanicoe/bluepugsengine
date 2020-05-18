package core

import (
	"fmt"
	"reflect"
	"testing"
)

type getURLCompTest struct {
	url string
	exp linkComp
}

var getURLCompTestList = []getURLCompTest{
	{"https://google.com",
		linkComp{protocol: "https",
			domain:  "google.com",
			path:    "",
			rootURL: "https://google.com",
			orgURL:  "https://google.com",
		},
	}, {"https://google.com/",
		linkComp{protocol: "https",
			domain:  "google.com",
			path:    "/",
			rootURL: "https://google.com",
			orgURL:  "https://google.com/",
		},
	}, {"https://google.com/dir1/dir2",
		linkComp{protocol: "https",
			domain:  "google.com",
			path:    "/dir1/dir2",
			rootURL: "https://google.com",
			orgURL:  "https://google.com/dir1/dir2",
		},
	}, {"",
		linkComp{},
	},
}

func TestGetURLComp(t *testing.T) {
	for _, test := range getURLCompTestList {
		result, _ := getURLComp(test.url)
		if result != test.exp {
			t.Fatal(result)
		}
	}
}

type cleanURLTest struct {
	url  linkComp
	href string
	exp  string
}

var cleanURLTestList = []cleanURLTest{
	{linkComp{protocol: "https",
		domain:  "google.com",
		path:    "",
		rootURL: "https://google.com",
		orgURL:  "https://google.com"},
		"//blog", "https://blog"},
	{linkComp{protocol: "https",
		domain:  "google.com",
		path:    "",
		rootURL: "https://google.com",
		orgURL:  "https://google.com"},
		"/blog", "https://google.com/blog"},
	{linkComp{protocol: "https",
		domain:  "google.com",
		path:    "",
		rootURL: "https://google.com",
		orgURL:  "https://google.com"},
		"./blog", "https://google.com/blog"},
	{linkComp{protocol: "https",
		domain:  "google.com",
		path:    "",
		rootURL: "https://google.com",
		orgURL:  "https://google.com"},
		"#blog", "https://google.com"},
	{linkComp{protocol: "https",
		domain:  "google.com",
		path:    "",
		rootURL: "https://google.com",
		orgURL:  "https://google.com"},
		"%20/blog", "https://google.com"},
	{linkComp{protocol: "https",
		domain:  "google.com",
		path:    "",
		rootURL: "https://google.com",
		orgURL:  "https://google.com"},
		"@google", "https://google.com"},
	{linkComp{protocol: "https",
		domain:  "google.com",
		path:    "",
		rootURL: "https://google.com",
		orgURL:  "https://google.com"},
		"?id=1", "https://google.com/?id=1"},
	{linkComp{protocol: "https",
		domain:  "google.com",
		path:    "",
		rootURL: "https://google.com",
		orgURL:  "https://google.com"},
		"mailto:email@blo.fr", "https://google.com"},
	{linkComp{protocol: "https",
		domain:  "google.com",
		path:    "",
		rootURL: "https://google.com",
		orgURL:  "https://google.com"},
		"javascript:void(0)", "https://google.com"},
	{linkComp{protocol: "https",
		domain:  "google.com",
		path:    "",
		rootURL: "https://google.com",
		orgURL:  "https://google.com"},
		"", "https://google.com"},
	{linkComp{protocol: "https",
		domain:  "google.com",
		path:    "",
		rootURL: "https://google.com",
		orgURL:  "https://google.com"},
		"https://google.com", "https://google.com"},
}

func TestCleanURL(t *testing.T) {
	for _, test := range cleanURLTestList {
		result := cleanURL(test.url, test.href)
		if result != test.exp {
			fmt.Println(result)
			t.Fatal(test.href)
		}
	}
}

// type CleanFileNameTest struct {
// 	str string
// 	exp string
// }
//
// var CleanFileNameTestList = []CleanFileNameTest{
// 	{"a.json", "a.json"},
// }
//
// func TestCleanFileName(t *testing.T) {
// 	for _, test := range CleanFileNameTestList {
// 		// str := &test.str
// 		CleanFileName(&test.str)
// 		if test.str != test.exp {
// 			fmt.Println(result)
// 			t.Fatal(test.str)
// 		}
// 	}
// }

func TestContainsEXT(t *testing.T) {
	pos := containsEXT(blacklistEXT, "https://www.domain.com/pres.pdf")
	if !pos {
		t.Errorf("Answer is %t", pos)
	}

	pos = containsEXT(blacklistEXT, "https://www.domain.com/pres.hgf")
	if pos {
		t.Errorf("Answer is %t", pos)
	}
}

func TestContainsURL(t *testing.T) {

	pos := containsURL(backlistURL, "https://www.twitter.com/legal")
	if !pos {
		t.Errorf("Answer is %t", pos)
	}

	pos = containsURL(backlistURL, "https://www.bichromatics.com/")
	if pos {
		t.Errorf("Answer is %t", pos)
	}

	pos = containsURL([]string{""}, "https://www.bichromatics.com/")
	if pos {
		t.Errorf("Answer is %t", pos)
	}
}

func TestContains(t *testing.T) {
	var testList = []string{"https://www.bichromatics.com/calculator", "https://www.bichromatics.com/calculator/download_file", "https://www.bichromatics.com/calculator", "https://www.bichromatics.com/calculator", "https://www.bichromatics.com/calculator/download_file", "https://www.bichromatics.com/calculator"}
	pos := contains(testList, "https://www.bichromatics.com/calculator")
	if !pos {
		t.Errorf("Answer is %t", pos)
	}

	pos = contains(testList, "d")
	if pos {
		t.Errorf("Answer is %t", pos)
	}
}

func TestUnique(t *testing.T) {
	var testList = []string{"a", "a", "c"}
	var expected = []string{"a", "c"}
	pos := unique(testList)
	bol := reflect.DeepEqual(pos, expected)
	if !bol {
		t.Errorf("Answer is %t, %s", bol, pos)
	}

}

// func benchmarkFib(b *testing.B) {
// 	for n := 0; n < b.N; n++ {
// 		main()
// 	}
// }
