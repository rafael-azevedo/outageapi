# **Outage API** <sup><sub>_API Endpoints for the Outage Service_</sub></sup>
----

1. [Overview](#overview)
2. [Getting Started](#getting-started)
  * [Creating Config File](#creating-config-file)
  * [Setting Enviroment](#setting-enviroment)
  * [Installation](#installation)
  * [Running the Outage API](#running-the-outage-api)
3. [Using the EncryptCLI tool](#using-the-encryptcli-tool)


## Overview

**Outage API** creates an api endpoint for the outage service the following enpoints are included. 

```
 host:port/outage : GET
 host:port/assign : POST
 host:port/deassign : POST
 host:port/status/{id} : Get
```
Currently runs on linux and osx

## **Getting Started**

### Creating Config File

The config file must be called app.toml  

#### object

* oracle - Oracle object in the config file that holds parameters required to connect to an oracle database 

#### parameters

* username - username used to connect to a oracle db 
* password - Password encrypted by the included EncyrptCLI tool
* hostname - hostname where oracle database is located 
* port - oracle database port
* servicename - oracle database service name 


for example `app.toml` with following content:
```toml
---
[oracle]
username = "admin"
password = "188801b31601742f6910b69cc4d195571f9d7e455d70062309f01bd0"
hostname = "oracledb.staples.com"
port = 	"1521"
servicename = "ODBSERVICE"
```

### Setting Enviroment
Configure the following variables in your enviroment  

* **OUTAGECONF** - Directory containting the config file
* **OUTAGELOGDIR** - Directory where the api will log the requests and request status
* **ORACLEDB** - Key Provided by the EncryptCLI tool  

either export log location to your current enviroment 
```
$ export OUTAGECONF=~/.go/src/github.com/rafael-azevedo/outageapi
$ export OUTAGELOGDIR=~/.go/src/github.com/rafael-azevedo/outageapi/log
$ export ORACLEDB=rhyNpuwtTG+fWuH8X41VUg==
```

### Installation 
The build script can be found at ~/.go/src/github.com/rafael-azevedo/outageapi/scripts/build.sh

```
$ cd github.com/rafael-azevedo/outageapi/scripts 
$ ./build.sh
2016-12-27 18:39:29 UTC [info] project path: /Users/azera001/.go/src/github.com/rafael-azevedo/outageapi
2016-12-27 18:39:29 UTC [info] Building outageapi to /Users/azera001/.go/src/github.com/rafael-azevedo/outageapi/bin/
```

### Running The Outage API
```
$ cd github.com/rafael-azevedo/outageapi/bin
$ ./outageapi

	 .d8888b.  888                      888
	d88P  Y88b 888                      888
	Y88b.      888                      888
	 "Y888b.   888888  8888b.  88888b.  888  .d88b.  .d8888b
	    "Y88b. 888        "88b 888 "88b 888 d8P  Y8b 88K
	      "888 888    .d888888 888  888 888 88888888 "Y8888b.
	Y88b  d88P Y88b.  888  888 888 d88P 888 Y8b.          X88
	 "Y8888P"   "Y888 "Y888888 88888P"  888  "Y8888   88888P'
	                           888
	                           888
	                           888

	  01010011 01110100 01100001 01110000 01101100 01100101
	  01110011 00100000 01001111 01110101 01110100 01100001
	  01100111 01100101 00100000 01000001 01010000 01001001


2016/12/27 14:25:19 Outage Api Started at http://127.0.0.1:1234
[Active Endpoints :
 127.0.0.1:1234/outage : GET
 127.0.0.1:1234/assign : POST
 127.0.0.1:1234/deassign : POST
 127.0.0.1:1234/status/{id} : Get
```

## **Using the EncryptCLI tool**

