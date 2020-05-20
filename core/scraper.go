package core

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/mcnijman/go-emailaddress"
	log "github.com/prometheus/common/log"
)

type emailSource map[string]interface{}
type EmailValid struct {
	Email string
	Valid bool
}

type linkComp struct {
	protocol string
	domain   string
	path     string
	rootURL  string
	orgURL   string
}

type dirtyComp struct {
	urlStruct linkComp
	href      string
}

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

func getLinks(resp string, urlComp linkComp) []string {

	var cleanLinks []string
	soupHTML := soup.HTMLParse(resp)
	links := soupHTML.FindAll("a")

	dirty := make(chan dirtyComp, len(links))
	clean := make(chan string, len(links))

	for i := 0; i < len(links); i++ {
		go cleanURLWorker(dirty, clean)
	}

	for _, l := range links {
		pack := dirtyComp{
			urlStruct: urlComp,
			href:      l.Attrs()["href"],
		}
		dirty <- pack
	}
	close(dirty)

	for i := 0; i < len(links); i++ {
		l := <-clean
		cleanLinks = append(cleanLinks, l)
	}
	close(clean)

	cleanLinks = unique(cleanLinks)

	return cleanLinks
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
					s.EmailValid = append(s.EmailValid, EmailValid{Email: email, Valid: valid})
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

func cleanEmail(email string) string {
	switch {
	case strings.HasPrefix(email, "u0"):
		email = email[5:]
	case strings.HasPrefix(email, "%20"):
		email = email[3:]
	}
	return email
}

func getEmails(html, sourceUrl string) []emailSource {

	var emailList []emailSource

	text := []byte(html)

	for _, e := range emailaddress.FindWithIcannSuffix(text, false) {
		email, err := emailaddress.Parse(e.String())
		if err != nil {
			message := fmt.Sprintf("Error invalid email: %s %s", e, err)
			log.Debug(message)
		}
		em := cleanEmail(email.String())
		if !containsURL(blacklistEmails, em) {
			emailList = append(emailList, emailSource{em: sourceUrl})
		}
	}
	return emailList

}

func validateEmail(e string) bool {
	email, err := emailaddress.Parse(e)
	if err != nil {
		msg := fmt.Sprint("error whilst passing emails")
		log.Debug(msg)
		return false
	}
	err = email.ValidateHost()
	if err != nil {
		msg := fmt.Sprint("invalid host")
		log.Debug(msg)
		return false
	}
	return true
}

type timeOutR struct {
	link   []string
	emails []emailSource
}

func scrap(targetURL string) ([]string, []emailSource) {

	c1 := make(chan timeOutR)
	urlComp, err := getURLComp(targetURL)
	if err != nil {
		message := fmt.Sprintf("Error when parsing URL: %s %s", targetURL, err)
		log.Debug(message)
		return []string{}, []emailSource{}
	}
	var op timeOutR
	var links []string
	var emails []emailSource
	// log.Print(targetURL)

	go func() {
		defer close(c1)
		// time.Sleep(2 * time.Second)
		resp, err := soup.Get(urlComp.orgURL)
		if err != nil {
			message := fmt.Sprintf("Error when getting: %s %s", targetURL, err)
			log.Debug(message)

			c1 <- timeOutR{link: []string{}, emails: []emailSource{}}
		}

		emails := getEmails(resp, urlComp.orgURL)
		links := getLinks(resp, urlComp)

		c1 <- timeOutR{link: links, emails: emails}
	}()

	select {
	case op = <-c1:

		links = op.link
		emails = op.emails
	case <-time.After(20 * time.Second):
		// message := fmt.Sprintf("Timeout on url: %s", targetURL)
		// log.Print(message)
		return []string{}, []emailSource{}
	}

	return links, emails
}
