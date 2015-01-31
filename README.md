# mozoval

## Overview

The mozoval project is a set of experimental OVAL security processing
modules developed by Mozilla.

The primary tool set under active development is a go based OVAL library
and associated command line processor. This processor can read certain
elements of published OVAL checks and return results.

## Objective

The objective of the project is to eventually integrate OVAL checking into
the Mozilla Investigator (mig), to permit agent based vulnerability checks
on various systems. The intent is to implement enough of the OVAL specification
to allow for execution of the most applicable types of checks. In addition,
to extend support where required to OVAL to support simplified analysis of
other types of vulnerabilities.

## Status

The project is currently in it's infancy, supporting a limited subset of
checks to be expanded in the future.
