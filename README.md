# votar
AUR vote utility

## As a CLI Util

*go install*

```sh
go install github.com/Jguer/votar/cmd/votar
```

Usage:

```
Usage: main [--vote VOTE] [--unvote UNVOTE] [--user USER] [--password PASSWORD]

Options:
  --vote VOTE, -v VOTE   packages to vote for
  --unvote UNVOTE, -u UNVOTE
                         packages to unvote for
  --user USER [env: AUR_USER]
  --password PASSWORD [env: AUR_PASSWORD]
  --help, -h             display this help and exit
```

Examples:

*vote for a package*

```sh
votar --user <user> --password <password> -v <package>
# or 
export AUR_USER=<user>
export AUR_PASSWORD=<password>
votar -v <package> -v <package 2>
```

*unvote for a package*

```sh
votar --user <user> --password <password> -u <package>
# or 
export AUR_USER=<user>
export AUR_PASSWORD=<password>
votar -u <unvote package> -v <vote package 2>
```

## As a library

```sh
go get github.com/Jguer/votar/pkg/vote
```

```go
package main

import (
	"context"
	"log"

	"github.com/Jguer/votar/pkg/vote"
)

client, err := vote.NewClient()
if err != nil {
    log.Println("Failed to create client")
}

client.SetCredentials("user", "password")


if err = client.Vote(context.Background(), "package");err != nil {
    log.Println("Failed to vote for", "package")
}
```
