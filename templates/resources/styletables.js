tags = ["tr", "th", "td"];
for (var j=0; j < tags.length; j++ ) {
	var elems = document.getElementsByTagName(tags[j]);
	for (var i=0; i < elems.length; i++) {
	     elems[i].className = tags[j] + "Class";
	}
}