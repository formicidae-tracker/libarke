# FORmicidae Tracker (FORT) : Inter-Sevice Communication Protocol

[![DOI](https://zenodo.org/badge/176954505.svg)](https://zenodo.org/doi/10.5281/zenodo.10019094)


The [FORmicidae Tracker (FORT)](https://formicidae-tracker.github.io) is an advanced online tracking system designed specifically for studying social insects, particularly ants and bees, FORT utilizes fiducial markers for extended individual tracking. Its key features include real-time online tracking and a modular design architecture that supports distributed processing. The project's current repositories encompass comprehensive hardware blueprints, technical documentation, and associated firmware and software for online tracking and offline data analysis.

This repository holds specifications, device (AVR) and host (Linux/Go) implementation of the CAN based communication protocol used by electronic devices of the FORT Project.


## Specifications

Complete specifications can be found in [specs/specs.md](https://github.com/formicidae-tracker/libarke/blob/master/specs/specs.md)

## Implementation

Two implementations of the protocol are currently available:
* AVR C implementation, currently only supporting AT90CANXXX processors.
* Golang implementation based on socketcan, only supporting the linux platform.

## Installation


## License

This project is licensed under the GPL version 3


