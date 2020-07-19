import React from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faSpinner } from "@fortawesome/free-solid-svg-icons";
import copy from "copy-to-clipboard";
import { Link } from "react-router-dom";

class Share extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            copied: false,
        };
    }

    getUrl() {
        const origin = window.location.origin;
        return origin + '/#/game/join/' + this.props.gameCode;
    }

    getShareButton() {
        let btnText = 'Копировать ссылку';
        if (this.state.copied) {
            btnText = `\u2713 Скопировано в буфер`;
        }
        setTimeout(() => this.setState({ copied: false }), 5000);
        return (<Link id="btn-reset-game" className="btn btn-block btn-primary" to="#" onClick={this.copyShareUrl.bind(this)} size="sm">{btnText}</Link>)
    }

    copyShareUrl() {
        copy(this.getUrl());
        this.setState({ copied: true });
    }

    render() {
        return (
            <div>
                <div className="container">
                    <FontAwesomeIcon size="lg" icon={faSpinner} spin={true} />
                    <br />
                    <p className="loading-text">Ждём оппонента...</p>
                    <kbd className="sharing-screen-link">{this.getUrl()}</kbd>
                    <br />
                    <br />
                    <p className="sharing-screen-info">Поделитесь этой ссылкой и начните игру</p>
                    {this.getShareButton()}
                </div>
            </div>
        );
    }
}

export default Share;
