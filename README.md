# loggo

loggo is my personal logging program. It's an experiment and I don't know what
direction it will take long-term. It's open source, but don't rely on it being
stable or not changing much for now. It will likely change rapidly.

Although it is in _very_ early development, here are some planned features:

* WSJT-X UDP Server integration
* LoTW integration
* Contest modes
* SKCC integration
* POTA and SOTA integration

The base idea of *loggo* is the idea of *tags*. The "comment" field of a QSO can
contain *tags* which are of the form `key:value` or `key:"multi-word
value"`. This allows you to dynamically tag "important" bits of your QSO for
filtering and exporting later on.

The only _required_ fields (which are _not_ tags) are: callsign, frequency,
mode, signal reports, and timestamp. All other fields end up being tags,
_in some way_, though this is still all experimental. The current thought is to
allow *loggo* to run in different "modes" which show different form fields
based on which mode is selected, and have the form fields correspond to tags.
In "general QSO" mode, tags can just be added in the comment field, so that if
K3XYZ suddenly tells you his SKCC number, you can just quickly type `skcc:` and
the number. This is all just in "idea" phase right now. General QSO logging
and ADIF exporting works but that is all for now.
