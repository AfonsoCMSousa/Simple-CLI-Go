package comand

type Properties interface {
	Run(args []string) string
	Help() string
	Usage() string
}

type Command struct {
	Name        string
	Description string
	UsageText   string
	Execute     func(args []string) (content string)
}

func (c *Command) Run(args []string) string {
	return c.Execute(args)
}

func (c *Command) Help() string {
	return c.Description
}

func (c *Command) Usage() string {
	return c.UsageText
}
