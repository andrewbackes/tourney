# Specifications

## Worker API Calls

Worker should be able to fetch a game that needs to be played based on tags:
```
GET /games?status=pending&contestant=engine1&hash=1234567&limit=1
[
    {
        "id": "...",
        "tournamentId": "...",
        "status": "pending",
        "tags": {
            "result": "*",
            "black": "engine1",
            "white": "engine2",
            "fEN": "...",
            "whiteControl": "40/40",
            "blackControl": "40/40",
            ...
        }
        "positions": [
            {
                "board": "...",
                "moveNumber": 1,
                "activeColor: "white",
                ...
            }
        ]
    }
]
```

Worker should be able to update positions for a game:
```
PUT /games/{id}/positions/{moveNumber}/{color}
{
    "board": "...",
    "moveNumber": 1,
    "activeColor": "black",
    ...
} 
```

Worker should be able to update tags for a game:
```
PATCH /games/{id}/tags
{
    "result": "1-0"
}
```

Worker should be able to update status:
```
PUT /games/{id}/status
"complete"
```

## UI API Calls

User should be able to add a new tournament:
```
POST /tournaments
{
    "contestants": [
        {"id": "..."},
        {"name": "dirty-bit", "version": "0.39", "url": "..."}
    ],
    "book": {"id":"..."}
}
```

UI should be able to get a list of tournaments filtered by status/tag and sorted by post date.
```
GET /tournaments?status=pending
[
    {
        ...
    }
]
```

## UI

### Tournaments

User should be able view:
- Runing, pending, and completed tournaments
- filter by tag

### Games

After selecting a tournament, user should be able to view:
- Completed, running, and pending games.
- Detailed game analysis.
- Watch a game in progress.


## API

Game Status put'ed -> Summarize Tournament

Tournament should contain:
    Settings
    Summary
    Games
