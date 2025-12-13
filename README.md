# Advent of Code 2025

My attempt at cracking the Advent of Code 2025 puzzles.

```
    *
   /•\
  /≈≈≈\
 /o  • \
/_______\
   ||
```

## Notes

* Day 1:
  * Round 1 is straight forward. Round 2 has some tricky edge cases that are not featured in the example data. In particular, the actual output data has codes with full revolutions.
* Day 9:
  * I thought I was being smart about how I check only the perimeter and use the concept of orientation along with linear algebra to detect when we step out of bounds. But the algorithm was still taking a few minutes. I suspect we could prune pairs of red tiles based on where previous candidate checks went out of bounds.
* Day 12:
  * Every solution for round 1 can be trivially tranformed into another solution by mirroring it vertically and/or horizontally. So we could start building a solution into one direction.
  * Shapes can have various symmetries: horizontally, diagonally (2x), vertically, rotational (90 and 180 degrees). Detecting those prunes the search space considerably.
  * TODO: The current greedy solution was good enough to get me the star, but it fails the test. For this to do, I'd have to add backtracking (which would be prohibitively expensive without further optimization) or do something entirely different.
    * One idea is: We know the size of the regions and the size of the shapes. That would allow us to compute ahead of time how densely packed a hypothetical solution will have to be.
