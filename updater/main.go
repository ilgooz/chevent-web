package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	checkPeriod = flag.Duration("period", 3*time.Minute, "check prediod for new version of the repo")

	mongoAddr = flag.String("mongo", "mongodb://127.0.0.1:27017/chevent", "Mongodb Address")
	session   *mgo.Session

	webSHA  string
	jsonSHA string
)

func main() {
	flag.Parse()

	var err error
	session, err = mgo.Dial(*mongoAddr)
	if err != nil {
		log.Fatalln(err)
	}
	session.SetSafe(&mgo.Safe{})

	fmt.Println("updater is running")

	for {
		if err := do(); err != nil {
			fmt.Println(fmt.Sprintf("error while running: %s", err.Error()))
		}

		fmt.Println("waiting for some")
		time.Sleep(*checkPeriod)
	}
}

type eventJSON struct {
	Items []Event `json:"items"`
}

type Event struct {
	ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name        string        `json:"name"`
	Date        Time          `json:"date"`
	Free        bool          `json:"free"`
	Image       string        `json:"image"`
	URL         string        `json:"url"`
	Description string        `json:"description"`
	Quota       int           `json:"quota"`
	Speakers    []Speaker     `json:"speakers"`
}

type Speaker struct {
	Name    string `json:"name"`
	Subject string `json:"subject"`
}

type Time time.Time

func (t *Time) UnmarshalJSON(b []byte) error {
	s := string(b)
	tt, err := time.Parse("02.01.2006", s[1:len(s)-1])
	*t = Time(tt)
	return err
}

func do() error {
	// chevent-web
	fmt.Println("[chevent-web] checking")

	if _, err := exec.Command("git", "pull", "origin", "master").Output(); err != nil {
		return err
	}

	b, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		return err
	}
	newSHA := string(b)

	if webSHA == newSHA {
		fmt.Println("[chevent-web] up-to-date")
	} else {
		fmt.Println(fmt.Sprintf("[chevent-web] new version found %s building", newSHA))

		c := exec.Command("make", "all")
		c.Dir = "/go/src/github.com/codeui/chevent-web"
		if err := c.Start(); err != nil {
			return err
		}
		if err := c.Wait(); err != nil {
			return err
		}

		webSHA = newSHA

		fmt.Println("[chevent-web] done")
	}

	// chevent/events/events.json
	fmt.Println("[events.json] checking")

	c := exec.Command("git", "pull", "origin", "master")
	c.Dir = "/go/src/github.com/codeui/chevent-web/chevent"
	if _, err = c.Output(); err != nil {
		return err
	}

	b, err = exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		return err
	}
	newSHA = string(b)

	if jsonSHA == newSHA {
		fmt.Println("[events.json] up-to-date")
	} else {
		fmt.Println("[events.json] updating")

		file, err := os.Open("/go/src/github.com/codeui/chevent-web/chevent/events/events.json")
		if err != nil {
			return err
		}
		defer file.Close()

		var items eventJSON

		if err := json.NewDecoder(file).Decode(&items); err != nil {
			return err
		}

		for _, event := range items.Items {
			m := bson.M{
				"name": event.Name,
				"date": event.Date,
			}

			if _, err := session.DB("").C("events").Upsert(m, &event); err != nil {
				log.Println(err)
				return err
			}
		}

		jsonSHA = newSHA

		fmt.Println("[events.json] done")
	}

	return nil
}
