import React from 'react';
import GamePage from "./GamePage";

class GameJoinPage extends React.Component {
    constructor(props) {
        super(props);

        let gameCode;

        if (props.params && props.params.gameCode) {
            gameCode = props.params.gameCode;
        }

        this.state = {
            username: null,
            gameCode: gameCode,
            isFormSubmittable: false
        }
    }

    render() {
        return this.joinGame()
    }

    joinGame() {
        if (this.state.gameCode && this.state.username) {
            return(<GamePage gameCode={this.state.gameCode} username={this.state.username} />)
        }

        return (
            <div className="container">
                <GamePage handleSubmit={this.handleSubmit.bind(this)} submitTrans="Join" isSubmitable={this.state.isFormSubmittable}>
                    {this.getInputFields()}
                </GamePage>
            </div>
        )
    }

    getInputFields () {
        if (this.state.gameCode) {
            return(
                <fieldset className="form-group"><input type="text" className="form-control" id="username" placeholder="Enter Username" ref={c => this._username = c} onChange={this.handleFormChange.bind(this)} maxlength="20" /></fieldset>
            );
        }

        return (
            <div>
                <fieldset className="form-group"><input type="text" className="form-control" id="username" placeholder="Enter Username" ref={c => this._username = c} onChange={this.handleFormChange.bind(this)} maxlength="20" /></fieldset>
                <fieldset className="form-group"><input type="text" className="form-control" id="gameCode" placeholder="Enter GameCode" ref={c => this._gameCode = c} onChange={this.handleFormChange.bind(this)} /></fieldset>
            </div>
        );
    }

    handleSubmit(event) {
        event.preventDefault();

        if (this._username.value.length > 20) {
            alert("Выберите имя покороче");
        }

        if (this._username.value.length < 3) {
            alert("Выберите имя подлиннее");
        }

        this.setState({username: this._username.value});
        this._username.value = '';

        if (this._gameCode && this._gameCode.value) {
            this.setState({gameCode: this._gameCode.value});
            this._gameCode.value = '';
        }
    }

    handleFormChange() {
        if (! this._username.value) {
            return;
        }

        if (this._gameCode && ! this._gameCode.value) {
            return;
        }

        this.setState({isFormSubmittable: true});
    }
}

export default GameJoinPage;
