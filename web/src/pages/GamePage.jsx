import React from "react";
import "bootstrap/dist/css/bootstrap.css";
import Row from "../components/Row";
import "../styles/App.scss";
import {
    Alert,
    Jumbotron,
} from "reactstrap";
import {
    Link,
    withRouter,
} from "react-router-dom";
import {AIScore, patterns, symbolsMap} from "../board/board";
import Header from "../components/Header";

class GamePage extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            active: true,
            alertShow: false,
            alertText: '',
            alertType: '',
            boardState: new Array(9).fill(2),
            connAlertShow: false,
            connAlertText: '',
            connAlertType: '',
            gameId: this.props.gameCode,
            isInitiator: false,
            lastMoveUserId: '',
            mode: 'AI',
            opponentId: '',
            shareButtonShow: false,
            turn: 0,
            userId: '',
            ws: null,
        };
        this.handleNewMove = this.handleNewMove.bind(this);
        this.handleReset = this.handleReset.bind(this);
        this.processBoardState = this.processBoardState.bind(this);
        this.makeAIMove = this.makeAIMove.bind(this);
    }

    processBoardState = () => {
        const {
            boardState,
            mode,
            turn,
        } = this.state;
        let won = false;

        patterns.forEach(pattern => {
            const firstMark = boardState[pattern[0]];

            if (firstMark !== 2) {
                const marks = boardState.filter((mark, index) => {
                    return pattern.includes(index) && mark === firstMark; //looks for marks matching the first in pattern's index
                });

                if (marks.length === 3) {
                    pattern.forEach(index => {
                        const id = index + "-" + firstMark;
                        document.getElementById(id).parentNode.style.background = "#d4edda";
                    });
                    won = true;
                    this.setState({
                        alertText: `${String.fromCharCode(symbolsMap[marks[0]][1])} выиграли!`,
                        alertShow: true,
                        alertType: 'success',
                        active: false,
                    });
                }
            }
        });

        if (!boardState.includes(2) && !won) {
            this.setState({
                alertText: `Конец игры - ничья`,
                alertShow: true,
                alertType: 'warning',
                active: false
            });
        } else if (mode === 'AI' && turn === 1 && !won) {
            this.makeAIMove();
        }
    }

    makeAIMove = () => {
        const {boardState} = this.state;
        const empty = [], scores = [];

        boardState.forEach((mark, index) => {
            if (mark === 2)
                empty.push(index);
        });

        // MiniMax
        empty.forEach(index => {
            let score = 0;
            patterns.forEach(pattern => {
                if (pattern.includes(index)) {
                    let xCount = 0, oCount = 0;
                    pattern.forEach(p => {
                        if (boardState[p] === 0) xCount += 1;
                        if (boardState[p] === 1) oCount += 1;
                        score += p === index ? 0 : AIScore[boardState[p]];
                    });
                    if (xCount >= 2) score += 10;
                    if (oCount >= 2) score += 20;
                }
            });
            scores.push(score);
        });

        let maxIndex = 0;
        scores.reduce((maxVal, currentVal, currentIndex) => {
            if (currentVal >= maxVal) {
                maxIndex = currentIndex;
                return currentVal;
            }
            return maxVal;
        });

        this.handleNewMove(empty[maxIndex]);
    }

    handleReset = e => {
        if (e) e.preventDefault();
        this.setState({
            boardState: new Array(9).fill(2),
            turn: 0,
            active: true,
            alertShow: false,
            alertType: '',
            alertText: '',
        });
    }

    handleNewMove = (id) => {
        this.setState(
            prevState => {
                return {
                    boardState: prevState.boardState
                        .slice(0, id)
                        .concat(prevState.turn)
                        .concat(prevState.boardState.slice(id + 1)),
                    turn: (prevState.turn + 1) % 2
                };
            },
            () => {
                this.processBoardState();
            }
        );
    }

    render() {
        const {
            active,
            alertShow,
            alertText,
            alertType,
            boardState,
            turn,
        } = this.state;
        const rows = [];

        for (let i = 0; i < 3; i++)
            rows.push(
                <Row
                    key={i}
                    row={i}
                    boardState={boardState}
                    onNewMove={this.handleNewMove}
                    active={active}
                />
            );
        return (
            <div>
                <Jumbotron
                    className="container"
                >
                    <Header/>
                    <div>
                        {this.state.mode === "AI"
                            ?
                            <Link id="btn-reset-game" className="btn btn-block btn-primary" to="#" onClick={this.handleReset} size="sm">Сбросить игру</Link>
                            : null
                        }
                    </div>
                    <br/>
                    <p>Очередь за {String.fromCharCode(symbolsMap[turn][1])}</p>
                    <br/>
                    <div className="board">
                        {rows}
                    </div>
                    <div className="alert-container">
                        <Alert color={alertType} show={alertShow} isOpen={alertShow}>
                            {alertText}
                        </Alert>
                    </div>
                </Jumbotron>
            </div>
        );
    }
}

export default withRouter(GamePage);
