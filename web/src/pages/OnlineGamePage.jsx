import React from "react";
import "bootstrap/dist/css/bootstrap.css";
import Row from "../components/Row";
import "../styles/App.scss";
import {
    Alert,
    Jumbotron,
} from "reactstrap";
import {
    withRouter,
} from "react-router-dom";
import Share from "../components/Share";
import {patterns, symbolsMap} from "../board/board";
import Header from "../components/Header";

class OnlineGamePage extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            active: false,
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
            opponentId: '',
            shareButtonShow: false,
            turn: 0,
            userId: '',
            ws: null,
        };
        this.handleNewMove = this.handleNewMove.bind(this);
        this.processBoard = this.processBoard.bind(this);
        this.connectWebSocket = this.connectWebSocket.bind(this);
        this.handleWebSocketMessage = this.handleWebSocketMessage.bind(this);
        this.handleResetConnAlert = this.handleResetConnAlert.bind(this);
        this.sendMoveToServer = this.sendMoveToServer.bind(this);
    }

    componentDidMount() {
        const path = window.location.hash;
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

    connectWebSocket = (callback) => {
        let wsUrl = 'ws:';
        if (window.location.protocol === 'https:') {
            wsUrl = 'wss:';
        }
        wsUrl += window.location.host + '/ws';
        let ws = new WebSocket(wsUrl);

        ws.onopen = () => {
            console.log("Connected to WebSocket...");
            this.setState({
                ws: ws,
                connAlertShow: true,
                connAlertText: 'Установлено соединение',
                connAlertType: 'success',
            }, () => {
                setTimeout(this.handleResetConnAlert,3000);
                callback && callback();
            });
        };

        ws.onclose = () => {
            this.setState({
                connAlertShow: true,
                connAlertText: 'Не удалось установить соединение с сервером',
                connAlertType: 'danger'
            }, () => {
                setTimeout(this.handleResetConnAlert,3000);
            });
        };

        ws.onmessage = this.handleWebSocketMessage;
    };

    handleResetConnAlert = () => {
        this.setState({
            connAlertShow: false,
            connAlertText: '',
        });
    }

    handleWebSocketMessage = (e) => {
        const jsonData = JSON.parse(e.data);
        if (jsonData.command) {
            switch (jsonData.command) {
                case "GENERATE_NEW_GAME":
                    if (jsonData.code === 201 && jsonData.gameInfo) {
                        this.setState({
                            gameId:          jsonData.gameInfo.gameId,
                            userId:          jsonData.gameInfo.firstUserId,
                            boardState:      jsonData.gameInfo.state.split(',').map(Number),
                            shareButtonShow: true,
                            isInitiator:     true,
                        });
                    } else {
                        console.log(jsonData)
                    }
                    break;
                case "JOIN_GAME":
                    if (jsonData.code === 200 && jsonData.gameInfo) {
                        if (this.state.isInitiator) {
                            this.setState({
                                gameId:          jsonData.gameInfo.gameId,
                                opponentId:      jsonData.gameInfo.secondUserId,
                                boardState:      jsonData.gameInfo.state.split(',').map(Number),
                                connAlertText:   'Оппонент подключился и готов играть',
                                connAlertType:   'success',
                                connAlertShow:   true,
                                shareButtonShow: false,
                                active:          true,
                            }, () => {
                                setTimeout(this.handleResetConnAlert, 3000);
                            });
                        } else {
                            this.setState({
                                gameId:        jsonData.gameInfo.gameId,
                                opponentId:    jsonData.gameInfo.firstUserId,
                                userId:        jsonData.gameInfo.secondUserId,
                                boardState:    jsonData.gameInfo.state.split(',').map(Number),
                                connAlertShow: true,
                                connAlertText: 'Вы подключились к игре',
                                connAlertType: 'success',
                                active:        true,
                            }, () => {
                                setTimeout(this.handleResetConnAlert, 3000);
                            });
                        }
                    } else {
                        console.log(jsonData)
                        if (jsonData.error === 'sql: no rows in result set') {
                            this.setState({
                                connAlertShow: true,
                                connAlertText: 'Игра с данным кодом не найдена. Попробуйте другой код.',
                                connAlertType: 'danger'
                            }, () => {
                                setTimeout(this.handleResetConnAlert,3000);
                            });
                        }
                    }
                    break;
                case "NEW_MOVE":
                    if (jsonData.code === 200 && jsonData.gameInfo) {
                        this.setState(prevState => {
                            return {
                                gameId:          jsonData.gameInfo.gameId,
                                boardState:      jsonData.gameInfo.state.split(',').map(Number),
                                lastMoveUserId:  jsonData.gameInfo.lastMoveUserId,
                                turn:            (prevState.turn + 1) % 2
                            };
                        }, () => {
                            this.processBoard();
                        });
                    } else {
                        console.log(jsonData)
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
            ws.send(JSON.stringify(message));
            this.setState({isInitiator: true, active: false});
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
            ws.send(JSON.stringify(message));
            this.setState({active: false});
        }
    }

    sendMoveToServer = () => {
        const {ws} = this.state;
        if (ws || ws.readyState === WebSocket.OPEN) {
            const message = {
                command: "NEW_MOVE",
                gameInfo: {
                    gameId:  this.state.gameId,
                    state:   this.state.boardState.join(),
                    lastMoveUserId: this.state.userId,
                }
            }
            console.log(JSON.stringify(message))
            ws.send(JSON.stringify(message));
        }
    }

    processBoard = () => {
        let won = false;
        const {
            boardState,
        } = this.state;

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
        }
    }

    handleNewMove = (id) => {
        this.setState(
            prevState => {
                return {
                    boardState: prevState.boardState
                        .slice(0, id)
                        .concat(prevState.turn)
                        .concat(prevState.boardState.slice(id + 1)),
                };
            },
            () => {
                this.processBoard();
                this.sendMoveToServer();
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
            connAlertShow,
            connAlertText,
            connAlertType,
            lastMoveUserId,
            turn,
            userId,
        } = this.state;
        const rows = [];

        const isActive = active && (userId !== lastMoveUserId);

        for (let i = 0; i < 3; i++)
            rows.push(
                <Row
                    key={i}
                    row={i}
                    boardState={boardState}
                    onNewMove={this.handleNewMove}
                    active={isActive}
                />
            );
        return (
            <div>
                <Jumbotron
                    className="container"
                >
                    <Header/>
                    <p>Очередь за {String.fromCharCode(symbolsMap[turn][1])}</p>
                    <br/>
                    <div className="board">
                        {rows}
                    </div>
                    <div className="alert-container">
                        <Alert color={alertType} show={alertShow} isOpen={alertShow}>
                            {alertText}
                        </Alert>
                        <Alert color={connAlertType} show={connAlertShow} isOpen={connAlertShow}>
                            {connAlertText}
                        </Alert>
                    </div>
                </Jumbotron>
                {
                    this.state.shareButtonShow ? <Share gameCode={this.state.gameId}/> : null
                }
            </div>
        );
    }
}

export default withRouter(OnlineGamePage);
