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
* Assignments
    * This updates all links for netacad assignments in a course
    * You must download the html of the assignments page, and reference as command arg
        * The html comes from the "Course Home", with all the links
    * Run the command like `canvas-api-automations assignments course.html 99999` where 99999 is the
      id of the canvas course you found running the classes subcommand
