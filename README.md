<h1 align="left">Trash-rm</h1>

Trash-rm is a little tool for the linux terminal written in golang. With this tool you can
kind of soft delete files and folders. What do i mean with soft delete?
It just means when you delete a file as example and it creates, in the same directory, a folder
named .trash (hidden), compress the file and move it in there. It also creates a database
file with sqlite and makes an entry. The point is that you can decide later to restore that file,
when you changed your mind and you need it again.

With rm on linux you delete files and folders permantly and you cant recover it. So i created this
little tool.

<h1 align="left">How do i use it?</h1>
Clone that repository, make sure you have golang installed and execute the build.sh script. It
compiles the project to a binary in the bin/ directory.
Then you can use
```
./trm help
```
to show all available commands.
