# **Moodle Sync LDAP**

moodle-sync-ldap - utility for synchronizing courses and users as added to courses.

## **Installation**
- install [golang](https://go.dev/)
- go get github.com/xxandev/moodle-sync-ldap
- cd ..../moodle-sync-ldap
- make [ build | arm6 | arm7 | arm8 | linux64 | linux32 | win64 | win32 | win64i | win32i ] or go build .

## **Run**
```bash
.../moodle-sync-ldap \
    -url https://moodle.domain/webservice/rest/server.php \
    -token you-token-afsadfsdf231232134 \
    -dn dc=example,dc=lan
```

## **Setup**
**. . .**