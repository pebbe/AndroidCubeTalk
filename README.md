**under development**

## API ##



**request**          | **reply**
---------------------|---------------------------------------
`join` {id}          | `.`
`reset`              | `.`
`lookat` {x} {y} {z} | `enter` {id} {n0} <BR>
                     | `exit` {id} {n0} <BR>
                     | `self` {id} {n1} {x} {y} {z} <BR>
                     | `moveto` {id} {n2} {x} {y} {z} <BR>
                     | `lookat` {id} {n3} {x} {y} {z} <BR>
                     | `color` {id} {n4} {red} {green} {blue} <BR>
                     | `uncolor` {id} {n4} <BR>
                     | `.`

`join` must be the first request for every connection.

`reset` must be used once after the app has started.

`enter` must be used as reply before any other replies with the same `id`.

