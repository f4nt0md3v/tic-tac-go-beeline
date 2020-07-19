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
import Share from "../components/Share";

const symbolsMap = {
    2: ["marking", "32"],
    0: ["marking marking-x", 9587], // represents x
    1: ["marking marking-o", 9711], // represents o
};

const patterns = [
    //horizontal
    [0, 1, 2],
    [3, 4, 5],
    [6, 7, 8],
    //vertical
    [0, 3, 6],
    [1, 4, 7],
    [2, 5, 8],
    //diagonal
    [0, 4, 8],
    [2, 4, 6]
];

const AIScore = {2: 1, 0: 2, 1: 0};

class GamePage extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            active: true,
            alertShow: false,
            alertText: '',
            alertType: '',
            boardState: new Array(9).fill(2),
            gameId: this.props.gameCode,
            mode: '',
            opponentId: '',
            shareButtonShow: false,
            turn: 0,
            userId: '',
            ws: null,
        };
        this.handleNewMove = this.handleNewMove.bind(this);
        this.handleReset = this.handleReset.bind(this);
        this.handleModeChange = this.handleModeChange.bind(this);
        this.processBoard = this.processBoard.bind(this);
        this.makeAIMove = this.makeAIMove.bind(this);
        this.connectWebSocket = this.connectWebSocket.bind(this);
        this.handleWebSocketMessage = this.handleWebSocketMessage.bind(this);
    }

    componentDidMount() {
        const path = window.location.hash;
        switch (path) {
            case "#/game/ai":
                this.setState({mode: 'AI'})
                break;
            default:
                this.setState({mode: '2P'})
                this.connectWebSocket(() => {
                    if (path === "#/game/start") {
                        if (this.state.ws && this.state.ws.readyState === WebSocket.OPEN) {
                            this.generateNewGame();
                        }
                    }
                    if (path.includes('#/game/join') && this.props.match.params.gameId !== '') {
                        this.setState({gameId: this.props.match.params.gameId}, () => {
                            if (this.state.ws && this.state.ws.readyState === WebSocket.OPEN) {
                                this.joinGame();
                            }
                        });
                    }
                });
        }
    }

    connectWebSocket = (callback) => {
        let wsUrl = 'ws:';
        if (window.location.protocol === 'https:') {
            wsUrl = 'wss:';
        }
        wsUrl += "localhost:8081/ws";
        let ws = new WebSocket(wsUrl);
        ws.onopen = () => {
            console.log("Connected to WebSocket...");
            this.setState({ws: ws}, () => {
                callback && callback();
            });
        };
        ws.onmessage = this.handleWebSocketMessage;
    };

    handleWebSocketMessage = (e) => {
        const jsonData = JSON.parse(e.data);
        alert(jsonData.command);
        if (jsonData.command) {
            switch (jsonData.command) {
                case "GENERATE_NEW_GAME":
                    if (jsonData.code === 201 && jsonData.gameInfo) {
                        this.setState({
                            gameId:          jsonData.gameInfo.gameId,
                            userId:          jsonData.gameInfo.firstUserId,
                            boardState:      jsonData.gameInfo.state.split(','),
                            shareButtonShow: true,
                        });
                    }
                    break;
                case "JOIN_GAME":
                    if (jsonData.code === 200 && jsonData.gameInfo) {
                        this.setState({
                            gameId:     jsonData.gameInfo.gameId,
                            opponentId: jsonData.gameInfo.firstUserId,
                            userId:     jsonData.gameInfo.secondUserId,
                            boardState: jsonData.gameInfo.state.split(','),
                        });
                    }
                    break;
                default:
                    break;
            }
        } else {
            if (jsonData.error) {
                console.log(jsonData)
            }
        }
    }

    generateNewGame = () => {
        const {ws} = this.state;
        if (ws || ws.readyState === WebSocket.OPEN) {
            const message = {
                command: "GENERATE_NEW_GAME"
            }
            ws.send(JSON.stringify(message))
        }
    }

    joinGame = () => {
        const {ws} = this.state;
        if (ws || ws.readyState === WebSocket.OPEN) {
            const message = {
                command: "JOIN_GAME",
                gameInfo: {
                    gameId:  this.state.gameId,
                }
            }
            ws.send(JSON.stringify(message))
        }
    }

    processBoard = () => {
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
        scores.reduce(function(maxVal, currentVal, currentIndex) {
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

    handleNewMove = id => {
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
                this.processBoard();
            }
        );
    }

    handleModeChange = e => {
        e.preventDefault();
        if (e.target.value === "AI") {
            this.setState({ mode: "AI" });
            this.handleReset(null);
        } else if (e.target.value === "2P") {
            this.setState({ mode: "2P" });
            this.handleReset(null);
        }
    }

    render() {
        const {
            active,
            alertText,
            alertType,
            alertShow,
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
                    className={"container"}
                >
                    <h3>Игра "Крестики-Нолики"</h3>
                    <hr/>
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
                        <Alert variant={alertType} show={alertShow} isOpen={alertShow}>
                            {alertText}
                        </Alert>
                    </div>
                </Jumbotron>
                {
                    this.state.mode === '2P' && this.state.shareButtonShow ? <Share gameCode={this.state.gameId}/> : null
                }
            </div>
        );
    }
}

export default withRouter(GamePage);
