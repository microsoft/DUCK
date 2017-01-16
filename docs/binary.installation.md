DUCK Installation
===================

DUCK is distributed as an execuable binary. To install:

1. send your GitHub username to babakj@microsoft.com to be added to the project, if you don't have one create a new one: https://github.com/join
1. install CouchDB: http://couchdb.apache.org/#download
2. install swi-prolog: http://www.swi-prolog.org/download/devel
1. Logon to GitHub 
3. extract the zip from https://github.com/Microsoft/DUCK/releases. You can only access this page when you are logged in.
4. run backend.exe from where the binaries were extracted. One example, on Windows machine, would be: C:\Users\\[username]\Downloads\duck_win
7. When Windows Firewall come up, allow backend.exe to run over all possible network types.
5. open a browser window and load up: [http://localhost:3000](http://localhost:3000)
5. A manual can be found here: https://github.com/Microsoft/DUCK/blob/master/docs/usermanual.md

**Updating**: Since this app has only the two dependencies CouchDB and prolog outside of the downloadable zip it is *completely safe to  delete the previous downloaded version and use the files from this one* without any side effects. It is also possible to keep the old files and run the new version from another place without problems.

_If you have the GO Programming Language installed and/or the  GOPATH environment variable defined_ you might have to change the configuration Keys `webdir` and `rulebasedir`. If these have a relative path value the program will look for the file relative to the GOPATH. _Only if a GOPATH is set._ The program will print the full path it is using to the command line.
