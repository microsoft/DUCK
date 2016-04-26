# DUCK Application

This is a Gulp-powered build system with these features:

- Sass compilation and prefixing
- JavaScript concatenation
- Built-in BrowserSync server
- For production builds:
  - CSS compression
  - JavaScript compression
  - Image compression

## Installation

To use this template, your computer needs:

- [NodeJS](https://nodejs.org/en/) (0.12 or greater)
- [Git](https://git-scm.com/)

This template can be installed with the Foundation CLI, or downloaded and set up manually.

### Setup

First clone with Git:

```bash
git clone https://github.com/Metaform/duck duck
```

Then open the folder in your command line, and install the needed dependencies:

```bash
cd duck
npm install
bower install
```

Finally, run `npm start` to run Gulp. Your finished site will be viewable at:

```
http://localhost:8000
```

To create compressed, production-ready assets, run `npm run build`.
