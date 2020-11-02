Go 3270 Server Library
======================

This library allows you to write Go servers for tn3270 clients by building 3270 data streams from fields and processing the client's response to receive the attention keys and field values entered by users.

**The library is incomplete, likely buggy, and under heavy development: the interface is UNSTABLE until this notice is removed from this readme and version 1.0 is released.**

Everything I know about 3270 data streams I learned from [Tommy Sprinkle's tutorial][sprinkle]. The tn3270 telnet negotiation is gleaned from [RFC 1576: TN3270 Current Practices][rfc1576], [RFC 1041: Telnet 3270 Regime Option][rfc1041], and [RFC 854: Telnet Protocol Specification][rfc854].

[sprinkle]: http://www.tommysprinkle.com/mvs/P3270/
[rfc1576]: https://tools.ietf.org/html/rfc1576
[rfc1041]: https://tools.ietf.org/html/rfc1041
[rfc854]: https://tools.ietf.org/html/rfc854

Known Problems
--------------

 - The telnet data is not checked for the special telnet byte value, 0xFF, which requires escaping while sending and unescaping while receiving data. If your 3270 data streams contain an FF character, things will probably break.
 - The telnet negotiation does not check for any errors or for any responses from the client. We just assume it goes well and we're actually talking to a tn3270 client.

License
-------

This library is licensed under the MIT license; see the file LICENSE for details.