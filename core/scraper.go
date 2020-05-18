package core

import (
	"fmt"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/mcnijman/go-emailaddress"
	log "github.com/prometheus/common/log"
)

type emailSource map[string]interface{}

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

func getLinks(resp string, urlComp linkComp) []string {

	var cleanLinks []string
	soupHTML := soup.HTMLParse(resp)
	links := soupHTML.FindAll("a")

	dirty := make(chan dirtyComp, len(links))
	clean := make(chan string, len(links))

	for i := 0; i < len(links); i++ {
		go CleanURLWorker(dirty, clean)
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

	cleanLinks = Unique(cleanLinks)

	return cleanLinks
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
		emailList = append(emailList, emailSource{email.String(): sourceUrl})
	}
	return emailList

}

type timeOutR struct {
	link   []string
	emails []emailSource
}

func scrap(targetURL string) ([]string, []emailSource) {

	c1 := make(chan timeOutR)
	urlComp, err := GetURLComp(targetURL)
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
