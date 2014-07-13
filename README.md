Project Tourney
===============


What is it?
-----------

A simple cross platform console application to runs chess tournaments and view results remotely.


What can it do?
---------------

Launched from terminal, it will prompt you for what kind of tourney you want to set up and with what engines.Optionally, this information can be kept in a separate file and included in the command line argument. Once the tournament begins, minimal data will start to be printed on the screen. This may include: current tournament standings and current game data (possibly: time, scores, score chart, moves, etc.). Commands can be inputted to get additional data and modify tournament settings. It will be easy to get all of the data regarding a certain move in a particular game (getting this data may not need to be in the scope of this program, but the way the data is saved needs to allow for easy access). All output from the engines will be saved to logs or databases. Results and periodic updates will be saved to a web server for easy remote viewing.


What is the goal for the first version?
---------------------------------------

 * Run tournaments for adjusted 40/5 and 40/40 time controls
 	* Detect checkmates, stalemates, repetitions, illegal moves, time, etc.
 	* Use standard format opening books
 * Round Robin, Gauntlet, and Multi-Gauntlet
 * Print tourney data to screen in real time
 * Accept commands to control the tourney in real time
 * Accept commands to print more advanced data to screen. i.e. game specific details
 * Save logs of complete engine output
 * upload basic score results to a web server
 * Use a *tournament file* to run a tourney with specific settings. Without this file, prompt for info on what kind of tourney to run.
 * Written in a *protocol fashion* so that a windowed or web based GUI can easily be added as a front end


Questionable
------------

Should there be a mode to open a previously played tournament and go through the results? If so, this needs to be easy to use yet powerful. It may be better to do this portion with a web app.


How can it be extended?
-----------------------
 
 * More time controls
 * Detect blunders
 * Can be controlled remotely. First maybe by connecting via terminal, then later maybe by a web interface
 * Tournament queue
 * After a tourney finishes, runs a script or executable before running the next tourny in queue. The idea here is to run an engine specific tuning program between tourneys.
 * Tourney can be distributed to run on multiple machines


What else?
----------

It should be very easy and straightforward to use. Users should not have to read directions on how to use it. The program itself should guide users.

Observation: When running tournaments in Arena there is usually only very specific data I am interested in depending on the stage of engine design. For the most part, I want to see the current scores in the tournament. It is only if things are not going as expected that I want to see more information. Next, I would want to see what kind of losses there were: illegal move, absolute slaughter, close game, etc, was there a blunder? At that point, game specific data is usually looked at. That would be move list, pv, and score. Maybe engine log.


More About Arena
----------------

Useful features: 
 * Score chart
 * debug log console
 * tournament results page

Difficulties: 
 * Once a move is identified on the score chart or move list, viewing _all_ of the engine's output regarding that move in the debug console (pv, best-move, scores, errors, 'info string' output) takes a lot of effort.

Incapabilities: 
 * Can not check the results of a tournament remotely. Much less start, reset, or modify a tournament remotely.
 * Can not set a gauntlet for more than one engine to play against all the rest. I call this multi-gauntlet.

