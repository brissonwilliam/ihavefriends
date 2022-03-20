import React, {Component} from 'react';
import Button from 'react-bootstrap/Button';
import { useCookies } from 'react-cookie';
import { w3cwebsocket as W3CWebSocket } from "websocket";
import "./buttonCounter.css";

const BACKEND_HOST = process.env.REACT_APP_BACKEND_HOST;

export default function ButtonCounter() {

    const [cookies, setCookie] = useCookies(['jwt', 'jwtExpiration']);

    var wsUrl = 'wss://' + BACKEND_HOST + '/api/bonneFete/ws?token=' + cookies.jwt
    var ws;

    class ButtonCounter extends Component {
        constructor() {
            super();
            this.state = {
                analytics: {
                    total: 0,
                    totalByUsers: [
                        {
                            name: "",
                            count: 0
                        }
                    ]
                }
            };
            this.handleClick = this.handleClick.bind(this);
            this.handleClickReset = this.handleClickReset.bind(this);
        }

        componentDidMount() {
            ws = new W3CWebSocket(wsUrl);
            ws.onopen = () => {
                console.log("websocket connected");
            };
            ws.onmessage = (message) => {
                const dataFromServer = JSON.parse(message.data);
                this.updateStateFromResponse(dataFromServer)
            };
            ws.onclose = (e) => {
                console.log("websocket closed");
            };

            const requestOptions = {
                method: 'GET',
                headers: {'Authorization': 'Bearer ' + cookies.jwt}
            }
            fetch('/api/bonneFete', requestOptions)
                .then(res => res.json())
                .then(res => this.updateStateFromResponse(res));
        }

        componentWillUnmount() {  
            if (ws != null) {
                ws.close();
            }
        }
    
        updateStateFromResponse(res) {
            this.setState({
                analytics: res
            });
        }
    
        handleClick(e) {
            e.preventDefault();
            const requestOptions = {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': 'Bearer ' + cookies.jwt
                },
                body: JSON.stringify()
            };
            fetch('/api/bonneFete', requestOptions)
                .then(res => res.json())
                .then(res => this.updateStateFromResponse(res));
        }

        handleClickReset(e) {
            e.preventDefault();
            if (!window.confirm('Voulez-vous vraiment remettre votre compteur à 0?')) {
                return;
            }
            const requestOptions = {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': 'Bearer ' + cookies.jwt
                },
                body: JSON.stringify()
            };
            fetch('/api/bonneFete/reset', requestOptions)
                .then(res => res.json())
                .then(res => this.updateStateFromResponse(res));
        }

        render() {
            let analyticsHtml = this.state.analytics.totalByUsers.map((userTotal, index) =>  
                <div key={index} className="row justify-content-center">
                    <div className="col-6 col-lg-4">
                        <p className="m-0 m-md-1 text-start">{userTotal.name}</p>
                    </div>
                    <div className="col-6 col-lg-4">
                        <p className="m-0 m-md-1 text-center">{userTotal.count}</p>
                    </div>
                </div>
            );
    
    
            return (
                <div className="col-12">
                    <div className="Button-Counter-Container text-wrap fs-1 mt-5 mb-0 lh-sm">
                        <p className="m-0">Bonne fête raph!</p>
                        <p className="m-0">Tu as 
                            <span className="Button-Counter-Count m-0"> {this.state.analytics.total} </span>
                            ans!
                        </p>
                    </div>  

                    <Button className="col-8 col-lg-3 m-3 mt-md-1 Btn-Counter-Btn" onClick={this.handleClick} variant="dark">+1</Button>

                    <div className="container mt-md-5 mt-3">
                        <div className="row justify-content-center">
                            <div className='col-md-5 col-12 border border-2 rounded-3 align-self-center'>
                                <h2 className="p-3">Les meilleurs souhaiteurs de fête</h2>
                                <div className='container pb-3'>
                                    {analyticsHtml}
                                </div>
                            </div>
                        </div>
                    </div>

                    <br/>
                    <div>
                        <Button className="m-3" onClick={this.handleClickReset} variant="dark">Reset</Button>
                    </div>
                </div>
            );
        }
    }
    return <ButtonCounter/>  
};