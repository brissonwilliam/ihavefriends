import React, {Component} from 'react';
import ButtonCounter from './buttonCounter';
import "./dashboard.css"

class Dashboard extends Component {
    render() {
        return (
            <div className="Dashboard-Body">
                <div>
                    <ButtonCounter />
                </div>
            </div>
        );
    }
}

export default Dashboard;