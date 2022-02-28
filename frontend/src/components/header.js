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
            <div className="Header">
                <div className="row d-flex pt-2 pb-2 ms-2 me-2">
                    <div className="col-10 text-start">
                        <span>SP MASTER</span>
                    </div>
                    <div className="col-2 d-flex justify-content-end">
                        <Button className={isLoggedIn? "btn Header-logout-btn": "d-none"} onClick={this.handleLogoutClick}>Logout</Button>
                    </div>
                </div>
            </div>);
        }
    }


    return <Header/>
}

