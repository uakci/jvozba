# jvozba

An O(*n*) implementation of the lujvo-making algorithm to save the world.

## What's the big deal with O(*n*), anyway?

All the jvozba I've seen over the years are of exponential complexity (O(*c*^*n*), where 1 ≤ *c* ≤ 4), because the ‘algorithm’ they implement is basically collecting all possible combinations of rafsi in an array, mapping the array with a score function, and sorting. This means that prefixing an input tanru with just one `bloti` will *quadruple* the time and memory it takes for the lujvo to compute. To put this into perspective: in order to find the lujvo for `bloti bloti bloti bloti bloti bloti bloti bloti bloti bloti`, the algorithm will have to call the score function a *million* times. Double the input length and your 32-bit machine will explode. (Or wake up the OOM killer.)

This jvozba, on the other hand, is linear in complexity, which means it can compute even a million-`bloti` lujvo in about a second. ‘How does it achieve *that*?’, I hear you ask. Simply put, it goes through each tanru unit, keeping track of the best lujvo ‘so far’ alongside its score, with a separate tally for tosmabru words for soundness. There's a bunch more performance tweaks in the code – I encourage you to perhaps read it.

## Usage

`main.go` should give you an idea of how to use the (extremely simple) basic API. If you want to customise stuff, dig into the code and you should find the right procedures to call.
