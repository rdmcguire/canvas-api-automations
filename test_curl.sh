#!env sh
curl ${CANVAS_URL}v1/courses/${CANVAS_COURSE_ID}/modules/140294/items/146076 \
	-0 -v -X PUT \
	-H "Authorization: Bearer ${CANVAS_TOKEN}" \
	-H "Content-Type: application/json" \
	-H "Accept: application/json" \
	--data-binary @- << EOF
{
	"module_item": {
		"external_url": "https://lms.netacad.com/mod/lti/view.php?id=82047086",
		"title": "Chapter 1"
	}
}
EOF
