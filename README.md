# turtle-bunny
A Financial Transaction Database inspired by Tiger Beetle

## Info
Turtle Bunny uses a sqlite3 database alongside Golang.
Turtle Bunny can be used as a cli or can be imported as a Go Package.

Turtle Bunny attempts to behave like Tiger Beetle where possible.

Because sqlite3 only supports 64-bit signed integers while Tiger Beetle supports up to 128-bit unsigned integers the following design choices have been made.
- Turtle Bunny stores uint64 and uint128 values as TEXT fields.
- Turtle Bunny enforces the text field so that it can only be set as a valid unsigned 64 or 128 bit integer.

Turtle Bunny does not support two-phase transactions.
