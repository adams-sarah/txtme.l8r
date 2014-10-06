txtme.l8r
=========

a snooze button.. for your todos

If txtme.l8r were a production application, it would probably be wise to have some persistent backup of the waiting-to-be-sent messages. As currently implemented, an application restart (after a fresh deploy, after a crash, etc) would kill all running goroutines and their messages would never be sent.


Enjoy!
