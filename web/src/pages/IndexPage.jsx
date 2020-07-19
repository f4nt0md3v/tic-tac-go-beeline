import React from 'react';
import {Link} from "react-router-dom";

class IndexPage extends React.Component {
    render() {
        return (
            <div className="landing-page">
                <section className="jumbotron text-xs-center">
                    <div className="container">
                        <h1 className="jumbotron-heading">Игра "Крестики-Нолики"</h1>
                        <p className="lead text-muted">Добро пожаловать в игру "Крестики-Нолики"</p>
                        <p>
                            <Link to="/game/ai" className="btn btn-block btn-primary">Играть против компьютера</Link>
                            <Link to="/game/start" className="btn btn-block btn-primary">Играть против соперника</Link>
                            <Link to="/game/join" className="btn btn-block btn-secondary">Присоединиться к игре с соперником</Link>
                        </p>
                    </div>
                </section>
            </div>
        );
    }
}

export default IndexPage;
