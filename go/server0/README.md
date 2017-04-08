
## Global layout ##

In the coordinate system, positive x points to the right, positive y
points up, and positive z points forward (towards the viewer).

The global layout is defined in file `user.go` after the comment `layout
is built from this list`. Here the cubes are positioned, usually in a
circle equidistant from the center, usually at level `y=0`, but in this
example one cube is located below and one above the zero level.

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
