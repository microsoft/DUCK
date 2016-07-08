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

- GO (1.6.2 or later)
- [NodeJS](https://nodejs.org/en/) (0.12 or greater)
- [Git](https://git-scm.com/)

### Setup

First clone the project repository (make sure it is under the GOPATH directory):

```bash
git clone https://github.com/Microsoft/DUCK duck
```

Then, from the cloned directory, install the required dependencies:

```bash
cd duck
npm install
npm install -g bower
npm install -g gulp
bower install
```

Finally, run `npm start` to execute the build. The application will be accessible at:

```
http://localhost:8000
```
Dynamic reloading will be enabled. Both frontend (Javascript, CSS, HTML) and backend (go) assets are watched for changes, which will automatically trigger an
application update.  

To create compressed, production-ready assets, run `npm run build`.
