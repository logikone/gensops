package config

type Config struct {
	AWS AWSConfig `yaml:"aws,flow"`
	Map []struct {
		RegEx string              `yaml:"regex"`
		Tags  []map[string]string `yaml:"tags,flow"`
	}
}

type AWSConfig struct {
	Accounts    []AWSAccount `yaml:"accounts,flow"`
	IncludeTags []string     `yaml:"include_tags,flow"`
}

type AWSAccount struct {
	ID      string     `yaml:"id"`
	IAMRole AWSIAMRole `yaml:"iam_role,flow"`
	Regions []string   `yaml:"regions,flow"`
}

type AWSIAMRole struct {
	Name   string `yaml:"name"`
	Prefix string `yaml:"prefix"`
}
