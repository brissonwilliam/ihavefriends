import React, {Component} from 'react';
import ButtonCounter from './buttonCounter';
import Billing from './billing';
import "./dashboard.css"


class Dashboard extends Component {
    render() {
        return (
            <div className="Dashboard-Body">
                <ButtonCounter />
                <br></br>
                <hr></hr>
                <Billing />
            </div>
        );
    }
}

export default Dashboard;