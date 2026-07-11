# Interface test

This exercise tests your understanding of Go interfaces.

interfaces.go defines 2 interfaces, Input and Output.

The test has an array of interface implementations. It iterates through those,
finds the Inputs, and calls their NextValue method.
Then, it iterates through the interfaces again to find all the
Outputs and sends each value to every Output using the SendValue
method.

At the end of the test, the values sent to one of the Outputs
are checked. As provided, some Interfaces do not behave as expected
and must be fixed. The test file (interfaces_test.go) does not
need to be modified. The list package is a helper and also
does not need any changes, however looking at what it does
might help solve the exercise.

In short, make the test pass ;-).
