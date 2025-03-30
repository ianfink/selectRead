selectRead executes a command-line program and reads stdout and stderr from the command.
While the Go (Golang) packages unix and unixutils are provided, implementing something
like a C “select loop” in Go can be challenging.  selectRead is an example of a Go program
that can provide the same functionality of a C “select loop” in Go that can be (easily)
tailored to the needs of a programmer.  Two go functions read stdout and stderr from the
command and provide the data from stdout and stderr via channels.  A timeout is also
provided, which shows an example of “doing something” if stdout or stderr exceed a time
limit.<br/>
testStderr is used to test selectRead.  testStderr provides output to stdout and stderr
at random times.  testStderr can be found among this author’s other software on github.com.

