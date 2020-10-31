# jvozba

[![](https://godoc.org/github.com/uakci/jvozba/v2?status.svg)](https://godoc.org/github.com/uakci/jvozba/v2)

An O(*n*) implementation of the lujvo-making algorithm to save the world.

Please use version 2 by importing the versioned path
`github.com/uakci/jvozba/v2`.

## What's the big deal with O(*n*), anyway?

All the jvozba I've seen over the years are exponential in complexity
(O(*c*^*n*), where 1 ≤ *c* ≤ 4) — the ‘algorithm’ they implement is basically
collecting all possible combinations of rafsi in one large array, mapping the
array with a score function, and sorting. This means that prefixing an input
tanru with just one `bloti` (a word which happens to have four affix forms)
*quadruples* the time (and memory, if the implementation is exponentially
naïve) it takes for such an algorithm to complete. To put this into
perspective: in order to find the lujvo for `bloti bloti bloti bloti bloti
bloti bloti bloti bloti bloti`, the algorithm has to call the score function a
*million* (4^10) times.

This jvozba, on the other hand, is linear in complexity, which means it can
compute even a million-`bloti` lujvo in about a second. It goes through each
tanru unit, keeping track of the best rafsi chains ‘so far’ alongside their
score, with a separate tally for tosmabru words for soundness. The algorithm
bears some resemblance to A\*.

Please don't read the code — it's pretty messy.
