

E-invoices is a cloud-based, cache-enabled, invoice intermediator.

You can send invoices to providers in Estonia and Finland. 
You can download invoices as xml to send them manually. The standarts are E-arve 1.2 and Finvoice 3.0
### Operator explanations
> * OMNIVA - estonian e-invoice provider
> * MAVENTA - finnish e-invoice provider
> * BYCOUNTRY - provider is chosen automatically by sale document customer country

### Requests
<details><summary>send</summary>

#### /send - export sale documents to operator.
##### Request body:
```json
{
	"operator": "all caps operator name.[OMNIVA, EARVELDAJA, MAVENTA, BYCOUNTRY]",
	"clientId": "client ID",
	"invoiceIds": "string array specifying invoice IDs from Back Office.", 
	"sessionKey": "Session key for API"
}
```
example request 
```json
{
    "invoiceIds": [
        "2031"
    ],
    "clientId": "1231231",
    "sessionKey": "test",
    "operator": "BYCOUNTRY"
}
```
##### Responses with explanation:
 * If everything goes smoothly you should receive this message 

```json
{
    "message": "documents were sent for processing",
    "statusCode": 200,
    "providerResponse": {
        "MAVENTA": {
            "message": "success",
            "statusCode": "201"
        },
        "OMNIVA": {
            "message": "success",
            "statusCode": "201"
        }
    }
}
```
 * If you did not add any credentials:
 ```json
 {
    "message": "Internal error",
    "statusCode": 500,
    "externalErrorResponse": {
        "message": "no provider credentials were added",
        "statusCode": "500"
    }
}
 ```

 * If your credentials are incorrect you will see next error responses
 ```json
{
    "message": "Provider service error",
    "statusCode": 400,
    "providerResponse": {
        "OMNIVA": {
            "message": "Illegal authentication phrase.",
            "statusCode": "ns0:60"
        }
    }
}
 ```
 ```json
 {
    "message": "Provider service error",
    "statusCode": 400,
    "providerResponse": {
        "MAVENTA": {
            "message": "ERROR: ERROR: VENDOR API KEY DISABLED",
            "statusCode": "400"
        }
    }
}
 ```
 * If some documents were declined on the phase of sending, you receive next response

```json

{
    "message":"documents were sent for processing",
    "statusCode":200,
    "providerResponse":{
        "EARVELDAJA":{
            "message":"success",
            "statusCode":"200",
            "documentError":{
                "887":{
                    "Status":"FAIL",
                    "Message":"Ettevõtja registrikoodiga 10090688 ei ole meie klient.",
                    "Provider":"EARVELDAJA"
                }
              }
        }
    }
}
```
 * If there is an error coming from external service, such as API, Provider, eBusinessRegister, etc.. this response structure will be send. Current case: session key is expired or wrong
 ```json
{
    "message": "api error",
    "statusCode": 400,
    "externalErrorResponse": {
        "message": "API: getSalesDocuments: error status: 1055",
        "statusCode": "1055"
    }
}
 ```
</details>
<details><summary>status</summary>

#### /status - check status of exported documents.
##### Request body:
```json
{
	"clientId": "client ID",
	"invoiceIds": "string array specifying invoice ID from Back Office.",
	"sessionKey": "Session key for API"
}
```
example request
```json
{
    "invoiceIds": [
        "86",
        "98",
        "198",
        "298"
    ],
    "clientId": "12312",
    "sessionKey": ""
}
```

##### Responses with explanation:
After send document request, status update request is scheduled after some delay. When document is sent it gets status PENDING,
which means it is being proccessed by operator. EXPORTED means documents were exported. DECLINED means documents has some wrong format or etc.
 * Successful response
 ```json
{
    "message": "invoices status report",
    "statusCode": 200,
    "exported": [
        86
    ],
    "pending": [
            98
    ],
    "declined": [
        198
    ],
    "notExported": [
            298
    ]
}
```
</details>
<details><summary>download</summary>

#### /download -  download documents in zip archive as XML files. The standarts are E-arve 1.2 and Finvoice 3.0.
##### Request body:
```json
{
    "invoiceIds": [
        "2083"
    ],
    "clientId": "12331",
    "sessionKey": "fdsfs",
    "operator": "BYCOUNTRY"
}
```
##### Response:
```json
{
    "message": "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAALAAAAMTkwMDA0MS54bWzMWMtu27wS3p+nELxuROpiXQpG+BVbOb/Pn7iGrQbtqmCsSSJEEg2JSuJ3OYsu+iTti/2gbpZsyU5vQLOJOPNcDgcfhxKkog3S55YuAbpJY6S7O1LFp6PHjjfvEXo+flZftZklt4jFWMFfbi+Wq0fIKZnYZJxmqxhJL1k4duEzWkM2YauodRfsTXlIUvOR3AWlu4/PUGqyKr8kgUj5z9S8Uf+BhpAWg8L0ZRycFSs2GfYPlMxQYWgjbgMI5gFjh3cKio17LV1q+sBKLe3pkJQpWzjbyDNQpY4iqwSVA/qCFA3BFIno4p6FpyP8EhK4X6ex7eQno8UVdPHhmmNpAyiCNIl3Ce1SsG2rZvaqDN95XFBUx5C1lYV6lXhRWi3+7pCLzLrXLv/c+fT9yvfXf5z4y5d6d23/xNUqPpslnW4Th0SQTtZn8WN6+8AnqfpmqprY31sGpZBUEfbZz5hCadrPqWHz0VgW+GPSA+hrFlnggnObeMcty2hf3mEHa+/zfaQamErmq2/SdqAfss+qZMvRBqYPpdXDkk+lO3ymXnrJoaqgRXDLG9z5WF5kqXwkVOhM83yc+MKuvhfa93r0cCNpNSVD9nd75NwAA//9QSwcIPF3XvMcFAAC0FwAAUEsBAhQAFAAIAAgAAAAAADxd17zHBQAAtBcAAAsAAAAAAAAAAAAAAAAAAAAAADE5MDAwNDEueG1sUEsFBgAAAAABAAEAOQAAAAAGAAAAAA==",
    "statusCode": 200
}
```
You will then receive the files in zip archive if everything goes smoothly. In this case just 1 file

</details>
<details><summary>credential/add</summary>

#### /credential/add - add operator credentials
##### Request body:
```json
{
    "JWT": "identity admin application token",
    "credentials": {
        "authPhrase": "omniva authphrase",
        "omnivaUrl": "omniva data exchange url",
        "vendor": "maventa vendorAPIKey",
        "user": "maventa user secret",
        "company": "maventa company secret"
    }
}
```

##### Responses with explanation:

 * Successful response
 ```json
{
    "message": "operator credentials were added",
    "statusCode": 200
}
```
</details>

## How to Run & Build
<details><summary>Run</summary>

### Run
Simple execution.

1) Go to project's directory in your terminal and execute
```sh
$ sudo apt install make
$ make #now you can see the commands you can use
$ make setup 
```
If asked enter root password of mysql-server

2) Install database using our script
```console
$ sudo mysql
> source <YOUR_FULL_PATH_TO>/e-invoices/app/db/scripts/Dump20190807.sql
```

3) Set up right user
In the configuration file `config.json` there are fields `DbUser`, `DbPassword` that you can set.

```console
sudo mysql
CREATE USER 'DbUser'@'localhost' IDENTIFIED BY 'DbPassword';
GRANT ALL PRIVILEGES ON invoice_db. * TO 'DbUser'@'localhost';
FLUSH PRIVILEGES;
````

4) Install MariaDB (link is for Ubuntu 18.04 LTS)
https://downloads.mariadb.org/mariadb/repositories/#distro=Ubuntu&distro_release=bionic--ubuntu_bionic&mirror=ukfast&version=10.4

5) Run the application
```console
$ make run
```
Command above will generate logs under logs/ directory. I use this commands to monitor the system:
```sh
$ tail -f logs/frontend.log
$ tail -f logs/backend.log
$ tail -f logs/backend_error.log
```
Note: logs will be terminated when you execute `make clean` command.
Make sure to stop your program after using it with 
```sh
$ make stop
$ make clean
```
</details>
<details><summary>Build</summary>

### Build
For this you will need Go and make.
After successful installation go to project's directory in your terminal and execute

```sh
$ make build
$ make run
```

</details>

### Dependency versions

* MariaDB: 10.4.7
* Go: 1.13

## Changelog

<details><summary>v3</summary>
    <details><summary>3.2.0</summary>
* moved Maventa vendor key to config.json 
    </details>
    <details><summary>3.1.0</summary>
* no provider credentials were added error for /send
* do not check lastModified on getExportRecords from DB. 
* e-invoice: calculate itemSum using amount*price
* e-invoice: addition: calculate sum: amount * price * disc / 100
    </details>
    <details><summary>3.0.0</summary>
* Omniva data exchange URL in UI
* Omniva data exchange URL added to credential/add request body
    </details>
</details>
<details><summary>v2</summary>
    <details><summary>2.0.0</summary>

* download zip in base64 encoding
* better finvoice field mapping
    </details>
</details>
<details><summary>v1</summary>
    <details><summary>1.1.0</summary>

* access to UI from back office
* better finvoice field mapping
    </details>

    <details><summary>1.0.7</summary>

* 404 error page
    </details>
    <details><summary>1.0.4</summary>

* fixed stopping scripts
* commented install_dependencies script
* removed customer's IBAN validation
* Multiple provider handling: if only Omniva is set up -> send all there, same for Maventa. If both -> check address to select provider for invoice

    </details>
    <details><summary>1.0.3</summary>

* gitlab-ci

    </details>
    <details><summary>1.0.2</summary>

* download invoices as XML files
* pretty loader for UI
* new service responses with more explanation for Send endpoint
* clientID of type string in all requests
* Maventa: invoices are possible to send by e-mail / e-invoicing address / post
* Maventa: possible to send invoices to customers of type person

    </details>
    <details><summary>1.0.1</summary>

* API default port: 80
* disabled host check
* Added database IP, port to configuration. Production UI port is 443.
* Migrated from MYSQL to MariaDB

    </details>
    </details>
