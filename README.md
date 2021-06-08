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

To give you a head start, here's some basics about XBRL:

## What is XBRL?

At a high level, XBRL is an XML spec intended to be the standard for digital business reporting.

I can't give you a more precise definition because "digital business reporting" can mean many things, 
and XBRL needs to work for all of them.

XBRL is useless on its own, it relies on supplemental taxonomies that describe facts for specific use cases.
For example the [US GAAP Financial Reporting Taxonomy](https://xbrl.us/xbrl-taxonomy/2021-us-gaap/) defines schemas for
facts that relate to US GAAP financial reporting, which is used heavily in quarterly reports submitted to the SEC, among many other things.

If you're here, I'm assuming you have a specific reason to want to parse XBRL documents, so I'm not going to go into more detail on this.

Let's jump right into the main components of an XBRL document... 

### Facts

The core of any XBRL document: facts (or [items](https://www.xbrl.org/Specification/XBRL-2.1/REC-2003-12-31/XBRL-2.1-REC-2003-12-31+corrected-errata-2013-02-20.html#_4.6) as the spec calls them)
represent a single business measurement.
Here's an example of a fact whose schema is defined in the [US GAAP Financial Reporting Taxonomy](https://xbrl.us/xbrl-taxonomy/2021-us-gaap/):

```xml
<us-gaap:EarningsPerShareBasic contextRef="c1" decimals="2" unitRef="u1">1.41</us-gaap:EarningsPerShareBasic>
```

A fact by itself is only a fragment of a useful piece of information. 
In the above example we see that earnings per share is `1.41`,
but we need more context around when this fact was true, and how it was measured.

Let's start with the context...

### Contexts

A [Context](https://www.xbrl.org/Specification/XBRL-2.1/REC-2003-12-31/XBRL-2.1-REC-2003-12-31+corrected-errata-2013-02-20.html#_4.7)
describes a business entity, period of time, and an optional scenario (this library doesn't currently support scenarios, so we're going to gloss over them).  

When a fact references a context, it gives the fact more detail to help us understand what it means.

Note that many facts can reference the same context where it makes sense to do so.

The fact in the above example references a context called "c1", let's see what that context might look like:
```xml
<context id="c1">
    <entity>
        <identifier scheme="http://www.sec.gov/CIK">0000320193</identifier>
    </entity>
    <period>
        <startDate>2020-12-27</startDate>
        <endDate>2021-03-27</endDate>
    </period>
</context>
```

With the information in this context, we now know that in Q1 of 2021 (between 2020-12-27 and 2021-02-37), Apple Inc.'s (CIK 0000320193) EPS was 1.41.

We're closer to having a useful piece of information now, but there's one thing we're still missing.
EPS is 1.41...what? What unit are we measuring it in?

### Units

A [Unit](https://www.xbrl.org/Specification/XBRL-2.1/REC-2003-12-31/XBRL-2.1-REC-2003-12-31+corrected-errata-2013-02-20.html#_4.8)
describes a unit of measure for a numeric fact.

A unit can represent something simple like number of shares, 
something slightly more complex like dollars per share, 
or any other kind of unit you can think of.

Just like contexts, more than one fact can reference the same unit when it makes sense to do so.

Note that only numeric facts have units. 
Sometimes a fact is a block of text, which doesn't make sense to have a unit.

Let's look at what a simple unit like number of shares might look like:
```xml
<unit id="shares">
    <measure>shares</measure>
</unit>
```

That's great and all...but the fact in the above example references a unit called "u1", 
let's see what that more complex unit might look like:
```xml
<unit id="u1">
    <divide>
        <unitNumerator>
            <measure>iso4217:USD</measure>
        </unitNumerator>
        <unitDenominator>
            <measure>shares</measure>
        </unitDenominator>
    </divide>
</unit>
```

This unit is a ratio of two simple units: USD / shares.

And with that we now fully understand the fact from the example above:

In Q1 of 2021 (between 2020-12-27 and 2021-02-37), Apple Inc.'s (CIK 0000320193) EPS was 1.41 dollars per share.

---

### Wrapping Up

That was a _very_ brief overview of the XBRL format, 
hopefully it empowers you enough to understand the basics of how to get information out of an XBRL document.

If you need to dig a little deeper, the models in this library are well documented and contain links to their definitions in the XBRL spec for your reference.
Beyond that, the [official spec](https://www.xbrl.org/Specification/XBRL-2.1/REC-2003-12-31/XBRL-2.1-REC-2003-12-31+corrected-errata-2013-02-20.html)
is a bit large, but it's pretty clear and will almost definitely have the information you need (and probably a lot you don't need too!)
