import React from 'react';
import "bootstrap/dist/css/bootstrap.css";
import {Route} from "react-router";
import {HashRouter} from "react-router-dom";
import IndexPage from "./pages/IndexPage";
import GamePage from "./pages/GamePage";
import GameJoinPage from "./pages/GameJoinPage";


class App extends React.Component {
    render() {
        return (
            <HashRouter>
                <Route exact path="/" component={IndexPage}/>
                <Route path="/game/start" component={GameJoinPage} />
                <Route path="/game/join" component={GamePage} />
                <Route path="/game/ai" component={GamePage} />
            </HashRouter>
        );
    }
}

export default App
