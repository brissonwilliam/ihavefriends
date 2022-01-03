import React, {Component} from 'react';
import Button from 'react-bootstrap/Button';
import "./buttonCounter.css"
import { IsJWTValid } from './login';

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
    }

    componentDidMount() {
        fetch('/api/bonneFeteRaph')
            .then(res => res.json())
            .then(res => this.updateStateFromResponse(res));
    }

    updateStateFromResponse(res) {
        this.setState({
            analytics: res
        });
    }

    handleClick(e) {
        var userId = IsJWTValid
        const requestOptions = {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify()
        };
        fetch('/api/bonneFeteRaph', requestOptions)
            .then(res => res.json())
            .then(res => this.updateStateFromResponse(res));
    }

    render() {
        let analyticsHtml = this.state.analytics.totalByUsers.map(userTotal =>  
            <div className="row justify-content-center">
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

                <Button className="m-3 mt-md-1 Btn-Counter-Btn" onClick={this.handleClick} variant="dark">+1</Button>              

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

            </div>
        );
    }
}

export default ButtonCounter;