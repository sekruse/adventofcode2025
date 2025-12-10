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
  * I thought I was being smart about how I check only the perimeter and use the concept of orientation along with linear algebra to detect when we step out of bounds. But the algorithm was still taking aa few minutes. I suspect we could prune pairs of red tiles based on where previous candidate checks went out of bounds.
