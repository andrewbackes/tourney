import React, { Component } from 'react';
import Panel from 'components/panel';
import TournamentService from 'services/tournament';

export default class GameList extends Component {
  constructor(props) {
    super(props);
    this.state = {
      tournament: TournamentService.getTournament(this.props.match.params.tournamentId),
      gameList: TournamentService.getGameList(this.props.match.params.tournamentId),
      filterText: ""
    };
    this.handleFilterTextInput = this.handleFilterTextInput.bind(this);
  }

  componentDidMount() {
    if (this.state.tournament.status !== "Complete") {
      this.timerID = setInterval(
        () => this.refreshGameList(),
        10000
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
  

  refreshGameList() {
    this.setState({
      tournament: TournamentService.getTournament("something"),
      gameList: TournamentService.getGameList("something")
    });
    if (this.state.tournament.status === "Complete") {
      clearInterval(this.timerID);
    }
  }
  
  render() {
    return (
      <div>
        <Panel title="Search" mode="default" content={
          <Search onFilterTextInput={this.handleFilterTextInput} filterText={this.state.filterText}/>
        }/>
        <div className="panel panel-default">
          <div className="panel-body">
            <GameTable gameList={this.state.gameList} filterText={this.state.filterText}/>
          </div>
        </div>
      </div>
    );
  }
}

class GameTable extends Component {
  
  filterGame(game, filterText) {
    if (filterText === "") {
      return true;
    }
    return JSON.stringify(game).includes(filterText);
  }

  render() {
    let rows = [];
    this.props.gameList.forEach( (game) => {
      if (this.filterGame(game, this.props.filterText)) {
        rows.push(<GameTableRow key={game.id} game={game}/>)
      }
    })
    return (
      <table className="table table-hover table-condensed">
        <thead>
          <tr>
            <th>Round</th>
            <th>White</th>
            <th>Black</th>
            <th>Result</th>
            <th>Winner</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          { rows }
        </tbody>
      </table>
    );
  }
}

function engineLabel(engine) {
  return engine.name + " " + engine.version;
}

class GameTableRow extends Component {
  render() {
    return (
      <tr>
        <td>{this.props.game.round}</td>
        <td>{engineLabel(this.props.game.contestants["0"])}</td> 
        <td>{engineLabel(this.props.game.contestants["1"])}</td> 
        <td>-</td>
        <td>-</td>
        <td>{this.props.game.status}</td>
      </tr>
    );
  }
}
class Search extends Component {
  constructor(props) {
    super(props);
    this.handleFilterTextInputChange = this.handleFilterTextInputChange.bind(this);
  }

  handleFilterTextInputChange(e) {
    this.props.onFilterTextInput(e.target.value);
  }

  render() {
    return (
      <form>
        <input
          type="text"
          placeholder="Search..."
          value={this.props.filterText}
          onChange={this.handleFilterTextInputChange}
        />
      </form>
    );
  }
}