# json-to-golang


## About

[json-to-golang](http://json-to-golang.com) is a simple web app to convert raw json to a golang struct.

The main logic is in `server/lib/json_str.go`.

Tests can be run from `./server` with `go test -v ./lib`


## Technical Notes

### Preserving Key Order of JSON

On the back-end, when using the provided JSON decoder to unmarshal a json string to 
a map[string]interface{} data type, the order of the keys will be lost when iterating
over all the keys in the resulting map.

To fix this, a package called "go-ordered-json" is used which retains the order of keys.

Here is some more info about that package
* Code: https://gitlab.com/c0b/go-ordered-json
* Article: https://medium.com/@ty0h/preserving-json-object-keys-order-in-javascript-python-and-go-language-170eaae0de03
* Docs: https://godoc.org/gitlab.com/c0b/go-ordered-json


## Manual Testing

### Good JSON to Use

```
{
	"id": 1296269,
	"owner": {
		"login": "octocat",
		"id": 1,
		"avatar_url": "https://github.com/images/error/octocat_happy.gif",
		"gravatar_id": "somehexcode",
		"url": "https://api.github.com/users/octocat",
		"html_url": "https://github.com/octocat",
		"followers_url": "https://api.github.com/users/octocat/followers",
		"repos_url": "https://api.github.com/users/octocat/repos",
		"events_url": "https://api.github.com/users/octocat/events{/privacy}",
		"received_events_url": "https://api.github.com/users/octocat/received_events",
		"type": "User",
		"site_admin": false
	},
	"name": "Hello-World",
	"full_name": "octocat/Hello-World",
	"has_wiki": true,
	"has_downloads": true,
	"pushed_at": "2011-01-26T19:06:43Z",
	"created_at": "2011-01-26T19:01:12Z",
	"updated_at": "2011-01-26T19:14:43Z",
	"permissions": {
		"admin": false,
		"push": false,
		"pull": true
	},
	"subscribers_count": 42,
	"organization": {
		"login": "octocat",
		"id": 1,
		"avatar_url": "https://github.com/images/error/octocat_happy.gif",
		"gravatar_id": "somehexcode",
		"site_admin": false
	},
	"parent": {
		"id": 1296269,
		"owner": {
			"login": "octocat",
			"id": 1,
			"avatar_url": "https://github.com/images/error/octocat_happy.gif",
			"gravatar_id": "somehexcode",
			"type": "User",
			"site_admin": false
		},
		"name": "Hello-World",
		"full_name": "octocat/Hello-World",
		"description": "This your first repo!",
		"private": false,
		"fork": true,
		"homepage": "https://github.com",
		"language": null,
		"forks_count": 9,
		"stargazers_count": 80,
		"watchers_count": 80,
		"size": 108,
		"default_branch": "master",
		"open_issues_count": 0,
		"has_issues": true,
		"has_wiki": true,
		"has_downloads": true,
		"pushed_at": "2011-01-26T19:06:43Z",
		"created_at": "2011-01-26T19:01:12Z",
		"updated_at": "2011-01-26T19:14:43Z",
		"permissions": {
			"admin": false,
			"push": false,
			"pull": true
		}
	},
	"source": {
		"id": 1296269,
		"owner": {
			"login": "octocat",
			"id": 1,
			"received_events_url": "https://api.github.com/users/octocat/received_events",
			"type": "User",
			"site_admin": false
		},
		"name": "Hello-World",
		"full_name": "octocat/Hello-World",
		"description": "This your first repo!",
		"private": false,
		"fork": true,
		"mirror_url": "git://git.example.com/octocat/Hello-World",
		"homepage": "https://github.com",
		"language": null,
		"forks_count": 9,
		"watchers_count": 80,
		"size": 108,
		"default_branch": "master",
		"open_issues_count": 0,
		"has_issues": true,
		"has_wiki": true,
		"has_downloads": true,
		"pushed_at": "2011-01-26T19:06:43Z",
		"created_at": "2011-01-26T19:01:12Z",
		"permissions": {
			"admin": false,
			"pull": true
		}
	}
}
```
