import $ from 'jquery';

export default class TournamentService {
  
  static apiHost = 'http://localhost:9090/api/v2';

  static getTournament(tournamentId, callback) {
    $.get(this.apiHost + '/tournament/' + tournamentId, callback);
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
  
  static getGameList(tournamentId, callback) {
    $.get(this.apiHost + '/tournaments/' + tournamentId + '/games',callback);
  }

  static getGame(tournamentId, gameId, callback) {
    $.get(this.apiHost + '/tournaments/' + tournamentId + '/games/' + gameId,callback);
  }
}
