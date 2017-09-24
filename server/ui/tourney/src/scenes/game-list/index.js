import React, { Component } from 'react';
import Panel from 'components/panel';
import TournamentService from 'services/tournament';

import Search from 'scenes/game-list/search';
import GameTable from 'scenes/game-list/game-table';

export default class GameList extends Component {
  constructor(props) {
    super(props);
    this.state = {
      tournament: {},
      gameList: [],
      filterText: ""
    };
    this.handleFilterTextInput = this.handleFilterTextInput.bind(this);
    this.setTournament = this.setTournament.bind(this);
    this.setGameList = this.setGameList.bind(this);
    this.refreshState()
  }

  componentDidMount() {
    if (this.state.tournament.status !== "Complete") {
      this.timerID = setInterval(
        () => this.refreshState(),
        1000
      );
    }
  }

  componentWillUnmount() {
    clearInterval(this.timerID);
  }

  handleFilterTextInput(filterText) {
    this.setState({
      filterText: filterText
    });
  }
  
  setTournament(tournament) {
    if (this.timerID) {
      this.setState({ tournament: tournament });
      if (this.state.tournament.status === "Complete") {
        clearInterval(this.timerID);
      }
    }
  }

  setGameList(gameList) {
    if (this.timerID) {
      this.setState({ gameList: gameList });
    }
  }

  refreshState() {
    TournamentService.getTournament(this.props.match.params.tournamentId, this.setTournament)
    TournamentService.getGameList(this.props.match.params.tournamentId, this.setGameList)
  }
  
  render() {
    return (
      <div>
        <Panel title="Search" mode="default" content={
          <Search onFilterTextInput={this.handleFilterTextInput} filterText={this.state.filterText}/>
        }/>
        <div className="panel panel-default">
          <div className="panel-body">
            <GameTable gameList={this.state.gameList} filterText={this.state.filterText} history={this.props.history}/>
          </div>
        </div>
      </div>
    );
  }
}
