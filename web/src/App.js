import React from 'react';
import "bootstrap/dist/css/bootstrap.css";
import {Route} from "react-router";
import {HashRouter} from "react-router-dom";
import IndexPage from "./pages/IndexPage";
import GamePage from "./pages/GamePage";
import JoinGamePage from "./pages/JoinGamePage";
import OnlineGamePage from "./pages/OnlineGamePage";

class App extends React.Component {
    render() {
        return (
            <HashRouter>
                <Route path="/" component={IndexPage} exact />
                <Route path="/game/start" component={OnlineGamePage} />
                <Route path="/game/join" component={JoinGamePage} exact />
                <Route path="/game/join/:gameId" component={OnlineGamePage} />
                <Route path="/game/ai" component={GamePage} />
            </HashRouter>
        );
    }
}

export default App;
