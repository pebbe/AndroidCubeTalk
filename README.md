**under development**

## API ##



**request**          | **reply**
---------------------|---------------------------------------
`join` <id>          | `.`
---------------------|---------------------------------------
`reset`              | `.`
---------------------|---------------------------------------
`lookat` <x> <y> <z> | `enter` <id> <n0>
                     | `exit` <id> <n0>
                     | `self` <id> <n1> <x> <y> <z>
                     | `moveto` <id> <n2> <x> <y> <z>
                     | `lookat` <id> <n3> <x> <y> <z>
                     | `color` <id> <n4> <red> <green> <blue>
                     | `uncolor` <id> <n4>
                     | `.`

`join` must be the first request for every connection.

`reset` must be used once after the app has started.

`enter` must be used as reply before any other replies with the same `id`.

