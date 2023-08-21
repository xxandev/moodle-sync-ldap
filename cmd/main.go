package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"moodle-sync-ldap/internal/archive"
	"moodle-sync-ldap/internal/utils"
	"os"
	"path/filepath"
	"text/template"
)

var config Config

func init() {
	flag.StringVar(&config.Moodle.URL, "url", "", "moodle api url [https://moodle.domain/webservice/rest/server.php]")
	flag.StringVar(&config.Moodle.Token, "token", "", "moodle token [xxxxxxxxxxxxxx]")
	flag.StringVar(&config.LDAP.DN, "dn", "", "ldap dn [ou=example,dc=example,dc=com]")
	flag.StringVar(&config.fileLDIF, "fd", "moodle.ldif", "file ldif [example.ldif]")
	flag.Parse()
	if err := config.Check(); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	fmt.Println("start moodle-sync-ldap")
	rp, err := os.Executable()
	if err != nil {
		log.Fatalf("error get root path: %v\n", err)
	}
	if err := archive.Unpack(filepath.Join(filepath.Dir(rp), "moodle-sync-ldap-repository")); err != nil {
		log.Printf("warning unpack archive: %v\n", err)
	}
	mdata, err := config.Moodle.GetData()
	if err != nil {
		log.Fatalf("error get moodle data: %v\n", err)
	}
	fmt.Println("==============================================================")
	res, _ := json.MarshalIndent(&mdata, "", "\t")
	fmt.Println(string(res))

	content, err := os.ReadFile(filepath.Join(filepath.Dir(rp), "moodle-sync-ldap-repository", "moodle", "ldif.tmpl"))
	if err != nil {
		log.Fatalf("error read ldif file: %v\n", err)
	}
	var buffer bytes.Buffer
	template.Must(template.New("moodle-ldap").Parse(fmt.Sprintf(string(content), config.LDAP.DN))).Execute(&buffer, &mdata)
	fmt.Println("==============================================================")
	fmt.Println(buffer.String())
	utils.CreateFile(config.fileLDIF, buffer.Bytes())
}
