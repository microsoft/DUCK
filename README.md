# DUCK Application

This is a Gulp-powered build system with the following features:

- Sass compilation and prefixing
- JavaScript transpilation based on Babel and concatenation
- Go compilation
- Dynamic browser reloading using BrowserSync
- For production builds:
  - CSS compression
  - JavaScript compression
  - Image compression

## Installation

The project build requires:

- [Git](https://git-scm.com/)
- [Go] (https://golang.org/)  (1.7 or later)
- [NodeJS](https://nodejs.org/en/) (0.12 or greater, 4.4.7 LTS recommended)
- [CouchDB](http://couchdb.apache.org/) (1.6 or greater)
- [SWI Prolog] (http://swi-prolog.org) (Development Release 7.3)

### Setup

First clone the project repository (make sure it is under the GOPATH directory):

```bash
git clone https://github.com/Microsoft/DUCK DUCK
```

Then, from the cloned directory, install the required dependencies:

```bash
cd DUCK
npm install
npm install -g bower
npm install -g gulp
bower install
```
Update go dependencies by executing 'go get' in the DUCK directory:

```bash
go get github.com/carneades/carneades-4
go get  gopkg.in/yaml.v2
```

Make sure Couch DB is running.

Finally, run `npm start` to execute the build. The application will be accessible at:

```
http://localhost:8000
```
Dynamic reloading will be enabled. Both frontend (Javascript, CSS, HTML) and backend (go) assets are watched for changes, which will automatically trigger an
application update.  

To create compressed, production-ready assets, run `npm run build`.


### Configuration

This project reads its configuration from the file backend/configuration.json, environment variables and also command-line flags. 
The following precendence order is used. Each item takes precedence over the item below it:


-flag
-env
-config
-default

#### Default
The default configuration has these values:

```yaml
  database: 
      location: "http://127.0.0.1"
      port: 5984
      name: "duck"
	  username: ""
	  password: ""
  jwtkey: "secret"
  webdir: "/src/github.com/Microsoft/DUCK/frontend/dist"
  rulebasedir: "/src/github.com/Microsoft/DUCK/RuleBases"
  gopathrelative: true
  loadtestdata: false
```
#### env
The environment variable names are prefixed with DUCK_ and all uppercase. Fields in the database object are referenced using the `.` operator, e.g. 
`DUCK_DATABASE.NAME`.

#### flags
The flags are handled in the go standard way described in https://golang.org/pkg/flag/. 
Main points are:

>Command line flag syntax:
>```
-flag
-flag=x
-flag x  // non-boolean flags only
```
>One or two minus signs may be used; they are equivalent. The last form is not permitted for boolean flags because the meaning of the command
>
>`cmd -x *`
>
>will change if there is a file called 0, false, etc. You must use the -flag=false form to turn off a boolean flag.
>
>Integer flags accept 1234, 0664, 0x1234 and may be negative. Boolean flags may be:
>
>`1, 0, t, f, T, F, true, false, TRUE, FALSE, True, False`
>
>It is *not* possible to configure the database via flags.
