# fsisolatecli

Test app for [fsisolate](https://github.com/odacremolbap/fsisolate)

# Usage

Command pattern:

`isocli [OPTIONS] IMAGE COMMAND [command args ...]`

`IMAGE` can be one of the following:
- A non compressed tarball in a local directory
- A URL to a non compressed tarball
- A local directory

The chrooted process will contain this image.

`COMMAND` is the command to execute inside the chroot isolation. Hence, it must be contained inside the image.

**Options**

- `-a`, time to wait seconds after finishing process
- `-b`, time to wait in seconds before executing process
- `-d`, enable debug messages
- `-r`, directory to place the new root

By default there is a 2 seconds delay before executing the chrooted process.

If `IMAGE` is a directory, `-r` parameter will be ignored, and the chroot process will run in the directory

If `IMAGE`is URL or file, and `-r` is not present, a temp directory will be created to host the chrooted isolation.

# Running

Once `fsisolatecli` is running you can get chrooted process infor by writing:

- <kbd>s</kbd>+<kbd>Enter</kbd> to get process status
- <kbd>h</kbd>+<kbd>Enter</kbd> to send SIGHUP
- <kbd>i</kbd>+<kbd>Enter</kbd> to send SIGINT
- <kbd>k</kbd>+<kbd>Enter</kbd> to send SIGKILL
- <kbd>u</kbd>+<kbd>Enter</kbd> to send SIGUSR1

(Sorry for the need to press <kbd>Enter</kbd>, I'm due to try [termbox](https://github.com/nsf/termbox-go))

# Tests

`test` folder contains a `simple`subfolder containing a simple `loop-darwin` and `loop-linux` application that iterates a number of times an listens to OS signals.

You can also in that directory find the tarball file `simple.tar`.

The [loop](https://github.com/odacremolbap/loop) application executes a 1 second sleep iteration a number of times determined by flag `-i`, and in the end exits with status determined by `-e`

## Linux

This command will execute the loop app for 10 iterations, and will exit with code 0
It will also wait 5 seconds before executing the process, and 5 after the process has finished
```
sudo ./fsisolatecli -b 5 -a 5 test/simple /loop-linux -- -i 10 -e 0
```

To test a failing process you can change the exit code to a non 0 value
```
sudo ./fsisolatecli -b 5 -a 5 test/simple /loop-linux -- -i 10 -e 2
```

You can ask for status or send signals anytime. If you try to communicate with the process before it's running of after it has finished (use `-a`and `-b` delays), you should receive a message indicating that the process isn't running.

You can use a tarball file instead of a directory
```
sudo ./fsisolatecli -b 5 -a 5 test/simple.tar /loop-linux -- -i 10 -e 2
```

Or a URL
```
sudo ./fsisolatecli -b 5 -a 5 https://raw.githubusercontent.com/odacremolbap/fsisolatecli/master/test/simple.tar /loop-linux -- -i 10 -e 2
```

## Darwin

Darwin tests can also use `test\simple`, choosing the darwin compiled `loop` application.
```
sudo ./fsisolatecli -b 5 -a 5 test/simple /loop-darwin -- -i 10 -e 0
```
Directory with exit status 2, and no delay after executing the process
```
sudo ./fsisolatecli -b 5 test/simple /loop-darwin -- -i 5 -e 2
```
Tarball  with exit status 2
```
sudo ./fsisolatecli -b 5 -a 5 test/simple.tar /loop-darwin -- -i 10 -e 2
```
URL with exit status 0
```
sudo ./fsisolatecli -b 5 -a 5 https://raw.githubusercontent.com/odacremolbap/fsisolatecli/master/test/simple.tar /loop-darwin -- -i 10 -e 0
```

...


This is an output example

```
$ sudo ./fsisolatecli -b 5 -a 5 https://raw.githubusercontent.com/odacremolbap/fsisolatecli/master/test/simple.tar /loop-darwin -- -i 10 -e 2

[ISOCLI]Preparing the image, this might take a while ...
[ISOCLI]When running the process enter:
[ISOCLI]	's + enter' to get process state
[ISOCLI]	'i + enter' to send SIGINT
[ISOCLI]	'h + enter' to send SIGHUP
[ISOCLI]	'k + enter' to send SIGKILL
[ISOCLI]	'u + enter' to send SIGUSR1
s
[ISOCLI] PID: 0 STATE: not started
s
[ISOCLI] PID: 0 STATE: not started
[CHROOT]Iteration #1
[CHROOT]Iteration #2
[CHROOT]Iteration #3
s
[ISOCLI] PID: 23852 STATE: running
[CHROOT]Iteration #4
[CHROOT]Iteration #5
[CHROOT]Iteration #6
u
[CHROOT]>>>>>>>>> Signal received: user defined signal 1
[CHROOT]Iteration #7
[CHROOT]Iteration #8
[CHROOT]Iteration #9
k
[ISOCLI]Process finished with error: Error waiting process: signal: killed
```
