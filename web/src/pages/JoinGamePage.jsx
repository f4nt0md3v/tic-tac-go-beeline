import React from "react";
import {Jumbotron} from "reactstrap";
import {Link} from "react-router-dom";
import {Redirect} from "react-router";

class JoinGamePage extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            gameId: '',
            canProceed: false,
        }
    }

    joinGamePage() {
        if (this.state.gameId !== '' && this.state.canProceed) {
            return (<Redirect to={{pathname: `/game/join/${this.state.gameId}`}} />)
        }

        return (
            <div>
                <Jumbotron
                    className={"container"}
                >
                    <h3>Игра "Крестики-Нолики"</h3>
                    <hr/>
                    <div>
                        <div>
                            <fieldset className="form-group">
                                <input
                                    type="text"
                                    className="form-control"
                                    placeholder="Введите код игры"
                                    ref={c => this.gameId = c}
                                    onChange={this.handleFormChange.bind(this)}/>
                            </fieldset>
                            <Link
                                className="btn btn-block btn-primary" to="#"
                                onClick={this.handleFormSubmit.bind(this)} size="sm"
                            >
                                Подключиться
                            </Link>
                        </div>
                    </div>
                </Jumbotron>
            </div>
        )
    }

    handleFormChange() {
        if (this.gameId && !this.gameId.value) {
            return;
        }
        this.setState({gameId: this.gameId.value});
    }

    handleFormSubmit() {
        if (this.gameId && !this.gameId.value) {
            return;
        }
        this.setState({canProceed: true});
    }

    render() {
        return this.joinGamePage()
    }
}

export default JoinGamePage;
