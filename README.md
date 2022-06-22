# itp - Beautiful Prime Numbers from Images

<samp>itp</samp> is a command line tool to generate large prime numbers which look like any arbitrary image.

<div align="center">
  <a href="https://gist.github.com/raklaptudirm/fc6fd8452fefc246951efca730e9e214">
    <img src="https://user-images.githubusercontent.com/68542775/174952364-39ce9031-6ab0-4e7b-a15f-2665bdc88056.png">
  </a>
</div>


### Installation and Usage

```
go install laptudirm.com/x/itp@latest
itp [image file]
```

### Working

<samp>itp</samp> finds primes which look similar to an image using the following steps:

- `step 1` Convert the image into a number. Brightness chart is "7772299408".
- `step 2` If last digit is even or 5 change it, otherwise number can't be a prime.
- `step 3` If first digit is 0 change it. No redundant leading zeros allowed
- `step 4` Check if number is prime. If it is a prime, print it and exit.
- `step 5` If number is not prime, switch one of the digits with a similar one.
- `step 6` Goto step 4.

### References

- https://github.com/TotalTechGeek/pictoprime

- https://en.wikipedia.org/wiki/Baillie%E2%80%93PSW_primality_test
- https://en.wikipedia.org/wiki/Miller%E2%80%93Rabin_primality_test
