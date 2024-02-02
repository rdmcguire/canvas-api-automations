# Canvas API Automations

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

## Weird Stuff

It's great that Canvas produced an API. Thanks Canvas. Pay attention Cisco.

There are, however, some crappy things about the Canvas API:

* Postman schema is written for form data instead of json payload
    * This breaks most auto-generated OpenAPI code
    * See note about custom RoundTripper below
* Data types are downright wrong in some places
    * Manually correcting auto-generated code sucks, but I had no choice
    * Some fields (such as `PointsPossible` are generated as an int, but can be floats.
    * You may find more. You'll have to fix them. I fixed what I needed.
* Canvas wrote a GraphAPI, but didn't bother producing an SDL or providing an introspection endpoint
    * Because of this, you can't auto-generate code for their GraphAPI.. So I used the REST API

Because of this, I had to do two icky things:

1. Modify auto-generated code to correct field types (int -> float32)
1. Override the default http Transport with a custom RoundTripper

### About this RoundTripper

Basically the RoundTripper guts any PUT request Body, instead
loading it up with FormData. The OpenAPI spec leads the code generator to
believe it's producing `application/json` payload, when instead it should
be producing `application/x-www-form-urlencoded`.

What this looks like, is the code itself would have a json body such as:
```json
{
  "field_1[sub_field_1]": "someVal",
  "field_1[sub_field_2]": "otherVal"
}
```

What it SHOULD have is:
```json
{
  "field_1": {
      "sub_field_1": "someVal",
      "sub_field_2": "otherVal"
  }
}
```

Well crap. Also, for some reason Canvas doesn't even accept the second, I couldn't
even force mangle this bullshit into a map[string]any with nested fields as a
map[string]map[string]any without the thing cranking back a 400 on me.

So, I gut the request, replace it with a new one, load it up with form data,
and off to the races we go. You don't have to know this. You can just use
the package through the client. You're welcome.

The RoundTripper also sets auth, something the auto-gen code failed to accomodate.
