## GoExpect
### What am I?
GoExpect runs through scripts similar to the classic expect command. This version uses lua scripts instead of tcl and works on windows, mac and (probably) the other supported golang os.
### Build
* Ensure that the latest version of Golang is installed (1.8 at the time of writing this)
* Build as usual.
### Scripting
* Uses Lua 5.1 using the github.com/yuin/gopher-lua wrapper.
* Extra functions:
``` 
--starts the subprocess with the given name and arguments
spawn("path to command name", "command arguments") 
```
```
--sets the number of seconds to wait before timing out when expect()-ing
timeout(15)
```
``` 
--sends the given string to the subprocess
send("string to send")
```
```
--returns false if the timeout is reached and the given string is not found in the output
--returns true if the given string is in the output before the timeout period
expect("string to expect")
```
```
--sleeps the script for a specified number of milliseconds
sleep(1000)
```
```
--closes the currently running spawned sub process
exit()
```
```
--argv holds the command line arguments as specified when calling GoExpect
ip = argv[0]
```
