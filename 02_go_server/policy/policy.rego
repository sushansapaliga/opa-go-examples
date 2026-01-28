package authz

default allow = false

# Allow access to public endpoint
allow {
	input.method == "GET"
	input.path == "/public"
}

# Allow admin user to access anything
allow {
	input.user == "admin"
}
