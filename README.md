<div id="top"></div>


<!-- ABOUT THE PROJECT -->
# About The Project

ihavefriends is a personal project made just **for fun** for my friends and I. Most of the content and functionalities
of the website are not to be taken seriously. They are not meant to be understood or used by the general public
as they are built upon inside jokes and fun memories of the clique.

A live build of this project is in production [here](https://sourpusss.com). Obviously, you will need user credentials to navigate the site though. Either contact me or build the app and try for yourself :)

### Built With

* [Go](https://go.dev/)
* [NodeJS](https://nodejs.org/en/)
* [React.js](https://reactjs.org/)
* [Bootstrap](https://getbootstrap.com)
* [WebSockets](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket)


## TODO List

This project is currently in development. Here are some important missing parts of the project.

### Technical tasks
- UNIT TESTS
- MORE UNIT TESTS, anything under 90% is a shame
- Build with docker
- GitHub Actions CI to deploy on PROD

### Features to implement
- Refactor code so storage is seperated from components
- Users can change their password
- User creation (admin only) (missing frontend only)
- User deletion (admin only)


<!-- GETTING STARTED -->
## Getting Started

### Prerequisites
* npm
* go

### Installation

1. Clone the repo
   ```
   git clone git@github.com:brissonwilliam/ihavefriends.git
   ```
2. Install npm dependencies
   ```
   cd frontend
   npm install
   ```
3. Install go dependencies
   ```
   cd backend
   go mod download
   ```

### Building and running the project
Frontend dev
```
npm start
```

Frontend build for production. This creates js chunks and cached static files that can be served in production. Also bakes in the PROD config (hosts, not secrets, duh!)
```
npm run build -- --production
```

Backend build
```
go build
```


<!-- LICENSE -->
## License

Distributed under the GNU-GPLv3 License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#top">back to top</a>)</p>
