
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
matches the personal layout. Each cube gets a forward vector pointing
towards the y-axis:

	rotate = atan2(x, z)
	forward = -sin(rotate), 0, -cos(rotate)
	
## How do you know if someone is looking at you? ##

To know if someone is looking at you, you need two things:

 * Where is the other one located relative to me?
 * What direction is the other one looking at, relative to the direction
   from him to me?
   
Suppose you are `A` and the other one is the first cube.

 1. You are at `0, 0, users["A"].selfZ`
 2. The first cube is located at `users["A"].cubes[0].pos`
 3. The first cube is currently looking at `users[users["A"].cubes[0].uid].lookat`
 4. The vector from 3 is relative to vector `users["A"].cubes[0].forward`
 
(You can see these values in the log file.)
 
You need to rotate (horizontally and sometimes vertically as well) so
the vector from number 4 is pointed in the same direction as the vector
pointing from number 2 to number 1. Then you need to rotate the vector
from number 3 with the same amount. Finally, you can calculate the cosine
of the angel between the vector from 2 to 1, and the rotated vector 4. If
this value is (nearly) 1, then the other one is looking at you.

The cosine of the angel between two unit vectors v and w is:

    v(x) * w(x) + v(y) * w(y) + v(z) * w(z)
	
(For non-unit vectors, you need to devide the above by the product of
the lengths of the two vectors.)

Remember, all vectors (`forward` and `lookat`) are normalised to unit
length.
