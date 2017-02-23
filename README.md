**under development**

## API ##



**request**          | **reply**
---------------------|---------------------------------------
`join` {id}          | `.`
`reset`              | `.`
`lookat` {x} {y} {z} | `self` {n0} {z} <BR> `enter` {id} {n1} <BR> `exit` {id} {n1} <BR> `moveto` {id} {n2} {x} {y} {z} <BR> `lookat` {id} {n3} {x} {y} {z} <BR> `color` {id} {n4} {red} {green} {blue} <BR> `info` {n5} {iid} {nr of lines} <BR> {lines} <BR> `.`

`join` must be the first request for every connection.

`reset` must be used once after the app has started.

`enter` must be used as reply before any other replies with the same `id`.

