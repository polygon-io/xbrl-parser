# xbrl-parser

[![GoDoc](https://godoc.org/github.com/polygon-io/xbrl-parser?status.svg)](https://godoc.org/github.com/polygon-io/xbrl-parser)

A Go library to parse xbrl documents into their facts, contexts, and units.

This library is based around the [XBRL 2.1 spec](https://www.xbrl.org/Specification/XBRL-2.1/REC-2003-12-31/XBRL-2.1-REC-2003-12-31+corrected-errata-2013-02-20.html).
It implements support for parsing basic facts (not tuples of facts), contexts and units through the `xml.Unmarshaler` interface.
 
See the package example in the godocs for how to unmarshal into the `XBRL` struct.

This library supports basic validation that checks for malformed facts and broken references between facts and contexts/units (see `XBRL.Validate()`),
but it does _not_ implement full semantic validation of XBRL documents.

There are no abstractions added on-top of the XBRL data structure, which makes this library flexible and simple,
but it also means you might have to read up a bit on how XBRL works to take full advantage of it.
