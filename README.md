# API for creating agreement between 2 users using Google Docs & Drive

-------------------------------------------------------------
## About
-------------------------------------------------------------

Using this API you can create agreement between 2 users. API uses **Google Docs** and **Drive** **APIs**. To work with this API you need ***Google service account*** and enabled **Google Docs API** and **Google Drive API**. API can be used as a backend for applications or CLI tools. For the purpose of simplicity SQLITE3 db system is used. 
**Routing** logic is managed using [Gin](https://github.com/gin-gonic/gin) library. **Database** logic is managed using [GORM](https://github.com/jinzhu/gorm/) library.
Only **POST** HTTP requests work with **API**. Examples of requests can be viewed in ***example-requests.json*** file. Requests should be sent in json format, so responses are in json format too.

## Build
-------------------------------------------------------------

1. Install and configure go in your system https://golang.org/doc/install
2. Download your service account key (in json format) and place it somewhere safe. Create a .db file with following query:
```
   CREATE TABLE documents (
	doc_id TEXT UNIQUE PRIMARY KEY NOT NULL,
	doc_title TEXT NOT NULL,
	doc_url TEXT NOT NULL,
	owner1 TEXT,
	owner2 TEXT,
	signed1 INTEGER,
	signed2 INTEGER,
	date_signed1 TEXT,
	date_signed2 TEXT
);
```
3. Set environmental variables:
   - SECRET_CREDENTIALS: should be a path to your *credentials.json* file (service account key). **E.g** /temp/example/credentials.json
   - DATABASE_FILE: should be a path to your database file (sqlite .db file). **E.g** /temp/example/docs.db
4. Build a project by using go tools: ***go build***
5. (Optional) Run an executable with -port argument (default port is 5050).
-------------------------------------------------------------

## Routes
-------------------------------------------------------------

server:port/document/create - creates a document and gives write permission to a user
server:port/document/perm - gives read & write permission to indicated user
server:port/document/sign - adds a user sign in db
server:port/document/list - lists all documents of a user

## ETC
-------------------------------------------------------------
API will have more features and more convenient interface to work with. Currently it's in ***alpha*** version. If you want to collaborate or to help improve API, please create an ***issue*** or ***pull request***.

API was created for the Mobile application, but for the purposes of privacy, name of an application and full version will not be uploaded to git.

> Contact: bejanhtc@gmail.com