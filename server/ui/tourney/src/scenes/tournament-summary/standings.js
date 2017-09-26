import React, { Component } from 'react';

export default class StandingsTable extends Component {
  render() {
    let contestantsById = {};
    this.props.contestants.forEach( (contestant) => {
      contestantsById[contestant.id] = contestant;
    });
    let records = [];
    for (let contestantId in this.props.stats) {
      records.push({id: contestantId, contestant: contestantsById[contestantId], entry: this.props.stats[contestantId]});
    }
    records.sort( (a, b) => {
      if (score(a.entry) > score(b.entry)) { return -1; }
      if (score(a.entry) < score(b.entry)) { return 1;}
      return 0;
    });
    let rows = [];
    records.forEach( (record) => {
      rows.push(<Row key={record.id} id={record.id} contestant={record.contestant} entry={record.entry}/>)
    })
    return (
      <table className="table table-condensed">
        <thead>
          <tr>
            <th>Position</th>
            <th>Name</th>
            <th>Record</th>
            <th>Incomplete</th>
          </tr>
        </thead>
        <tbody>
          { rows }
        </tbody>
      </table>
    );
  }
}

function score(record) {
  return record.wins + record.draws/2;
}

function formatRecord(record) {
  return record.wins + "-" + record.losses + "-" + record.draws;
}

function formatContestantLabel(contestant) {
  let l = contestant.name + " " + contestant.version;
  if (l !== "") {
    return l;
  }
  return contestant.id;
}

class Row extends Component {
  render() {
    return (
      <tr>
        <td>-</td>
        <td>{formatContestantLabel(this.props.contestant)}</td> 
        <td>{formatRecord(this.props.entry)}</td>
        <td>{this.props.entry.incomplete}</td> 
      </tr>
    );
  }
}