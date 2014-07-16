## Stdio-fork


#### Launching

Try `go run dan.go`.

#### Spawning children processes

From within Stdio-fork try `s [command]`.  E.g., `s python my_script.py`  Stdiofork will generate a pipe number (in order from 0) called the *address* of the pipe.

#### Sending commands to children processes

To send a command over stdin to the process piped to address `n`, try `pn [command]`.  E.g., `p0 run`.