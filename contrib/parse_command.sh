#!env bash
if [ $# -ne 1 ]; then
	echo "USAGE: $0 <file.html>"
	exit 1
fi

ASSIGNMENTS=$1
if [ ! -f $ASSIGNMENTS ]; then
	echo "File $ASSIGNMENTS does not exist"
	exit 2
fi

echo "Parsing ${ASSIGNMENTS}"
\cat ${ASSIGNMENTS} | sed 's/<a/\r\n<a/g' | grep 'launch' > parsed_${ASSIGNMENTS}
