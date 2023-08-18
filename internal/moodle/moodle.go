package moodle

import (
	"encoding/json"
	"io"
	"moodle-sync-ldap/internal/utils"
	"net/http"
	"net/url"
	"strconv"
)

type course struct {
	ID                int64  `json:"id,omitempty"`
	UID               int64  `json:"uid,omitempty"`
	ShortName         string `json:"shortname,omitempty"`
	ShortNameTranslit string `json:"shortname_translit,omitempty"`
	CategoryID        int    `json:"categoryid,omitempty"`
	CategorySortOrder int    `json:"categorysortorder,omitempty"`
	FullName          string `json:"fullname,omitempty"`
	DisplayName       string `json:"displayname,omitempty"`
	IDNumber          string `json:"idnumber,omitempty"`
	Users             []user `json:"users,omitempty"`
}

type user struct {
	ID        int64  `json:"id,omitempty"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty"`
	FullName  string `json:"fullname,omitempty"`
	Email     string `json:"email,omitempty"`
}

type Config struct {
	URL   string `json:"url" xml:"url" yaml:"url"`
	Token string `json:"token" xml:"token" yaml:"token"`
}

func (c *Config) GetData() (courses []course, err error) {
	u, err := url.Parse(c.URL)
	if err != nil {
		return
	}
	u.RawQuery = ""
	q := u.Query()
	q.Add("wstoken", c.Token)
	q.Add("moodlewsrestformat", "json")
	q.Add("wsfunction", "core_course_get_courses")
	u.RawQuery = q.Encode()
	resp, err := get(u)
	if err != nil {
		return
	}
	if err = json.Unmarshal(resp, &courses); err != nil {
		return
	}
	for n := range courses {
		courses[n].ShortNameTranslit = utils.Transliterate(courses[n].ShortName)
		courses[n].UID = courses[n].ID + 900000
		u.RawQuery = ""
		q := u.Query()
		q.Add("wstoken", c.Token)
		q.Add("moodlewsrestformat", "json")
		q.Add("wsfunction", "core_enrol_get_enrolled_users")
		q.Add("courseid", strconv.FormatInt(courses[n].ID, 10))
		u.RawQuery = q.Encode()
		if resp, err := get(u); err == nil {
			json.Unmarshal(resp, &courses[n].Users)
		}
	}
	return
}

func get(u *url.URL) ([]byte, error) {
	res, err := http.Get(u.String())
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()
	res.Close = true
	return io.ReadAll(res.Body)
}
