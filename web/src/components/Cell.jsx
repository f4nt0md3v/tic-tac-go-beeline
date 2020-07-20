import React from "react";
import {symbolsMap} from "../board/board";

class Cell extends React.Component {
    constructor(props) {
        super(props);
        this.handleNewMove = this.handleNewMove.bind(this);
    }

    handleNewMove(e) {
        if (this.props.marking === 2 && this.props.active)
            this.props.onNewMove(parseInt(e.target.id));
    }

    render() {
        return (
            <div className="col" onClick={this.handleNewMove}>
                <div className={symbolsMap[this.props.marking][0]} id={this.props.id}>
                    {String.fromCharCode(symbolsMap[this.props.marking][1])}
                </div>
            </div>
        );
    }
}

export default Cell;
