import $ from 'jquery';

export default class TournamentService {
  
  static apiHost = 'http://localhost:9090/api/v2';

  static getTournament(tournamentId, callback) {
    $.ajax({
      url: this.apiHost + '/tournaments/' + tournamentId,
      type: "GET",
      dataType: 'json',
      contentType: 'application/json',
      success: callback,
      error: function (jqXHR, status, err) {
        console.log("ajax error getting tournament.");
      }
    });
  }

  static getTournamentList(status, callback) {
    $.ajax({
      url: this.apiHost + '/tournaments?status=' + status,
      type: "GET",
      dataType: 'json',
      contentType: 'application/json',
      success: callback,
      error: function (jqXHR, status, err) {
        console.log("ajax error getting tournament list.");
      }
    });
  }
  
  static getGameList(tournamentId, callback, status) {
    let suffix = "";
    if (status) {
      suffix = "?status=" + status
    }
    $.ajax({
      url: this.apiHost + '/tournaments/' + tournamentId + '/games' + suffix,
      type: "GET",
      dataType: 'json',
      contentType: 'application/json',
      success: callback,
      error: function (jqXHR, status, err) {
        console.log("ajax error getting game list.");
      }
    });
  }

  static getGame(tournamentId, gameId, callback) {
    $.ajax({
      url: this.apiHost + '/tournaments/' + tournamentId + '/games/' + gameId,callback,
      type: "GET",
      dataType: 'json',
      contentType: 'application/json',
      success: callback,
      error: function (jqXHR, status, err) {
        console.log("ajax error getting game.");
      }
    });
  }
}
