
export default class TournamentService {
    static getTournament(tournamentId) {
      return MOCK_TOURNAMENT;
    }

    static getTournamentList() {
      console.log("Fetching Tournament List.");
      return MOCK_TOURNAMENT_LIST;
    }
    
    static getGameList(tournamentId) {
      return MOCK_GAME_LIST;
    }

    static getGame(tournamentId, gameId) {
      return MOCK_GAME;
    }
}

const MOCK_TOURNAMENT = {
  "id": "59a8c6787af8b4bdbeaeb888",
  "tags": null,
  "status": "Pending",
  "settings": {
    "testSeats": 2,
    "carousel": true,
    "rounds": 10,
    "engines": [
      {
        "Id": "",
        "name": "exacto",
        "version": "0.f",
        "protocol": "winboard",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.39/DirtyBit.0.39.5.t5.exe"
      },
      {
        "Id": "",
        "name": "dirty-bit",
        "version": "0.2",
        "protocol": "uci",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.20/DirtyBit.0.2.exe"
      }
    ],
    "timeControl": {
      "moves": 40,
      "repeating": true
    }
  },
  "summary": {}
};

const MOCK_TOURNAMENT_LIST = [
  MOCK_TOURNAMENT
];

const MOCK_GAME_LIST = [
  {
    "id": "59a8c6787af8b4bdbeaeb889",
    "tournamentId": "59a8c6787af8b4bdbeaeb888",
    "status": "Pending",
    "contestants": {
      "0": {
        "Id": "",
        "name": "exacto",
        "version": "0.f",
        "protocol": "winboard",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.39/DirtyBit.0.39.5.t5.exe"
      },
      "1": {
        "Id": "",
        "name": "stockfish",
        "version": "",
        "protocol": "uci",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.20/DirtyBit.0.2.exe"
      }
    },
    "timeControl": {
      "moves": 40,
      "repeating": true
    },
    "positions": []
  },
  {
    "id": "59a8c6787af8b4bdbeaeb88a",
    "tournamentId": "59a8c6787af8b4bdbeaeb888",
    "status": "Pending",
    "contestants": {
      "0": {
        "Id": "",
        "name": "dirty-bit",
        "version": "0.2",
        "protocol": "uci",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.20/DirtyBit.0.2.exe"
      },
      "1": {
        "Id": "",
        "name": "exacto",
        "version": "0.f",
        "protocol": "winboard",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.39/DirtyBit.0.39.5.t5.exe"
      }
    },
    "timeControl": {
      "moves": 40,
      "repeating": true
    },
    "positions": []
  },
  {
    "id": "59a8c6787af8b4bdbeaeb88b",
    "tournamentId": "59a8c6787af8b4bdbeaeb888",
    "status": "Pending",
    "contestants": {
      "0": {
        "Id": "",
        "name": "exacto",
        "version": "0.f",
        "protocol": "winboard",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.39/DirtyBit.0.39.5.t5.exe"
      },
      "1": {
        "Id": "",
        "name": "dirty-bit",
        "version": "0.2",
        "protocol": "uci",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.20/DirtyBit.0.2.exe"
      }
    },
    "timeControl": {
      "moves": 40,
      "repeating": true
    },
    "positions": []
  },
  {
    "id": "59a8c6787af8b4bdbeaeb88c",
    "tournamentId": "59a8c6787af8b4bdbeaeb888",
    "status": "Pending",
    "contestants": {
      "0": {
        "Id": "",
        "name": "dirty-bit",
        "version": "0.2",
        "protocol": "uci",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.20/DirtyBit.0.2.exe"
      },
      "1": {
        "Id": "",
        "name": "exacto",
        "version": "0.f",
        "protocol": "winboard",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.39/DirtyBit.0.39.5.t5.exe"
      }
    },
    "timeControl": {
      "moves": 40,
      "repeating": true
    },
    "positions": []
  },
  {
    "id": "59a8c6787af8b4bdbeaeb88d",
    "tournamentId": "59a8c6787af8b4bdbeaeb888",
    "status": "Pending",
    "contestants": {
      "0": {
        "Id": "",
        "name": "exacto",
        "version": "0.f",
        "protocol": "winboard",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.39/DirtyBit.0.39.5.t5.exe"
      },
      "1": {
        "Id": "",
        "name": "dirty-bit",
        "version": "0.2",
        "protocol": "uci",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.20/DirtyBit.0.2.exe"
      }
    },
    "timeControl": {
      "moves": 40,
      "repeating": true
    },
    "positions": []
  },
  {
    "id": "59a8c6787af8b4bdbeaeb88e",
    "tournamentId": "59a8c6787af8b4bdbeaeb888",
    "status": "Pending",
    "contestants": {
      "0": {
        "Id": "",
        "name": "dirty-bit",
        "version": "0.2",
        "protocol": "uci",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.20/DirtyBit.0.2.exe"
      },
      "1": {
        "Id": "",
        "name": "exacto",
        "version": "0.f",
        "protocol": "winboard",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.39/DirtyBit.0.39.5.t5.exe"
      }
    },
    "timeControl": {
      "moves": 40,
      "repeating": true
    },
    "positions": []
  },
  {
    "id": "59a8c6787af8b4bdbeaeb88f",
    "tournamentId": "59a8c6787af8b4bdbeaeb888",
    "status": "Pending",
    "contestants": {
      "0": {
        "Id": "",
        "name": "exacto",
        "version": "0.f",
        "protocol": "winboard",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.39/DirtyBit.0.39.5.t5.exe"
      },
      "1": {
        "Id": "",
        "name": "dirty-bit",
        "version": "0.2",
        "protocol": "uci",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.20/DirtyBit.0.2.exe"
      }
    },
    "timeControl": {
      "moves": 40,
      "repeating": true
    },
    "positions": []
  },
  {
    "id": "59a8c6787af8b4bdbeaeb890",
    "tournamentId": "59a8c6787af8b4bdbeaeb888",
    "status": "Pending",
    "contestants": {
      "0": {
        "Id": "",
        "name": "dirty-bit",
        "version": "0.2",
        "protocol": "uci",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.20/DirtyBit.0.2.exe"
      },
      "1": {
        "Id": "",
        "name": "exacto",
        "version": "0.f",
        "protocol": "winboard",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.39/DirtyBit.0.39.5.t5.exe"
      }
    },
    "timeControl": {
      "moves": 40,
      "repeating": true
    },
    "positions": []
  },
  {
    "id": "59a8c6787af8b4bdbeaeb891",
    "tournamentId": "59a8c6787af8b4bdbeaeb888",
    "status": "Pending",
    "contestants": {
      "0": {
        "Id": "",
        "name": "exacto",
        "version": "0.f",
        "protocol": "winboard",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.39/DirtyBit.0.39.5.t5.exe"
      },
      "1": {
        "Id": "",
        "name": "dirty-bit",
        "version": "0.2",
        "protocol": "uci",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.20/DirtyBit.0.2.exe"
      }
    },
    "timeControl": {
      "moves": 40,
      "repeating": true
    },
    "positions": []
  },
  {
    "id": "59a8c6787af8b4bdbeaeb892",
    "tournamentId": "59a8c6787af8b4bdbeaeb888",
    "status": "Pending",
    "contestants": {
      "0": {
        "Id": "",
        "name": "dirty-bit",
        "version": "0.2",
        "protocol": "uci",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.20/DirtyBit.0.2.exe"
      },
      "1": {
        "Id": "",
        "name": "exacto",
        "version": "0.f",
        "protocol": "winboard",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.39/DirtyBit.0.39.5.t5.exe"
      }
    },
    "timeControl": {
      "moves": 40,
      "repeating": true
    },
    "positions": []
  },
  {
    "id": "59a8c6787af8b4bdbeaeb893",
    "tournamentId": "59a8c6787af8b4bdbeaeb888",
    "status": "Pending",
    "contestants": {
      "0": {
        "Id": "",
        "name": "exacto",
        "version": "0.f",
        "protocol": "winboard",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.39/DirtyBit.0.39.5.t5.exe"
      },
      "1": {
        "Id": "",
        "name": "dirty-bit",
        "version": "0.2",
        "protocol": "uci",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.20/DirtyBit.0.2.exe"
      }
    },
    "timeControl": {
      "moves": 40,
      "repeating": true
    },
    "positions": []
  },
  {
    "id": "59a8c6787af8b4bdbeaeb894",
    "tournamentId": "59a8c6787af8b4bdbeaeb888",
    "status": "Pending",
    "contestants": {
      "0": {
        "Id": "",
        "name": "dirty-bit",
        "version": "0.2",
        "protocol": "uci",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.20/DirtyBit.0.2.exe"
      },
      "1": {
        "Id": "",
        "name": "exacto",
        "version": "0.f",
        "protocol": "winboard",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.39/DirtyBit.0.39.5.t5.exe"
      }
    },
    "timeControl": {
      "moves": 40,
      "repeating": true
    },
    "positions": []
  },
  {
    "id": "59a8c6787af8b4bdbeaeb895",
    "tournamentId": "59a8c6787af8b4bdbeaeb888",
    "status": "Pending",
    "contestants": {
      "0": {
        "Id": "",
        "name": "exacto",
        "version": "0.f",
        "protocol": "winboard",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.39/DirtyBit.0.39.5.t5.exe"
      },
      "1": {
        "Id": "",
        "name": "dirty-bit",
        "version": "0.2",
        "protocol": "uci",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.20/DirtyBit.0.2.exe"
      }
    },
    "timeControl": {
      "moves": 40,
      "repeating": true
    },
    "positions": []
  },
  {
    "id": "59a8c6787af8b4bdbeaeb896",
    "tournamentId": "59a8c6787af8b4bdbeaeb888",
    "status": "Pending",
    "contestants": {
      "0": {
        "Id": "",
        "name": "dirty-bit",
        "version": "0.2",
        "protocol": "uci",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.20/DirtyBit.0.2.exe"
      },
      "1": {
        "Id": "",
        "name": "exacto",
        "version": "0.f",
        "protocol": "winboard",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.39/DirtyBit.0.39.5.t5.exe"
      }
    },
    "timeControl": {
      "moves": 40,
      "repeating": true
    },
    "positions": []
  },
  {
    "id": "59a8c6787af8b4bdbeaeb897",
    "tournamentId": "59a8c6787af8b4bdbeaeb888",
    "status": "Pending",
    "contestants": {
      "0": {
        "Id": "",
        "name": "exacto",
        "version": "0.f",
        "protocol": "winboard",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.39/DirtyBit.0.39.5.t5.exe"
      },
      "1": {
        "Id": "",
        "name": "dirty-bit",
        "version": "0.2",
        "protocol": "uci",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.20/DirtyBit.0.2.exe"
      }
    },
    "timeControl": {
      "moves": 40,
      "repeating": true
    },
    "positions": []
  },
  {
    "id": "59a8c6787af8b4bdbeaeb898",
    "tournamentId": "59a8c6787af8b4bdbeaeb888",
    "status": "Pending",
    "contestants": {
      "0": {
        "Id": "",
        "name": "dirty-bit",
        "version": "0.2",
        "protocol": "uci",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.20/DirtyBit.0.2.exe"
      },
      "1": {
        "Id": "",
        "name": "exacto",
        "version": "0.f",
        "protocol": "winboard",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.39/DirtyBit.0.39.5.t5.exe"
      }
    },
    "timeControl": {
      "moves": 40,
      "repeating": true
    },
    "positions": []
  },
  {
    "id": "59a8c6787af8b4bdbeaeb899",
    "tournamentId": "59a8c6787af8b4bdbeaeb888",
    "status": "Pending",
    "contestants": {
      "0": {
        "Id": "",
        "name": "exacto",
        "version": "0.f",
        "protocol": "winboard",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.39/DirtyBit.0.39.5.t5.exe"
      },
      "1": {
        "Id": "",
        "name": "dirty-bit",
        "version": "0.2",
        "protocol": "uci",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.20/DirtyBit.0.2.exe"
      }
    },
    "timeControl": {
      "moves": 40,
      "repeating": true
    },
    "positions": []
  },
  {
    "id": "59a8c6787af8b4bdbeaeb89a",
    "tournamentId": "59a8c6787af8b4bdbeaeb888",
    "status": "Pending",
    "contestants": {
      "0": {
        "Id": "",
        "name": "dirty-bit",
        "version": "0.2",
        "protocol": "uci",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.20/DirtyBit.0.2.exe"
      },
      "1": {
        "Id": "",
        "name": "exacto",
        "version": "0.f",
        "protocol": "winboard",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.39/DirtyBit.0.39.5.t5.exe"
      }
    },
    "timeControl": {
      "moves": 40,
      "repeating": true
    },
    "positions": []
  },
  {
    "id": "59a8c6787af8b4bdbeaeb89b",
    "tournamentId": "59a8c6787af8b4bdbeaeb888",
    "status": "Pending",
    "contestants": {
      "0": {
        "Id": "",
        "name": "exacto",
        "version": "0.f",
        "protocol": "winboard",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.39/DirtyBit.0.39.5.t5.exe"
      },
      "1": {
        "Id": "",
        "name": "dirty-bit",
        "version": "0.2",
        "protocol": "uci",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.20/DirtyBit.0.2.exe"
      }
    },
    "timeControl": {
      "moves": 40,
      "repeating": true
    },
    "positions": []
  },
  {
    "id": "59a8c6787af8b4bdbeaeb89c",
    "tournamentId": "59a8c6787af8b4bdbeaeb888",
    "status": "Pending",
    "contestants": {
      "0": {
        "Id": "",
        "name": "dirty-bit",
        "version": "0.2",
        "protocol": "uci",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.20/DirtyBit.0.2.exe"
      },
      "1": {
        "Id": "",
        "name": "exacto",
        "version": "0.f",
        "protocol": "winboard",
        "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.39/DirtyBit.0.39.5.t5.exe"
      }
    },
    "timeControl": {
      "moves": 40,
      "repeating": true
    },
    "positions": []
  }
];

const MOCK_GAME = {
  "id": "59a8c6787af8b4bdbeaeb89c",
  "tournamentId": "59a8c6787af8b4bdbeaeb888",
  "status": "Pending",
  "contestants": {
    "0": {
      "Id": "",
      "name": "dirty-bit",
      "version": "0.2",
      "protocol": "uci",
      "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.20/DirtyBit.0.2.exe"
    },
    "1": {
      "Id": "",
      "name": "dirty-bit",
      "version": "0.39.5.t5",
      "protocol": "uci",
      "url": "https://github.com/andrewbackes/dirty-bit/releases/download/v0.39/DirtyBit.0.39.5.t5.exe"
    }
  },
  "timeControl": {
    "moves": 40,
    "repeating": true
  },
  "positions": []
};