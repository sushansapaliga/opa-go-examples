package authz

test_allow_public_get {
	allow with input as {"method": "GET", "path": "/public"}
}

test_deny_public_post {
	not allow with input as {"method": "POST", "path": "/public"}
}

test_allow_admin {
	allow with input as {"user": "admin", "method": "POST", "path": "/protected"}
}

test_deny_regular_user {
	not allow with input as {"user": "bob", "method": "POST", "path": "/protected"}
}
