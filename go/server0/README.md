
## Running server0

Run as:

    ./server0 config_file.json

To see what options are available, run this command:

    ./server0 layout-minimal.json

You can then see a full layout (in json format) at the top of
the log file.

In the config file, you either define a layout of cubes with
the option `cubes`, or you just give a list of user names in
`users`. Though they are both filled out in the log file, you
can't use both at the same time as input.

When you use `users`, a layout is created putting all users in
a circle, but if there is a robot that masks another user, it is
put away from the circle. The values for `default_color`,
`default_face`, `default_head`, and `default_skip_gui` are used to
fill in the details for the cubes.

When you use `cubes`, the list of users if created from the UIDs
of the cubes. When a cube has no color, the value of
`default_color` is used. The values for `default_face`,
`default_head`, and `default_skip_gui` are ignored.



## Global layout ##

In the coordinate system, positive x points to the right, positive y
points up, and positive z points forward (towards the viewer).

The global layout is created from the set of cubes defined in the config
file, or, if it is missing in the config file, a layout is generated for
the users listed in the config file.


## Per user layout ##

For each user, the global layout is transformed to a personal layout.
This is stored in the map `users`. The global layout is rotated, (and in
some cases raised or lowered), so the user is located at the positive
z-axis, at `0,0,sqrt(x*x+z*z)`. Forward is a vector `0,0,-1` (all vectors
are normalised to unit length), pointing to the origin. When you send
the `recenter` command from the GUI to a user, the whole scene he sees
is rotated so the y-axis is directly in front of where his head is
pointed to at that moment.

Locations for other cubes are shifted around so the global layout
matches the personal layout. Each cube gets a `forward` vector pointing
towards the y-axis:

    angle = atan2(x, z)
    forward = -sin(angle), 0, -cos(angle)

## How do you know if someone is looking at you? ##

Is A looking at B?

    a := labels["A"]
    b := labels["B"]
    v := users[a].lookat
    w := users[a].cubes[b].towards
    v.x * w.x + v.y * w.y + v.z * w.z > 0.99
