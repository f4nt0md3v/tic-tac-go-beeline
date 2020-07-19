import React  from "react";
import Cell from "./Cell";

class Row extends React.Component {
    render() {
        const cells = [];
        for (let i = 0; i < 3; i++) {
            const id = this.props.row * 3 + i;
            const marking = this.props.boardState[id];
            cells.push(
                <Cell
                    key={id + "-" + marking}
                    id={id + "-" + marking}
                    marking={marking}
                    onNewMove={this.props.onNewMove}
                    active={this.props.active}
                />
            );
        }
        return <div className="row">{cells}</div>;
    }
}

export default Row;
