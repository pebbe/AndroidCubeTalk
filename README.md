**under development**

## API ##



**request**          | **reply**
---------------------|---------------------------------------
`join` {id}          | `.`
`quit`               | nothing, the connection is terminated
`reset`              | `.`
`lookat` {x} {y} {z} {roll} {audio} [mark] <BR> `info` {infoID} {choice} | `self` {n0} {z} <BR> `recenter` <BR> `enter` {id} {n1} <BR> `exit` {id} {n1} <BR> `moveto` {id} {n2} {x} {y} {z} <BR> `lookat` {id} {n3} {x} {y} {z} {roll} <BR> `color` {id} {n4} {red} {green} {blue} <BR> `info` {n5} {nr of lines} <BR> {lines} <BR> `info` {n5} {nr of lines} {infoID} {choice1} {choice2} <BR> {lines} <BR> `cubesize` {n6} {width} {height} {depth} <BR> `head` {id} {n7} {number} <BR> `face` {id} {n8} {number} <BR> `.`

All requests are a single line. All responses to `lookat` and `info` can
be multiple lines. The last line is a single dot, except for the `quit`
command, which doesn't send a reply.

A connection can be terminated by the client with the `quit` command, but
that is optional. It can also just close the connection.

Each response command is a single line, except for {info}, which is followed by
the text to be displayed.

`join` must be the first request on every connection.

`reset` must be used once after the app has started. This tells the
server to send the initial layout of the other cubes.

`enter` must be used as reply before any other replies with the same `id`.

`exit` only hides a cube. It doesn't reset color or position. It can be
unhidden with `enter`.

All parameters are a single word, except for {lines}.

{n0}...{n8} are separate counters. The client should ignore any reply
command that has a counter value lower than already seen for that
particular counter. This is in case replies don't get processed in the
order they were sent.

{red}, {green}, {blue} between 0 and 1.

{roll} in degrees.

{x} to the right, {y} up, {z} to the front.

Every client sees itself placed on the positive z-axis, at a distance
from the origin set by the `self` command.

No cube should be placed more than 10 units from the origin.

Example configuration, two people:

    self:  0, 0,  4
    other: 0, 0, -4

Three people:

    self:     0,    0,  4
    other1:   3.46, 0, -2
    other2:  -3.46, 0, -2

Four people:

    self:     0, 0,  4
    other1:   4, 0,  0
    other2:   0, 0, -4
    other3:  -4, 0,  0

Each cube extends 1 unit in each direction from its origin. This can be changed with the `cubesize` command.

Info panels are located 3 units from the viewer.
