import React, { Component } from 'react';
import Panel from 'components/panel';
import TournamentService from 'services/tournament';

import MoveTable from 'scenes/game/move-list-table';
import EngineAnalysisTable from 'scenes/game/engine-analysis-table';
import Board from 'scenes/game/chessboard';

export default class GameDashboard extends Component {
  constructor(props) {
    super(props);
    this.state = {
      game: {},
      position: {
        fen: ""
      }
    };
    this.initGame = this.initGame.bind(this);
    this.setGame = this.setGame.bind(this);
    this.setPosition = this.setPosition.bind(this);
    TournamentService.getGame(this.props.match.params.tournamentId, this.props.match.params.gameId, this.initGame);
  }

  componentDidMount() {
    this.timerID = setInterval(
      () => this.refreshGame(),
      1000
    );
  }

  componentWillUnmount() {
    clearInterval(this.timerID);
  }

  initGame(game) {
    this.setState({ game: game });
    if (game.status === "Complete") {
      clearInterval(this.timerID);
      this.setPosition(game.positions[0]);
    } else {
      this.setPosition(this.state.game.positions[this.state.game.positions.length-1]);
    }
  }

  setGame(game) {
    let updatePosition = false;
    if (this.state.game && this.state.game.positions) {
      updatePosition = this.state.position === this.state.game.positions[this.state.game.positions.length-1]
    }
    if (this.timerID) {
      this.setState({ game: game });
    }
    if (updatePosition) {
      this.setPosition(this.state.game.positions[this.state.game.positions.length-1]);
    }
    
    if (this.state.game.status === "Complete") {
      clearInterval(this.timerID);
    }
  }

  setPosition(position) {
    if (this.timerID) {
      this.setState({ position: position });
    }
  }

  refreshGame() {
    TournamentService.getGame(this.props.match.params.tournamentId, this.props.match.params.gameId, this.setGame)
  }
  
  render() {
    let mode = 'default';
    if (this.state.game.status) {
      if (this.state.game.status.toLowerCase() === "running") {
        mode = 'success';
      } else if (this.state.game.status.toLowerCase() === "pending") {
        mode = 'info';
      }
    }
    return (
      <div>
        <div className="row">
          <div className="col-xs-4">
            <Panel title="Board" mode={ mode } content={<Board position={this.state.position}/>}/>
          </div>
          <div className="col-xs-3">
            <Panel title="Moves" mode="default" content={ <MoveTable game={this.state.game} setPosition={this.setPosition} currentPosition={this.state.position} /> }/>
          </div>
          <div className="col-xs-5">
            <Panel title="Analysis" mode="default"/>
          </div>
        </div>
        <div className="row">
          <div className="col-xs-12">
            <Panel title="Engine Output" mode='default' content={<EngineAnalysisTable analysis={this.state.position.analysis}/>}/>
          </div>
        </div>
      </div>
    );
  }
}
