import React, {Component} from 'react';
import Button from 'react-bootstrap/Button';
import "./buttonCounter.css"

class ButtonCounter extends Component {
    constructor() {
        super();
        this.state = {
            c: 0
        };
        this.handleClick = this.handleClick.bind(this);
    }

    componentDidMount() {
        fetch('/api/bonneFeteRaph')
            .then(res => res.json())
            .then(resJson => this.updateStateFromResponse(resJson));
    }

    updateStateFromResponse(json) {
        this.setState({
            c: json.nbBonneFete
        })
    }

    handleClick(e) {
        const requestOptions = {
            method: 'POST',
        };
        fetch('/api/bonneFeteRaph', requestOptions)
            .then(res => res.json())
            .then(resJson => this.updateStateFromResponse(resJson));
    }

    render() {
        return (
            <div className="d-grip col-12">
                <div className="Button-Counter-Container text-wrap fs-1 mt-5 mb-5 lh-sm">
                    <p className="m-0">Bonne fête raph!</p>
                    <p className="m-0">Tu as 
                        <span className="Button-Counter-Count m-0"> {this.state.c} </span>
                        ans!
                    </p>
                </div>  

                <Button className="mb-4 Btn-Counter-Btn" onClick={this.handleClick} variant="dark">+1</Button>              
            </div>
        );
    }
}

export default ButtonCounter;