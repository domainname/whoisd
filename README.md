
# Whois Daemon

*A quick and easy way to setup your own WHOIS server with Elasticsearch/MySQL storage*

Whois Daemon complies with the requirements of [ICANN](https://www.icann.org/resources/pages/approved-with-specs-2013-09-17-en)

[![Build Status](https://travis-ci.org/takama/whoisd.png?branch=master)](https://travis-ci.org/takama/whoisd)

**Whois Daemon** is a light server which provide fast way to present whois information. The daemon based on Elasticsearch storage (Mysql storage still in development). 

### Install

This package is "go-gettable", just do:

```sh
go get github.com/takama/whoisd
```

_NOTE_: you need go 1.2+ and need set PATH="$GOPATH/bin". Please check your installation with

```sh
go version
go env
```

### Testing

```sh
whoisd -t -config=test/testconfig.conf -mapping=test/testmapping.json
```

### Running

Start of whoisd:

```sh
sudo whoisd
```
or start as daemon
```sh
sudo whoisd install
sudo whoisd start
```

This will bring up whoisd listening on port 43 for client communication.

### Usage

```
whoisd - Whois Daemon

Usage:
  whoisd install | remove | start | stop | status
  whoisd [ -t | --test ] [ -option | -option ... ]
  whoisd -h | --help
  whoisd -v | --version

Commands:
  install           Install as service (is only valid for Linux and Mac Os X)
  remove            Remove service
  start             Start service
  stop              Stop service
  status            Check service status

  -h --help         Show this screen
  -v --version      Show version
  -t --test         Test mode

Options:
  -config=<path>    Path to config file (used in /etc/whoisd/whoisd.conf)
  -mapping=<path>   Path to mapping file (used in /etc/whoisd/conf.d/mapping.json)
  -host=<host/IP>   Host name or IP address
  -port=<port>      Port number
  -work=<number>    Number of active workers (default 1000)
  -conn=<number>    Number of active connections (default 1000)
  -storage=<type>   Type of storage (Elasticsearch, Mysql or Dummy for testing)
  -shost=<host/IP>  Storage host name or IP address
  -sport=<port>     Storage port number
  -base=<name>      Storage index or database name
  -table=<name>     Storage type or table name
```

### Config

The config file should be in /etc/whoisd/whoisd.conf. Of course possible to load config settings from any other place through -config option. If config file is absent, used default configuration below: 

```json
{
  "host": "0.0.0.0",
  "port": 43,
  "workers": 1000,
  "connections": 1000,
  "storage": {
    "storageType": "Dummy",
    "host": "localhost",
    "port": 9200,
    "indexBase": "whois",
    "typeTable": "domain"
  }
}
```
_NOTE_: Valid storage types: Elasticsearch, Mysql, Dummy. Dummy storage has two records for testing: "example.tld" and "google.com". You can test it: 
```sh
whois -h localhost example.tld
```
or
```sh
whois -h localhost google.com
```
These fixtures placed in "storage" package directory.


### Mapping

All required fields for whoisd must be defined in the mapping file. The mapping file represent all fields in your database as key names in the whoisd. The mapping file should be in /etc/whoisd/conf.d/mapping.json. It possible to load mapping file through -mapping option. The context of the mapping file is described below:

```json
{
  "Fields" : {
    "01": {
      "key": "Domain Name: ",
      "name": ["name"],
      "related": "name"
    },
    "02": {
      "key": "Registry Domain ID: ",
      "name": ["domainId"],
      "related": "name"
    },
    "03": {
      "key": "Registrar WHOIS Server: ",
      "value": ["whois.yourwhois.eu"]
    }

}
```

- "Fields" - a list of fields from "01" to last number "nn" in ascending order
- "key" - the label for the field (preinstalled config file has keys according to ICANN requirements)
- "value" - use it if the field has prearranged value (not use any field from the database)
- "name" - the name of the field in the database, if the field is not prearranged ("value" is not defined)
- "related" - the name of the field in the database through which the request for 

```json
{

    "06": {
      "key": "Creation Date: ",
      "name": ["creationDate"],
      "format": "{date}",
      "related": "name"
    }

}
```

- "format" - special instructions to indicate how to display field, examples of the use of this will be shown below
- "{date}" - used in the format to indicate that the field is a date and need special formatting of the date RFС3339

```json
{

    "12": {
      "key": "Domain Status: ",
      "name": ["domainStatus"],
      "multiple": true,
      "related": "name"
    },

    "52": {
      "key": "Name Server: ",
      "name": ["name"],
      "multiple": true,
      "related": "nsgroupId",
      "relatedBy": "nsgroupId",
      "relatedTo": "nameserver"
    }

}
```

- "multiple" - if this option is set to 'true', then for each value will be repeated label in whois output like that:
```
Name Server: ns1.example.com
Name Server: ns2.example.com
Name Server: ns3.example.com
```
- "relatedBy" - the name of the field in the database through which the related request for 
- "relatedTo" - the name of the table/type in the database through which made a relation

```json
{

    "13": {
      "key": "Registry Registrant ID: ",
      "name": ["handle"],
      "hide": true,
      "related": "ownerHandle",
      "relatedBy": "handle",
      "relatedTo": "customer"
    },

}
```

- "hide" - if this option is set to 'true', the value of the field will not shown in whois output

```json
{

    "14": {
      "key": "Registrant Name: ",
      "name": ["name.fullName"],
      "related": "ownerHandle",
      "relatedBy": "handle",
      "relatedTo": "customer"
    },

    "40": {
      "key": "Tech Name: ",
      "name": ["name.firstName", "name.lastName"],
      "related": "techHandle",
      "relatedBy": "handle",
      "relatedTo": "customer"
    }

}
```

- "name": ["name.fullName"] - use dot notation for embedded fields (MySQL storage not allowed)
- "name": ["name.firstName", "name.lastName"] - all these fields will be joined by default

```json
{

    "21": {
      "key": "Registrant Phone: ",
      "name": ["phone.countryCode", "phone.areaCode", "phone.subscriberNumber"],
      "format": "{string}.{string}{string}",
      "related": "ownerHandle",
      "relatedBy": "handle",
      "relatedTo": "customer"
    }

}
```

- "format": "{string}.{string}{string}" - indicate that the fields ["phone.countryCode", "phone.areaCode", "phone.subscriberNumber"] need special formatting with described format (they are not simple joined)
- {string} - represent one string field in format option

```json
{

    "55": {
      "key": "",
      "value": [""],
      "format": ">>> Last update of WHOIS database: {date} <<<"
    }
}
```

- is the example of the formating where used {date}, because the name of the field has not present, the result is {date} will be replaced by CURRENT date in RFC3339 format


### TODO

- in memory storage
- Rest API
- update storage records by Rest API


Copyright (c) 2014 Igor Dolzhikov

[MIT License](https://github.com/takama/whoisd/blob/master/LICENSE)
