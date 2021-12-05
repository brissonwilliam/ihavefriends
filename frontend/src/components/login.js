import React, { Component } from 'react';
import { useCookies } from 'react-cookie';
import  { Navigate } from 'react-router-dom';
import './login.css';


class Login extends Component {
    constructor(props) {
        super(props);

        this.state = {
        };
    }

    componentDidMount() {
        // check if we already have a valid auth in cookies
        if (IsLoginValid()) {
            console.log("valid login, redirecting")
        }
      }

    handleSubmit = async e => {
        e.preventDefault();

        let credentials = {username: "Fake", password: "Backend"};
        loginUser(credentials, this.handleLoginResult);
        
    }

    handleLoginResult(authResp) {
        
        // todo: redirect on success
    }

    render() {

    return(
        <div className="Login-Wrapper">
            <h1>Who are you?</h1>
            <form onSubmit={this.handleSubmit}>
            <label>
                <p>Username</p>
                <input type="text" />
            </label>
            <label>
                <p>Password</p>
                <input type="password" />
            </label>
            <div>
                <button type="submit">Submit</button>
            </div>
            </form>
        </div>
        );
    }
    
}

async function loginUser(credentials, callback) {
    let options = {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(credentials)
    };
    return fetch('/api/auth', options)
        .then(res => res.json())
        .then(callback);
}

export default Login;
  
// isLoginValid checks for information stored in cookie to determine if
// a new auth should be made.
// Not validating on the server side to make things lighter
// If user messes with his token or expiration the server will retun 401
// on other endpoints anyways
export function IsLoginValid() {
    const [cookies, setCookie] = useCookies(['user']);

    if (cookies.jwt === '' || cookies.jwtExpiration === '') {
        return false;
    }

    // FIXME: not sure if this works, need to test
    let expirationDate = new Date(Date.parse(cookies.jwtExpiration))
    if (expirationDate.getTime() > new Date().getUTCDate()) {
        return false
    }

    return true
};