import { Component } from "react";
import Button from "react-bootstrap/esm/Button";
import cfg from "../cfg";

class UserSelection extends Component {
    constructor(props) {
        super(props);
        this.state = {
            userSelection: [""]
        }
    }

    handleClick = (e) => {
        this.props.onUserSelected(e.target.id);
        e.preventDefault();
    }

    componentDidMount() {
        fetch(cfg.BACKEND_HOST + '/publicUsers')
            .then(res => res.json())
            .then(res => this.updateState(res));
    }

    updateState(resJSON) {
        this.setState({
            userSelection: resJSON
        });
    }

    render() {
        let userSelectionHTML = this.state.userSelection.map((val, index) =>  
            <Button id={val} key={index} className="btn btn-danger btn-lg mt-1 mb-1 m-2" role="button" data-bs-toggle="button" onClick={this.handleClick}>{val}</Button>
        );

        return (
            <div>
                {userSelectionHTML}
            </div>
        );
    }
}

export default UserSelection;