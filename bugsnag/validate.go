package bugsnag

import "fmt"

func validateValueFunc(validValues []string) func(val interface{}, key string) (warns []string, err []error) {
	return func(val interface{}, key string) (warns []string, err []error) {
		v := val.(string)
		isValid := false
		for _, t := range validValues {
			if t == v {
				isValid = true
				break
			}
		}
		if !isValid {
			return nil, []error{fmt.Errorf("unrecognized value '%s'", v)}
		}
		return nil, nil
	}
}

func validProjectTypes() []string {
	return []string{
		"rails",
		"django",
		"php",
		"laravel",
		"lumen",
		"magento",
		"silex",
		"symfony",
		"wordpress",
		"android",
		"ndk",
		"ios",
		"osx",
		"tvos",
		"cocos2dx",
		"reactnative",
		"expo",
		"sinatra",
		"rack",
		"node",
		"unity",
		"js",
		"java",
		"java_desktop",
		"spring",
		"python",
		"wsgi",
		"flask",
		"bottle",
		"tornado",
		"eventmachine",
		"ruby",
		"express",
		"connect",
		"koa",
		"restify",
		"angular",
		"ember",
		"backbone",
		"react",
		"vue",
		"heroku",
		"go",
		"go_net_http",
		"martini",
		"revel",
		"gin",
		"negroni",
		"dotnet",
		"dotnet_desktop",
		"aspnet",
		"dotnet_mvc",
		"webapi",
		"wpf",
		"aspnet_core",
		"other",
		"other_desktop",
		"other_mobile",
		"other_tv",
	}
}
