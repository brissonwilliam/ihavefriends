import { Component } from "react";
import Button from "react-bootstrap/esm/Button";

const userSelection = ["willow", "jowanne", "rose", "lauree", "stayo", "boeuf", "raph", "goulat"];

class UserSelection extends Component {
    handleClick = (e) => {
        this.props.onUserSelected(e.target.id);
        e.preventDefault();
    }

    render() {
        let userSelectionHTML = userSelection.map((val, index) =>  
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