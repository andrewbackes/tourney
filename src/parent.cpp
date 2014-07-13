#include <unistd.h>
#include <string>
#include <iostream>
#include <cstdlib>

using namespace std;

void error_out(int code, string message="Unhandled exception") {
  cerr << message << endl;
  exit(code);
}

string debuffer_pipe(int pip[2]) {
  string receive_output = "";
  char readbuffer[80];
  for (int bytes_read = (int)true; bytes_read;) {
    bytes_read = read(pip[0], readbuffer, sizeof(readbuffer)-1);
    readbuffer[bytes_read] = '\0';
    receive_output += readbuffer;
  }
  return receive_output;
}

int main() {
  int par_to_child_pipe[2], child_to_par_pipe[2];
  pid_t pid;
  string program_name = "./child";
  string command = "asdfx";


  if (pipe(par_to_child_pipe) || pipe(child_to_par_pipe))
    error_out(1, "Failed to pipe");
  pid = fork();

  if (pid < 0)
    error_out(-1, "Fork failed.");
  else if (pid == 0) {
    if (dup2(par_to_child_pipe[0], 0) != 0) {
      close(par_to_child_pipe[0]);
      close(par_to_child_pipe[1]);
      error_out(1, "Child failed to redirect stdin");
    }
    if (dup2(child_to_par_pipe[1], 1) != 1) {
      close(child_to_par_pipe[1]);
      close(child_to_par_pipe[0]);
      error_out(1, "Child failed to redirect stdout");
    }
    execl(program_name.c_str(), program_name.c_str(), (char *) 0);
    error_out(1, "Failed to execute child process");
  } else {
    close(par_to_child_pipe[0]);
    close(child_to_par_pipe[1]);

    int nbytes = command.length();
    if (write(par_to_child_pipe[1], command.c_str(), nbytes) != nbytes)
      error_out(1, "Failed to write to child");
    string receive_output = debuffer_pipe(child_to_par_pipe);
    close(par_to_child_pipe[1]);
    close(child_to_par_pipe[0]);
    cout << receive_output << endl;
  }
  return 0;
}
