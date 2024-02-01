# Canvaws API Automations

This contains various utilities for updating a Canvas course.

The canvas and canvasauto packages may be externally useful, but the core
purpose of this tool is for me to manage netacad courses.

## Netacad

Seriously, Cisco, it's 2024, make a GRPC or an API. It's not 1992. Grow up.

Since there is no fracking API for Netacad, some packages here expect a scrape
of straight up HTML, such as for swapping links for module items.

## Utilities

* Classes
	* Simply lists classes
* Students
	* Dumps students in a class, in the csv format Netacad expects
* Modules
	* Dumps module info
