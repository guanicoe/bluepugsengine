package core

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var backlistURL = []string{"t.co", "whatsapp", "github", "bing", "twitter", "flicker", "turner", "youtu", "flickr", "commailto", "pinterest", "linkedin", "zencart", "wufoo", "youcanbook", "instagram"}
var blacklistEXT = []string{"jpeg", "jpg", "gif", "pdf", "png", "ppsx", "f4v", "mp3", "mp4", "exe", "dmg", "zip", "avi", "wmv", "pptx", "exar1", "edx", "epub"}

func getURLComp(targetURL string) (linkComp, error) {
	u, err := url.Parse(targetURL)
	if err != nil || len(targetURL) == 0 {

		return linkComp{}, err
	}
	r := linkComp{protocol: u.Scheme,
		domain:  u.Host,
		path:    u.Path,
		rootURL: fmt.Sprintf("%s://%s", u.Scheme, u.Host),
		orgURL:  targetURL,
	}
	return r, nil
}

func cleanURL(urlComp linkComp, href string) string {
	var newHref string

	switch {
	case strings.HasPrefix(href, "//"):
		newHref = strings.Join([]string{urlComp.protocol, href}, ":")
	case strings.HasPrefix(href, "/"):
		newHref = strings.Join([]string{urlComp.rootURL, href}, "")
	case strings.HasPrefix(href, "./"):
		newHref = strings.Join([]string{urlComp.rootURL, href[1:]}, "")
	case strings.HasPrefix(href, "#"):
		newHref = urlComp.rootURL
	case strings.HasPrefix(href, "%20"):
		newHref = urlComp.rootURL
	case strings.HasPrefix(href, "@"):
		newHref = urlComp.rootURL
	case strings.HasPrefix(href, "?"):
		newHref = strings.Join([]string{urlComp.rootURL, href}, "/")
	case strings.HasPrefix(href, "mailto:"):
		newHref = urlComp.rootURL
	case strings.HasPrefix(href, "javascript:"):
		newHref = urlComp.rootURL
	case len(href) == 0:
		newHref = urlComp.rootURL
	default:
		newHref = href
	}

	cleanedUp := strings.TrimSuffix(newHref, "/")
	_, err := url.ParseRequestURI(strings.Join([]string{cleanedUp, "/"}, ""))
	if err != nil {
		return ""
	}
	// fmt.Println(cleanedUp)

	return cleanedUp
}

func cleanURLWorker(dirty <-chan dirtyComp, clean chan<- string) {

	for d := range dirty {
		clean <- cleanURL(d.urlStruct, d.href)
	}

}

func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsURL(s []string, e string) bool {

	for _, a := range s {
		match, _ := regexp.MatchString(a, e)

		if match && len(a) > 0 {
			return true
		}
	}
	return false
}

func containsEXT(s []string, e string) bool {
	for _, a := range s {
		if strings.HasSuffix(e, a) {
			return true
		}
	}
	return false
}

func setUniqueEmail(s *jobData) {
	for _, emailMap := range s.emailList {
		for email, _ := range emailMap {
			if !contains(s.emailListUnique, email) {

				s.emailListUnique = append(s.emailListUnique, email)
				if s.paramPointer.CheckEmails {
					valid := validateEmail(email)
					fmt.Println(email, "-", valid)
				} else {
					fmt.Println(email, "- not checked")
				}

			}
		}

	}
}

func newValidURL(l string, j *jobData, s string) bool {
	switch {
	case len(l) == 0:
		return false
	case contains(j.scrapedRecv, l):
		return false
	case contains(j.scrapedSent, l):
		return false
	case containsEXT(blacklistEXT, l):
		return false
	case containsURL(backlistURL, l):
		return false
	case !containsURL([]string{s}, l):
		return false
	default:
		return true

	}
}
