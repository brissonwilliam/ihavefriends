import React, {Component} from 'react';
import Button from 'react-bootstrap/Button';

class ButtonCounter extends Component {
    constructor() {
        super();
        this.state = {
            c: 0
        };
    }

    componentDidMount() {
        fetch("/api/bonneFeteRaph")
            .then(res => res.json())
            .then(resJson => this.updateStateFromResponse(resJson));
    }

    updateStateFromResponse(json) {
        this.setState({
            c: json.nbBonneFete
        })
    }

    render() {
        return (
            <div className="Button-Counter-Container">
                <Button className="mb-4" size="lg" variant="dark">+1</Button>
                
                <p class="h4">Bonne fête raph!</p>
                <p class="h2">Tu as <span className="Button-Counter-Count"> {this.state.c} </span> ans!</p>
            </div>
        );
    }
}

export default ButtonCounter;