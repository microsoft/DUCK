# DUCK Application

This is a Gulp-powered build system with the following features:

- Sass compilation and prefixing
- JavaScript concatenation
- Built-in BrowserSync server
- For production builds:
  - CSS compression
  - JavaScript compression
  - Image compression
- Go compilation

## Installation

The project build requires:

- [NodeJS](https://nodejs.org/en/) (0.12 or greater)
- [Git](https://git-scm.com/)

### Setup

First clone the project repository:

```bash
git clone https://github.com/Metaform/duck duck
```

Then, from the cloned directory, install the required dependencies:

```bash
cd duck
npm install
bower install
```

Finally, run `npm start` to execute the build. The application will be accessible at:

```
http://localhost:8000
```
Dynamic reloading will be enabled. Both frontend (Javascript, CSS, HTML) and backend (go) assets are watched for changes, which will automatically trigger an
application update.  

To create compressed, production-ready assets, run `npm run build`.
