import React, { Component } from 'react';
import  { Navigate } from 'react-router-dom';
import { useCookies } from 'react-cookie';
import './login.css';
import UserSelection from './userSelection';
import Button from 'react-bootstrap/esm/Button';
import cfg from '../cfg';


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
                showInvalidCredentials: false,
                selectedUser: "",
            }
        }
        
        handleSubmit = (e) => {
            e.preventDefault();

            let credentials = {username: this.state.selectedUser, password: e.target.elements.password.value};
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
                .then(jsonRes => {
                    let userInfo = {
                        jwt: jsonRes.jwt,
                        jwtExpiration: jsonRes.jwtExpiration,
                    };
                    console.log(userInfo)
                    this.setUserInfo(userInfo);
                })
                .catch(err => console.log(err)); 

                document.getElementById('password-input').value = "";
        }

        async setUserInfo(userInfo) {
            let expirationD = new Date(userInfo.jwtExpiration);
            let cookieOptions = {sameSite: true, expires: expirationD};
            setCookie('jwt', userInfo.jwt, cookieOptions);
            setCookie('jwtExpiration', userInfo.jwtExpiration, cookieOptions);
            this.setState({
                isLoggedIn: true,
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
            return fetch(cfg.BACKEND_HOST + '/auth', options)
        }

        // this callback is called from UserSelection child
        handleUserSelected = (username) => {
            this.setState({
                selectedUser: username
            });
        }
       
        render() {
            if (this.state.isLoggedIn) {
                return (<Navigate replace to="/dashboard"/>);
            }

            return(
                <div className="col-12 Login-Wrapper">
                    <h1 className="p-2">Tequilla?</h1>

                    <form onSubmit={this.handleSubmit}> 
                        <UserSelection onUserSelected={this.handleUserSelected}/>
                        <div className="p-3">
                            <div>
                                <p className={this.state.selectedUser !== "" ? "text-muted" : "d-none"}>Salut {this.state.selectedUser}</p>
                            </div>                        
                            <div>
                                <p className="mb-1">Password</p>
                                <input id="password-input" name="password" type="password"/>
                            </div>
                            <div>
                                <label className={this.state.showInvalidCredentials ? "text-danger mt-2" : "d-none"}>Nope! Try again!</label>
                            </div>
                        </div>
                       
                        <div>
                            <Button className="btn btn-light m-3" type="submit">Submit</Button>
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
    if (jwt == null || jwtExpiration == null || jwt === "" || jwtExpiration === "" || jwt === "undefined" || jwtExpiration === "undefined") {
        return false;
    }

    let expirationDate = new Date(Date.parse(jwtExpiration));
    let now = new Date();

    if (now > expirationDate) {
        return false
    }

    return true
};