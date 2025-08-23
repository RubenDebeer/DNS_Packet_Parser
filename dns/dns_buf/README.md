# Explain the Code dude.
// TODO: Draw out the Actual Packet so that I can better understand it.


# QNAME Decomposition
Query consists of 3 parts: the Query name, The Type and the Class. DNS encodes each name into a sequence of labels, with each label prepended by a single **byte** indicating it's length. IF We use "google.com" as a example "google" is 6 bytes and is thus preceded by **0x06**, while "com" is 3 bytes and is preceded by **0x03**. Finally, all names are terminated by a label of **zero length**, that is a null byte.

**Quote from FRC 1034**
"Each node has a label, which is zero to 63 octets in length.  Brother
nodes may not have the same label, although the same label can be used
for nodes which are not brothers.  One label is reserved, and that is
the null (i.e., zero length) label used for the root."

[3]www[6]google[3]com[0]

That’s: 0x03 'w' 'w' 'w'  0x06 'g' ... 'e'  0x03 'c' 'o' 'm'  0x00 → www.google.com