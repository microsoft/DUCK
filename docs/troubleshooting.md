# Trouble Shooting Guide

## CouchDB
- As of June 2017, version 2.0.0 of CouchDB does not work well with Windows 10 Creators Update. Use an older version of CouchDB.
- Make sure CouchDB is indeed started, manually or automatically, depending on the choice you selected during setup. If you suspect CouchDB is not started, try accessing http://127.0.0.1:5984 in your browser. The browser should render a text page that looks like the following:
`{"couchdb":"Welcome","version":"1.6.1","vendor":{"name":"The Apache Software Foundation"}}`
    - Another clue CouchDB has started is that you will see the process “erl.exe” in Windows task list.
    - For manual startup of CouchDB 2 on Windows, you may need to run the couchdb.cmd from an Admin cmd window (running in a plain cmd window won’t work even if you are logged in as admin)
- If CouchDB is installed as a Windows service, it will appear as “Apache CouchDB” in the Services list.

## Go & Duck.exe

- No need to install Go binaries if you are not planning to build DUCK
-   _If you have the GO Programming Language installed and/or the  GOPATH environment variable defined_  the output log will say `Found GOPATH, will use gopath for relative paths.` If this happens you might have to change the configuration Keys `webdir` and `rulebasedir`. The path to the configuration will be `$GOPATH/src/github.com/Microsoft/DUCK/backend/configuration.json`. If these have a relative path value the program will look for the file relative to the GOPATH. _Only if a GOPATH is set._ The program will print the full path it is using to the command line.
