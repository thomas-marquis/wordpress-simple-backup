# wordpress-simple-backup

A simple backup CLI for wordpress (content and database)

## Installation

TODO

## Configuration

TODO

## Usage

```bash
# Save the current wordpress content and database (this will create a new backup version)
wsb save

# List all the backup versions
wsb list

# Restore the last backup
wsb restore

# Restore a specific version of the backup
wsb restore -v 1
```

## Local development setup

TODO

```bash
go get .
go install go.uber.org/mock/mockgen@latest
go install github.com/spf13/cobra-cli@latest
go generate
```

## Launch tests

TODO

## External resources links

* [gomock](https://github.com/uber-go/mock)
* [cobra](https://github.com/spf13/cobra)
