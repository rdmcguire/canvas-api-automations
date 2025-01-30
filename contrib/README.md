# Updating Links

Each session you will have to do this bullsh!t. 1m thanks Netacad! An API would just be AWFUL. Welcome to 2025 with Cisco.

1. Log in to Netacad, and open the course
1. Right click the page and hit inspect
1. Go to the container surrounding all of the links
1. Right click the div and hit copy
1. Drop it into a file, run :LspStop and :syntax off, write it out
1. Run parse_command.sh on the file
1. Run canvas-api-automations modules update <parsed_file>
1. Flip Cisco off, two middle fingers if you got 'em'
