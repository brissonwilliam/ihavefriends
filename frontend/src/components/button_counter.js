import React, {Component} from 'react';
import Button from 'react-bootstrap/Button';

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
            <div className="Button-Counter-Container">
                <Button className="mb-4" onClick={this.handleClick} size="lg" variant="dark">+1</Button>
                
                <p class="h4">Bonne fête raph!</p>
                <p class="h2">Tu as <span className="Button-Counter-Count"> {this.state.c} </span> ans!</p>
            </div>
        );
    }
}

export default ButtonCounter;