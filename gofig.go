package gofig

//Private vars
var (
	loadedConfig *config = new(config)
)

//Public functions
func Str(id string) string {
	return loadedConfig.Val(id).(string)
}

//Private functions
func init() {
	loadedConfig.Populate()
}
