package membercache_test

var configAll = `
{
	"Config":{
		"Type":"",
		"Cache":{
			"Marshaler":"json",
			"Driver":"dummycache"
		}
	}
}
`

var configStatus = `
{
	"Config":{
		"Type":"status",
		"Cache":{
			"Marshaler":"json",
			"Driver":"dummycache"
		}
	}
}
`

var configAccounts = `
{
	"Config":{
		"Type":"accounts",
		"Cache":{
			"Marshaler":"json",
			"Driver":"dummycache"
		}
	}
}
`

var configToken = `
{
	"Config":{
		"Type":"token",
		"Cache":{
			"Marshaler":"json",
			"Driver":"dummycache"
		}
	}
}
`

var configRole = `
{
	"Config":{
		"Type":"role",
		"Cache":{
			"Marshaler":"json",
			"Driver":"dummycache"
		}
	}
}
`

var configData = `
{
	"Config":{
		"Type":"data",
		"Cache":{
			"Marshaler":"json",
			"Driver":"dummycache"
		}
	}
}
`

var configError = `
{
	"Config":{
		"Type":"error",
		"Cache":{
			"Marshaler":"json",
			"Driver":"dummycache"
		}
	}
}
`
