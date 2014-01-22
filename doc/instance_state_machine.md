Instance State Machine
------
This document describes the state machine for each instance status's processing.

Instance Status
======
Instace status include:
- Nil
- preAccepted
- accepted
- committed
- executed
- preparing

How Instance State Transition
======
*Nil* to
- *preparing*: 	dependency required
- *preAccepted*: receive proposal, receive preAccepted
- *accepted*: receive accepted
- *committed*: receive committed

*preAccepted* to
- *accepted*: slow path, receive accepted
- *committed*: fast path, receive committed
- *preparing*: timeout, dependency required

*accepted* to
- *committed*: majority agree, receive commmitted
- *preparing*: timeout, dependency required

*preparing* to
- *preAccepted*: (< N/2 preAccept), all noop
- *accepted*: (>= N/2 preAccept), (accepted in prepare reply), receive accepted
- *committed*: (committed in prepare reply), receive committed
- *preparing*: timeout

*committed* to
- *executed*: 

Transition Conditions Explained
======
Transition conditions:
- *dependency required*, happens when other committed instance depends on this instance and require it to be committed.
- *fast path*, happens when **initial** leader receives identical preAccept replies from fast quorum.
- *slow path*, happens when leader receives majority votes but it doesn't satisfy fast path conditions.
- *timeout*, happens when timeout in waiting for replies.
- *< N/2 preAccept*, happens when in preparing, less than N/2 replica reply preAccept messages.
- *>= N/2 preAccept*, happens when in preparing, majority reply preAccept messages.
- *all noop*, happens when all prepare replies of no-op commands.
- *majority agree*, happens when >= N/2 replica agree on the message you sent.
- *receive < XXX >*, happens when receiving message of < XXX > type.