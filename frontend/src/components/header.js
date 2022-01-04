import "./header.css"
import { Component } from "react"
import Button from 'react-bootstrap/Button';
import { useCookies } from "react-cookie";
import {IsJWTValid} from './login';


export default function Header() {
    const [cookies, setCookie, deleteCookie] = useCookies(['jwt', 'jwtExpiration']);

    class Header extends Component{
        constructor() {
            super();
            this.state = {
                redirectToLogin: false
            }
            this.handleLogoutClick = this.handleLogoutClick.bind(this);
        }
    
    
        handleLogoutClick(e) {
            deleteCookie('jwt');
            deleteCookie('jwtExpiration');
            this.setState({redirectToLogin: true});
        }
    
        render() {
            let isLoggedIn = IsJWTValid(cookies.jwt, cookies.jwtExpiration);

            if (this.state.redirectToLogin) {
                return (window.location.reload());
            }

            return (
            <header className="Header">
                <div className="row ps-4 pe-3 pt-1 pb-1 me-0 ms-0">
                    <div className="col-9 text-start">
                        <span>SP MASTER</span>
                    </div>
                    <div className="col-1"></div>
                    <div className={isLoggedIn ? "col-2" : "d-none"}>
                        <Button className="btn float-end p-1 Header-logout-btn" onClick={this.handleLogoutClick}>Logout</Button>
                    </div>
                </div>
            </header>);
        }
    }


    return <Header/>
}

