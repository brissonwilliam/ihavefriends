import React, { Component } from 'react';
import  { Navigate } from 'react-router-dom';
import { useCookies } from 'react-cookie';
import './login.css';

export default function Login() {
    const [cookies, setCookie] = useCookies(['jwt', 'jwtExpiration']);

    if (IsJWTValid(cookies.jwt, cookies.jwtExpiration)) {
        return (<Navigate replace to="/dashboard"/>);
    }

    class Login extends Component {
        constructor(props) {
            super(props);
            this.state = {
                isLoggedIn: false,
                showInvalidCredentials: false
            }
        }
        
        handleSubmit = (e) => {
            e.preventDefault();

            let credentials = {username: "fake", password: "backend"};
            this.postAuthenticate(credentials)
                .then(res => {
                    if (!res.ok) {
                        this.setState({
                            showInvalidCredentials: true
                        });
                        throw new Error("Invalid crendentials");
                    }   
                    return res.json();                 
                })
                .catch(err => console.log(err))
                .then(jsonRes => {
                    let userInfo = {
                        jwt: jsonRes.jwt,
                        jwtExpiration: jsonRes.exp,
                        name: jsonRes.username
                    };
                    this.setUserInfo(userInfo);
                }); 
        }

        async setUserInfo(userInfo) {
            let cookieOptions = {sameSite: true};
            setCookie('jwt', userInfo.jwt, cookieOptions);
            setCookie('jwtExpiration', userInfo.jwtExpiration, cookieOptions);
            this.setState({
                isLoggedIn: true
            });
        }

        postAuthenticate(credentials) {
            let options = {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(credentials)
            };
            return fetch('/api/auth', options)
        }
       
        render() {
            if (this.state.isLoggedIn) {
                return (<Navigate replace to="/dashboard"/>);
            }
    
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
                    <div className={this.state.showInvalidCredentials ? "" : ".d-none"}>
                        <label>Invalid credentials</label>
                    </div>
                    <div>
                        <button type="submit">Submit</button>
                    </div>
                    </form>
                </div>
            );
        }
        
    }

    return <Login/>
};

// IsJWTValid checks for jwt in cookies to determine if
// a new auth should be made.
// Not validating on the server side to make things lighter
// If user messes with his token or expiration the server will retun 401
// on other endpoints anyways
export function IsJWTValid(jwt, jwtExpiration) {
    console.log(jwt)
    console.log(jwtExpiration)

    if (jwt == null || jwtExpiration == null || jwt === "" || jwtExpiration === "" || jwt === "undefined" || jwtExpiration === "undefined") {
        return false;
    }

    let expirationDate = new Date(Date.parse(jwtExpiration));
    let now = new Date();

    if (now > expirationDate) {
        console.log("expired!");
        return false
    }

    return true
};