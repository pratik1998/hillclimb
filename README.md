# hillclimb

This repository contains code for running the hillclimb attack to decrypt the cipher text encrypted using the M4 enigma machine with some initial settings known.

Known settings are as follows:
- Rotors: ?? ?? IV III
- Rotors Initial Start Position: ?? ?? B Q
- Ringstellung: 1 1 1 16
- Reflector: C-Thin

So this engima decrypter gives you the plugboard setting and remaining rotor order and their start position.

## IMPORTANT
This repository uses the engima simulator code from https://github.com/emedvedev/enigma repository to encrypt message with initial enigma settings. So, most of the code inside the engima directory is from that repository.

# Installation

You need to have `go` installed in your system to run the code. After that you just need to run `go build` command to build the project.

# How to Run

You need to pass ciphertext file name in the argument to `hillclimb` executable file. For example, if ciphertext contained in a file named ct.txt then you need to run `./hillclimb ct.txt` command.

# References

- https://github.com/matthewdgreen/practicalcrypto/tree/master/Assignments
- https://cryptocellar.org/bgac/HillClimbEnigma.pdf
- http://practicalcryptography.com/cryptanalysis/letter-frequencies-various-languages/english-letter-frequencies/
