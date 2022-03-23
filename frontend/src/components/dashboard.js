import React, {Component} from 'react';
import ButtonCounter from './buttonCounter';
import Billing from './billing';
import "./dashboard.css"


class Dashboard extends Component {
    render() {
        return (
            <div className="Dashboard-Body">
                <Billing />
                <br></br>
                <ButtonCounter />
            </div>
        );
    }
}

export default Dashboard;